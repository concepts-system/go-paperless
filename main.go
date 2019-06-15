package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	worker "github.com/contribsys/faktory_worker_go"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/kpango/glg"

	"github.com/concepts-system/go-paperless/api"
	"github.com/concepts-system/go-paperless/auth"
	"github.com/concepts-system/go-paperless/common"
	"github.com/concepts-system/go-paperless/database"
	"github.com/concepts-system/go-paperless/documents"
	"github.com/concepts-system/go-paperless/migrations"
	"github.com/concepts-system/go-paperless/users"
	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var (
	version   string
	buildDate string
	release   string
)

func main() {
	start := time.Now()

	if version == "" {
		version = "DEV-SNAPSHOT"
	}
	if buildDate == "" {
		buildDate = time.Now().Format("2006-01-02")
	}

	rand.Seed(time.Now().UnixNano())

	glg.Infof("Starting application %s (%s)", version, buildDate)
	common.InitializeConfig(release == "true")
	createDirectories()

	defer database.DB().Close()

	if common.Config().MigrateDatabase() {
		glg.Info("Running migrations...")
		migrate()
	}

	glg.Info("Initializing job workers...")
	initializeWorkers()

	glg.Info("Preparing document index...")
	documents.PrepareIndex()
	defer documents.GetIndex().Close()

	ensureUserExists()

	server := initializeServer()
	glg.Successf("Start-up completed in %v", time.Since(start))

	endpoint := fmt.Sprintf(":%d", common.Config().GetPort())
	glg.Successf("Accepting connection on %s", endpoint)
	server.Start(endpoint)
}

func createDirectories() {
	glg.Info("Setting up directories...")

	// Create data directory
	if _, err := os.Stat(common.Config().GetDataPath()); os.IsNotExist(err) {
		os.MkdirAll(common.Config().GetDataPath(), os.ModePerm)
	}
}

func initializeServer() *echo.Echo {
	glg.Info("Initializing server...")

	r := echo.New()
	r.Debug = common.Config().IsDevelopment()
	r.HideBanner = true
	r.HidePort = true
	r.HTTPErrorHandler = api.ErrorHandler
	r.Validator = common.Validator{Validator: validator.New()}

	r.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} \t${method}\t${uri} -> status=${status} [${latency_human}] | ${error}\n",
	}))
	r.Use(middleware.Recover())
	r.Use(api.CustomContext)

	registerRoutes(r.Group(""))
	return r
}

func registerRoutes(r *echo.Group) {
	// Common routes
	auth.RegisterRoutes(r)

	// API routes
	api := r.Group("/api")
	users.RegisterRoutes(api)
	documents.RegisterRoutes(api)
}

func migrate() {
	migrator := migrations.BuildMigrator(database.DB())

	// Migrate to latest version by default
	if err := migrator.Migrate(); err != nil {
		glg.Fatalf("Error while executing DB migrations: %v", err)
		panic("Failed to execute database migrations!")
	}

	// Use 'migrator.MigrateTo(version)' to migrate to a specific version
}

func initializeWorkers() {
	glg.Info("Initializing workers...")

	manager := worker.NewManager()
	documents.RegisterWorkers(manager)

	go manager.Run()
}

func ensureUserExists() {
	_, count, err := users.Find(common.PageRequest{Offset: 0, Size: 1})

	if err != nil {
		glg.Fatalf("Error while checking for default user: %v", err)
		panic(err)
	}

	if count > 0 {
		return
	}

	glg.Info("No user present; creating default one...")
	defaultUser := users.UserModel{
		Username: "admin",
		Forename: "Default",
		Surname:  "User",
		IsAdmin:  true,
		IsActive: true,
	}

	defaultUser.SetPassword("admin")
	err = defaultUser.Create()

	if err != nil {
		glg.Fatalf("Error while creating default user: %v", err)
		panic(err)
	}
}
