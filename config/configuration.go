package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	log "github.com/kpango/glg"
)

const (
	envPrefix          = "PAPERLESS"
	profilesKey        = "PROFILES"
	productionProfile  = "production"
	developmentProfile = "development"
)

// Configuration holds all the whole application's configuration.
type Configuration struct {
	Profiles []string `ignored:"true"`
	Server   ServerConfiguration
	Database DatabaseConfiguration
	Security SecurityConfiguration
	Storage  StorageConfiguration
}

// ServerConfiguration holds all configuration values regarding the HTTP server.
type ServerConfiguration struct {
	PublicURL string `default:"http://localhost:8080"`
	Port      int    `default:"8080"`
}

// DatabaseConfiguration holds all configuration values regarding the database.
type DatabaseConfiguration struct {
	Type            string `default:"sqlite3"`
	URL             string `default:"database.db"`
	MigrateDatabase bool   `default:"true"`
}

// SecurityConfiguration holds all configuration values regarding security.
type SecurityConfiguration struct {
	JWTAlgorithm      string        `default:"HS256" split_words:"true"`
	JWTSecret         []byte        `split_words:"true"`
	JWTExpirationTime time.Duration `default:"5m" split_words:"true"`
	JWTRefreshTime    time.Duration `default:"24h" split_words:"true"`
}

// StorageConfiguration holds all configuration values regarding file storage.
type StorageConfiguration struct {
	DataPath string `default:"data" split_words:"true"`
}

// HasProfile returns a boolean value indicating whether the given profile is active.
func (c Configuration) HasProfile(profile string) bool {
	for _, configuredProfile := range c.Profiles {
		if configuredProfile == profile {
			return true
		}
	}

	return false
}

// IsProductionMode returns a boolean value indicating whether the applications
// runs in production mode (having active profile 'production').
func (c Configuration) IsProductionMode() bool {
	return c.HasProfile(productionProfile)
}

// Load loads the application configuration.
func Load(production bool) *Configuration {
	// Set 'production' profile based on production flag.
	if production {
		os.Setenv(
			profilesKey,
			fmt.Sprintf("%s %s", productionProfile, os.Getenv(profilesKey)),
		)
	}

	cfg := &Configuration{}
	cfg.Profiles = getProfiles()
	loadProfiles(cfg)

	if err := envconfig.Process(envPrefix, cfg); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return cfg
}

func loadProfiles(cfg *Configuration) {
	godotenv.Load(".env.local")
	godotenv.Load(".env")

	for _, profile := range cfg.Profiles {
		log.Infof("Loading configuration for profile '%s'...", profile)
		prefix := fmt.Sprintf(".env.%s", profile)
		godotenv.Load(prefix, prefix+".local")
	}
}

func getProfiles() []string {
	profiles := strings.TrimSpace(os.Getenv(profilesKey))

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
