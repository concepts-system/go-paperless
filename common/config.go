package common

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kpango/glg"
)

const (
	//JWTAlgorithmHS256 specifies the algorithm name for HS256.
	JWTAlgorithmHS256 = "HS256"
	//JWTAlgorithmHS384 specifies the algorithm name for HS384.
	JWTAlgorithmHS384 = "HS384"
	//JWTAlgorithmHS512 specifies the algorithm name for HS512.
	JWTAlgorithmHS512 = "HS512"

	keyProfiles          = "PROFILES"
	keyPublicURL         = "PUBLIC_URL"
	keyPort              = "PORT"
	keyDatabaseType      = "DB_TYPE"
	keyDatabaseURL       = "DB_URL"
	keyDatabaseMigrate   = "DB_MIGRATE"
	keyJWTAlgorithm      = "JWT_ALGORITHM"
	keyJWTKey            = "JWT_KEY"
	keyJWTExpirationTime = "JWT_EXPIRATION_TIME"
	keyJWTRefreshTime    = "JWT_REFRESH_TIME"
	keyDataPath          = "DATA_PATH"

	developmentProfile = "development"
	productionProfile  = "production"

	defaultPublicURL         = "http://localhost:8080"
	defaultPort              = 8080
	defaultDatabaseType      = "sqlite3"
	defaultDatabaseURL       = "database.db"
	defaultJWTAlgorithm      = "HS256"
	defaultJWTExpirationTime = 5 * time.Minute
	defaultJWTRefreshTime    = 72 * time.Hour
	defaultDataPath          = "data"
)

var (
	jwtSupportedAlgorithms = []string{
		JWTAlgorithmHS256,
		JWTAlgorithmHS384,
		JWTAlgorithmHS512,
	}
)

// Configuration is a struct which holds all configurable properties.
type Configuration struct {
	profiles          []string
	publicURL         *url.URL
	port              int
	databaseType      string
	databaseURL       *url.URL
	migrateDatabase   bool
	jwtAlgorithm      string
	jwtKey            []byte
	jwtExpirationTime time.Duration
	jwtRefreshTime    time.Duration
	dataPath          string
}

// IsProfileActive checks whether the given profile is configured as active.
func (c *Configuration) IsProfileActive(profile string) bool {
	profile = strings.TrimSpace(profile)

	for _, p := range c.profiles {
		if profile == p {
			return true
		}
	}

	return false
}

// IsProduction returns a boolean value indicating whether the development profile
// is currently active.
func (c *Configuration) IsProduction() bool {
	return c.IsProfileActive(productionProfile)
}

// IsDevelopment returns a boolean value indicating whether the development profile
// is currently active.
//
// Development profile is only recognized as active when 'production' profile is not set.
//
func (c *Configuration) IsDevelopment() bool {
	if c.IsProduction() {
		return false
	}

	return c.IsProfileActive(developmentProfile)
}

// GetPublicURL gets the configured server base URL for public access.
func (c *Configuration) GetPublicURL() *url.URL {
	return c.publicURL
}

// GetPort gets the configured listening port of the HTTP server.
func (c *Configuration) GetPort() int {
	return c.port
}

// GetDatabaseType gets the configured type of database. Can be 'sqlite3', 'postgres' or 'mysql'.
func (c *Configuration) GetDatabaseType() string {
	return c.databaseType
}

// GetDatabaseURL returns the URL for establishing the connection to the database.
func (c *Configuration) GetDatabaseURL() *url.URL {
	return c.databaseURL
}

// MigrateDatabase returns a boolean value indicating whether the migrations should be run.
func (c *Configuration) MigrateDatabase() bool {
	return c.migrateDatabase
}

// GetJWTAlgorithm gets the configured algorithm for issueing and verifying JWTs.
func (c *Configuration) GetJWTAlgorithm() string {
	return c.jwtAlgorithm
}

// GetJWTKey gets the configured encryption/decryption key for issueing and verifying JWTs.
func (c *Configuration) GetJWTKey() []byte {
	return c.jwtKey
}

// GetJWTExpirationTime gets the expiration time of access tokens.
func (c *Configuration) GetJWTExpirationTime() time.Duration {
	return c.jwtExpirationTime
}

// GetJWTRefreshTime gets the expiration time of refresh tokens.
func (c *Configuration) GetJWTRefreshTime() time.Duration {
	return c.jwtRefreshTime
}

// GetDataPath gets the base path for storing documents.
func (c *Configuration) GetDataPath() string {
	return c.dataPath
}

// Configuration singleton
var config *Configuration

// Config returns the singleton config instance.
func Config() *Configuration {
	return config
}

// InitializeConfig loads the application configuration using environment properties
// respecting '.env' file.
func InitializeConfig(release bool) {
	if release {
		os.Setenv(keyProfiles, fmt.Sprintf("%s %s", os.Getenv(keyProfiles), productionProfile))
	}

	godotenv.Load(".env.local")
	godotenv.Load(".env")

	profiles := getProfiles()

	for _, profile := range profiles {
		glg.Infof("Loading configuration for profile '%s'...", profile)
		loadProfile(profile)
	}

	// Initialize and validate configuration singleton
	config = &Configuration{
		profiles:          profiles,
		publicURL:         getPublicURL(),
		port:              getPort(),
		databaseType:      getDatabaseType(),
		databaseURL:       getDatabaseURL(),
		jwtAlgorithm:      getJWTAlgorithm(),
		jwtExpirationTime: getJWTExpirationTime(),
		jwtRefreshTime:    getJWTRefreshTime(),
		dataPath:          getDataPath(),
	}

	config.migrateDatabase = getMigrateDatabase(config)
	config.jwtKey = getJWTKey(config)
}

func loadProfile(profile string) {
	godotenv.Load(".env." + profile + ".local")
	godotenv.Load(".env." + profile)
}

// GetProfiles returns the currently active application profile(s).
func getProfiles() []string {
	profiles := strings.TrimSpace(os.Getenv(keyProfiles))

	// Default to 'development' profile
	if profiles == "" {
		return []string{developmentProfile}
	}

	profilesSlice := strings.Split(profiles, ",")

	for i, profile := range profilesSlice {
		profilesSlice[i] = strings.TrimSpace(profile)
	}

	return profilesSlice
}

func getPublicURL() *url.URL {
	publicURL := strings.TrimSpace(os.Getenv(keyPublicURL))

	if publicURL == "" {
		publicURL = defaultPublicURL
	}

	url, err := url.Parse(publicURL)

	if err != nil {
		glg.Warnf("Invalid base URL given: '%s'", publicURL)
		glg.Warnf("Falling back to default '%s'", defaultPublicURL)
		url, _ = url.Parse(defaultPublicURL)
	}

	return url
}

func getPort() int {
	portString := strings.TrimSpace(os.Getenv(keyPort))

	if portString == "" {
		return defaultPort
	}

	port, err := strconv.Atoi(portString)

	if err != nil {
		glg.Warnf("Invalid port given: '%s'!", portString)
		glg.Warnf("Using default port '%d'", defaultPort)
	}

	return port
}

func getDatabaseType() string {
	dbType := strings.TrimSpace(os.Getenv(keyDatabaseType))

	if dbType == "" {
		return defaultDatabaseType
	}

	return dbType
}

func getDatabaseURL() *url.URL {
	databaseURL := strings.TrimSpace(os.Getenv(keyDatabaseURL))

	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}

	url, err := url.Parse(databaseURL)

	if err != nil {
		glg.Warnf("Invalid database URL given: '%s'", databaseURL)
		glg.Warnf("Falling back to default '%s'", defaultDatabaseURL)
		url, _ = url.Parse(defaultDatabaseURL)
	}

	return url
}

func getMigrateDatabase(c *Configuration) bool {
	migrateDatabaseString := strings.ToLower(strings.TrimSpace(os.Getenv(keyDatabaseMigrate)))

	// Fallback to default -> migrate
	if migrateDatabaseString == "" {
		return true
	}

	if migrateDatabaseString == "true" {
		return true
	} else if migrateDatabaseString != "false" {
		glg.Fatalf(
			"Invalid value '%s' for configuration key '%s' given! Only 'true' and 'false' are allowed.",
			keyDatabaseMigrate,
			migrateDatabaseString,
		)
		panic("Invalid configuration!")
	}

	return false
}

func getJWTAlgorithm() string {
	algorithm := strings.TrimSpace(os.Getenv(keyJWTAlgorithm))

	if algorithm == "" {
		algorithm = defaultJWTAlgorithm
	}

	algorithm = strings.ToUpper(algorithm)

	if !isJWTAlgorithmSupported(algorithm) {
		panic("JWT algorithm '" + algorithm + "' is unknown or unsupported!")
	}

	return algorithm
}

func getJWTKey(c *Configuration) []byte {
	key := strings.TrimSpace(os.Getenv(keyJWTKey))

	if key == "" {
		key = RandomString(32)

		if c.IsDevelopment() {
			glg.Infof("No JWT key given, using random one: %s", key)
		} else {
			glg.Warn("No JWT key given, using random one.")
		}
	}

	return []byte(key)
}

func getJWTExpirationTime() time.Duration {
	expirationTimeString := strings.TrimSpace(os.Getenv(keyJWTExpirationTime))
	var expirationTime time.Duration

	if expirationTimeString == "" {
		return defaultJWTExpirationTime
	}

	expirationTime, err := time.ParseDuration(expirationTimeString)

	if err != nil {
		glg.Warnf("Invalid expiration time given: '%s'!", expirationTimeString)
		glg.Warnf("Using default expiration time '%s'", defaultJWTExpirationTime)

		return defaultJWTExpirationTime
	}

	return expirationTime
}

func getJWTRefreshTime() time.Duration {
	refreshTimeString := strings.TrimSpace(os.Getenv(keyJWTRefreshTime))
	var refreshTime time.Duration

	if refreshTimeString == "" {
		return defaultJWTRefreshTime
	}

	refreshTime, err := time.ParseDuration(refreshTimeString)

	if err != nil {
		glg.Warnf("Invalid expiration time given: '%s'!", refreshTimeString)
		glg.Warnf("Using default expiration time '%s'", defaultJWTRefreshTime)

		return defaultJWTRefreshTime
	}

	return refreshTime
}

func isJWTAlgorithmSupported(algorithm string) bool {
	for _, alg := range jwtSupportedAlgorithms {
		if alg == algorithm {
			return true
		}
	}

	return false
}

func getDataPath() string {
	dataPath := strings.TrimSpace(os.Getenv(keyDataPath))

	if dataPath == "" {
		return defaultDataPath
	}

	return dataPath
}
