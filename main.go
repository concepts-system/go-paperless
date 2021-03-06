package main

import (
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/domain"

	"github.com/concepts-system/go-paperless/infrastructure"

	"github.com/concepts-system/go-paperless/application"

	"github.com/concepts-system/go-paperless/config"
	"github.com/concepts-system/go-paperless/web"
)

var (
	version   string
	buildDate string
	release   string
)

const (
	documentsDirectory = "documents"
)

var log = common.NewLogger("main")

type bootstrapper struct {
	config   *config.Configuration
	database *infrastructure.Database
	server   *web.Server

	tubeMail domain.TubeMail

	users                domain.Users
	documents            domain.Documents
	documentArchive      domain.DocumentArchive
	documentPreprocessor domain.DocumentPreprocessor
	documentAnalyzer     domain.DocumentAnalyzer
	documentIndex        domain.DocumentIndex
	documentRegistry     domain.DocumentRegistry

	authService     application.AuthService
	userService     application.UserService
	documentService application.DocumentService

	tokenKeyResolver application.TokenKeyResolver
}

func main() {
	start := time.Now()
	bs := &bootstrapper{}

	loadConfiguration(bs)

	if version == "" {
		version = "DEV-SNAPSHOT"
	}

	if buildDate == "" {
		buildDate = time.Now().Format("2006-01-02")
	}

	rand.Seed(time.Now().UnixNano())
	log.Infof("Starting application %s (%s)", version, buildDate)

	prepareDatabase(bs)
	db, err := bs.database.DB.DB()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	setupDependencies(bs)
	initializeServer(bs)
	ensureUserExists(bs)

	log.Infof("Bootstrap completed in %v", time.Since(start))

	if err := bs.server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func loadConfiguration(bs *bootstrapper) {
	bs.config = config.Load(release == "true")
	config.ConfigureLogging(bs.config)
	log.Info("Configuration loaded")
	createDirectories(bs)
}

func prepareDatabase(bs *bootstrapper) {
	bs.database = infrastructure.NewDatabase(bs.config)
	if err := bs.database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if bs.config.Database.MigrateDatabase {
		log.Info("Running migrations...")

		// Migrate to most recent version by default.
		// Use 'bs.database.MigrateTo(version)' to migrate to a specific version.
		if err := bs.database.Migrate(); err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
	}
}

func createDirectories(bs *bootstrapper) {
	log.Info("Setting up directories...")
	config := bs.config

	// Create data directory
	if _, err := os.Stat(config.Storage.DataPath); os.IsNotExist(err) {
		err := os.MkdirAll(config.Storage.DataPath, os.ModePerm)

		if err != nil {
			log.Errorf("Failed to setup data directory: %v", err)
		}
	}
}

func setupDependencies(bs *bootstrapper) {
	bs.tokenKeyResolver = application.ConfigTokenKeyResolver(bs.config)
	bs.tubeMail = infrastructure.NewLocalAsyncTubeMailImpl()
	bs.users = infrastructure.NewUsers(bs.database)
	bs.documents = infrastructure.NewDocuments(bs.database)
	initializeDocumentArchive(bs)
	bs.documentPreprocessor = infrastructure.NewDocumentPreprocessorImpl(bs.documents, bs.documentArchive)
	bs.documentAnalyzer = infrastructure.NewTesseractOcrEngine(bs.documents, bs.documentArchive)
	initializeDocumentIndex(bs)

	bs.documentRegistry = domain.NewDocumentRegistry(
		bs.tubeMail,
		bs.documents,
		bs.documentPreprocessor,
		bs.documentAnalyzer,
		bs.documentIndex,
	)

	bs.userService = application.NewUserService(bs.users)
	bs.authService = application.NewAuthService(
		bs.config,
		bs.users,
		bs.tokenKeyResolver,
	)

	bs.documentService = application.NewDocumentService(
		bs.users,
		bs.documents,
		bs.documentArchive,
		bs.documentIndex,
		bs.documentRegistry,
	)
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

func initializeDocumentArchive(bs *bootstrapper) {
	documentArchive, err := infrastructure.NewDocumentArchiveFileSystemImpl(
		path.Join(bs.config.Storage.DataPath, documentsDirectory),
	)

	if err != nil {
		log.Fatalf("Failed to initialize document archive: %v", err)
	}

	bs.documentArchive = documentArchive
}

func initializeDocumentIndex(bs *bootstrapper) {
	documentIndex, err := infrastructure.NewBleveDocumentIndex(bs.config.Index.DocumentsPath, bs.documents)

	if err != nil {
		log.Fatalf("Failed to initialize document index: %v", err)
	}

	bs.documentIndex = documentIndex
}

func ensureUserExists(bs *bootstrapper) {
	_, count, err := bs.userService.GetUsers(domain.PageRequest{Offset: 0, Size: 1})

	if err != nil {
		log.Fatalf("Failed checking for default user: %v", err)
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

	_, err = bs.userService.CreateNewUser(defaultUser, "admin")

	if err != nil {
		log.Fatal("Failed to create default user", err)
	}
}
