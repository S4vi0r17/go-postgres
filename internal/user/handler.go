package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/savior/go-postgres/internal/platform/response"
	"github.com/savior/go-postgres/internal/platform/validation"
)

// Handler es el equivalente al Controller de NestJS.
// Recibe requests HTTP, llama al service, escribe la respuesta.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes monta las rutas de este feature en el router recibido.
// Así el paquete server no tiene que saber qué rutas tiene cada feature.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/users", h.List)
	r.Get("/users/{id}", h.Get)
	r.Post("/users", h.Create)
	r.Put("/users/{id}", h.Update)
	r.Delete("/users/{id}", h.Delete)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "Error fetching users")
		return
	}
	response.OKList(w, users, len(users))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	u, err := h.svc.Get(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}
	response.OK(w, u)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var dto CreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if errs := validation.Struct(dto); errs != nil {
		response.Validation(w, errs)
		return
	}

	u, err := h.svc.Create(r.Context(), dto)
	if err != nil {
		h.handleError(w, err)
		return
	}
	response.Created(w, "User created successfully", u)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	var dto UpdateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	if errs := validation.Struct(dto); errs != nil {
		response.Validation(w, errs)
		return
	}

	u, err := h.svc.Update(r.Context(), id, dto)
	if err != nil {
		h.handleError(w, err)
		return
	}
	response.OKMessage(w, "User updated successfully", u)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.handleError(w, err)
		return
	}
	response.OKMessage(w, "User deleted successfully", nil)
}

// handleError traduce errores de dominio a respuestas HTTP.
func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		response.Error(w, http.StatusNotFound, "User not found")
	case errors.Is(err, ErrEmailExists):
		response.Error(w, http.StatusConflict, "Email already exists")
	default:
		response.Error(w, http.StatusInternalServerError, "Internal error: "+err.Error())
	}
}
