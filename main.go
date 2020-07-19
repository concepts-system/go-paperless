package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/concepts-system/go-paperless/domain"

	"github.com/concepts-system/go-paperless/infrastructure"

	"github.com/concepts-system/go-paperless/application"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	log "github.com/kpango/glg"

	"github.com/concepts-system/go-paperless/config"
	"github.com/concepts-system/go-paperless/web"
)

var (
	version   string
	buildDate string
	release   string
)

type bootstrapper struct {
	config   *config.Configuration
	database *infrastructure.Database
	server   *web.Server

	users     domain.Users
	documents domain.Documents

	authService     application.AuthService
	userService     application.UserService
	documentService application.DocumentService

	tokenKeyResolver application.TokenKeyResolver
}

func main() {
	start := time.Now()
	bs := &bootstrapper{}

	if version == "" {
		version = "DEV-SNAPSHOT"
	}

	if buildDate == "" {
		buildDate = time.Now().Format("2006-01-02")
	}

	rand.Seed(time.Now().UnixNano())
	log.Infof("Starting application %s (%s)", version, buildDate)

	loadConfiguration(bs)
	prepareDatabase(bs)
	defer bs.database.Close()

	setupDependencies(bs)
	initializeServer(bs)
	ensureUserExists(bs)

	log.Successf("Start-up completed in %v", time.Since(start))
	bs.server.Start()
}

func loadConfiguration(bs *bootstrapper) {
	log.Info("Loading configuration...")
	bs.config = config.LoadConfiguration(release == "true")
	createDirectories(bs)
}

func prepareDatabase(bs *bootstrapper) {
	bs.database = infrastructure.NewDatabase(bs.config)
	bs.database.Connect()

	if bs.config.MigrateDatabase() {
		log.Info("Running migrations...")

		// Migrate to most recent version by default.
		// Use 'bs.database.MigrateTo(version)' to migrate to a specific version.
		if err := bs.database.Migrate(); err != nil {
			log.Fatalf("Error while migrating database: %v", err)
		}

	}
}

func createDirectories(bs *bootstrapper) {
	log.Info("Setting up directories...")
	config := bs.config

	// Create data directory
	if _, err := os.Stat(config.GetDataPath()); os.IsNotExist(err) {
		os.MkdirAll(config.GetDataPath(), os.ModePerm)
	}
}

func setupDependencies(bs *bootstrapper) {
	bs.tokenKeyResolver = application.ConfigTokenKeyResolver(bs.config)
	bs.users = infrastructure.NewUsers(bs.database)
	bs.documents = infrastructure.NewDocuments(bs.database)

	bs.userService = application.NewUserService(bs.users)
	bs.authService = application.NewAuthService(
		bs.config,
		bs.users,
		bs.tokenKeyResolver,
	)
	bs.documentService = application.NewDocumentService(bs.users, bs.documents)
}

func initializeServer(bs *bootstrapper) {
	bs.server = web.NewServer(bs.config, bs.authService)
	registerRouters(bs)
}

func registerRouters(bs *bootstrapper) {
	bs.server.Register(
		// Auth routes
		web.NewAuthRouter(
			bs.authService,
			application.ConfigTokenKeyResolver(bs.config),
		),

		// User routes
		web.NewUserRouter(bs.userService),

		// Document routes
		web.NewDocumentRouter(bs.documentService),
	)
}

// func initializeWorkers() {
// 	log.Info("Initializing workers...")

// 	manager := worker.NewManager()
// 	documents.RegisterWorkers(manager)

// 	go manager.Run()
// }

func ensureUserExists(bs *bootstrapper) {
	_, count, err := bs.userService.GetUsers(domain.PageRequest{Offset: 0, Size: 1})

	if err != nil {
		log.Fatalf("Error while checking for default user: %v", err)
		panic(err)
	}

	if count > 0 {
		return
	}

	log.Info("No user present; creating default one...")
	defaultUser := domain.NewUser(domain.User{
		Username: "admin",
		Forename: "Default",
		Surname:  "User",
		IsAdmin:  true,
		IsActive: true,
	})

	defaultUser, err = bs.userService.CreateNewUser(defaultUser, "admin")

	if err != nil {
		log.Fatalf("Error while creating default user: %v", err)
	}
}
