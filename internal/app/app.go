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

func NewApp() *App {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := databases.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	repo := mongo.NewMongoPersonRepo(database)
	svc := service.NewPersonService(repo)

	h := handlers.NewHandler(svc)

	r := routes.SetupRouter(h)

	return &App{Router: r, cfg: cfg}
}

func (a *App) Run() {
	addr := ":" + a.cfg.Port
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, a.Router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
