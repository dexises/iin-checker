package routes

import (
	handlers "github.com/dexises/iin-checker/internal/handler"
	"github.com/go-chi/chi/v5"
)

func SetupRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	// IIN validation endpoint
	r.Post("/iin_check", h.CheckIINHandler)

	// Person endpoints
	r.Post("/people/info", h.CreatePersonHandler)
	r.Get("/people/info/phone/{name}", h.SearchPersonHandler)
	r.Get("/people/info/{iin}", h.GetPersonHandler)

	return r
}
