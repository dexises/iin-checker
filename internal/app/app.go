package app

import (
	"github.com/dexises/iin-checker/internal/config"
	"github.com/dexises/iin-checker/internal/databases"
	"github.com/dexises/iin-checker/internal/databases/drivers/mongo"
	handlers "github.com/dexises/iin-checker/internal/handler"
	"github.com/dexises/iin-checker/internal/routes"
	"github.com/dexises/iin-checker/internal/service"
	"log"
	"net/http"
)

type App struct {
	Router http.Handler
	cfg    *config.Config
}

// NewApp initializes configuration, database connection, services, handlers and routes.
func NewApp() *App {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Debug: print MongoDB connection settings
	log.Printf("🔍 MongoDB URL: %s, Database: %s", cfg.Database.DSURL, cfg.Database.DSDB)

	// Connect to MongoDB
	database, err := databases.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	// Initialize repository and service
	repo := mongo.NewMongoPersonRepo(database)
	svc := service.NewPersonService(repo)

	// Initialize HTTP handlers
	h := handlers.NewHandler(svc)

	// Setup routes
	r := routes.SetupRouter(h)

	return &App{Router: r, cfg: cfg}
}

// Run starts the HTTP server.
func (a *App) Run() {
	addr := ":" + a.cfg.Port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, a.Router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
