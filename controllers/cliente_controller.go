package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/service"
)

type ClienteController struct {
	service *service.ClienteService
}

type createClienteRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewClienteController(service *service.ClienteService) *ClienteController {
	return &ClienteController{service: service}
}

func (c *ClienteController) Create(w http.ResponseWriter, r *http.Request) {
	var req createClienteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	cliente, err := c.service.CreateCliente(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, model.ErrEmailCadastrado) {
			writeError(w, err.Error(), http.StatusConflict)
			return
		}
		writeError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	writeJSON(w, cliente, http.StatusCreated)
}

func (c *ClienteController) FindByID(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID do cliente da URL
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		writeError(w, "id inválido", http.StatusBadRequest)
		return
	}

	cliente, err := c.service.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrClienteNaoEncontrado) {
			writeError(w, err.Error(), http.StatusNotFound)
			return
		}
		writeError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	writeJSON(w, cliente, http.StatusOK)
}

func (c *ClienteController) FindAll(w http.ResponseWriter, r *http.Request) {
	clientes, err := c.service.FindAll(r.Context())

	if err != nil {
		writeError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	writeJSON(w, clientes, http.StatusOK)
}
