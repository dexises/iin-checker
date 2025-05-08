package handlers

import (
	"encoding/json"
	"github.com/dexises/iin-checker/internal/models"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/dexises/iin-checker/internal/service"
)

type Handler struct {
	svc service.PersonService
}

func NewHandler(svc service.PersonService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CheckIINHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CheckIINRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, models.CheckIINResponse{Valid: false, Error: "invalid JSON payload"})
		return
	}

	date, gender, ok, err := h.svc.ValidateIIN(req.IIN)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, models.CheckIINResponse{Valid: false, Error: err.Error()})
		return
	}
	if !ok {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, models.CheckIINResponse{Valid: false, Error: "invalid IIN"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, models.CheckIINResponse{
		Valid:  true,
		Date:   date.Format("2006-01-02"),
		Gender: gender,
	})
}

func (h *Handler) CreatePersonHandler(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePersonRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, models.CreatePersonResponse{Success: false, Error: "invalid JSON payload"})
		return
	}

	_, err := h.svc.Create(r.Context(), req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, models.CreatePersonResponse{Success: false, Error: err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, models.CreatePersonResponse{Success: true})
}

func (h *Handler) SearchPersonHandler(w http.ResponseWriter, r *http.Request) {
	raw := chi.URLParam(r, "name")
	part, err := url.PathUnescape(raw)
	if err != nil {
		part = raw
	}
	log.Printf("SearchPersonHandler: looking for name part %q", part)
	persons, err := h.svc.FindByName(r.Context(), part)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, persons)
}

func (h *Handler) GetPersonHandler(w http.ResponseWriter, r *http.Request) {
	iin := chi.URLParam(r, "iin")
	person, err := h.svc.Get(r.Context(), iin)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, models.CreatePersonResponse{Error: "person not found"})
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, person)
}
