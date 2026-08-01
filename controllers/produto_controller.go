package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/service"
	"github.com/shopspring/decimal"
)

type ProdutoController struct {
	service *service.ProdutoService
}

type createProdutoRequest struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stockQuantity"`
}

func NewProdutoController(service *service.ProdutoService) *ProdutoController {
	return &ProdutoController{service: service}
}

func (c *ProdutoController) Create(w http.ResponseWriter, r *http.Request) {
	var req createProdutoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	// Converte o preço de float64 para decimal.Decimal
	price := decimal.NewFromFloat(req.Price)

	produto, err := c.service.CreateProduto(r.Context(), req.Name, req.Description, price, req.StockQuantity)
	if err != nil {
		writeError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	writeJSON(w, produto, http.StatusCreated)
}

func (c *ProdutoController) FindByID(w http.ResponseWriter, r *http.Request) {
	// Extrai o ID do produto da URL
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		writeError(w, "id inválido", http.StatusBadRequest)
		return
	}

	produto, err := c.service.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrProdutoNaoEncontrado) {
			writeError(w, err.Error(), http.StatusNotFound)
			return
		}
		writeError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	writeJSON(w, produto, http.StatusOK)
}

func (c *ProdutoController) FindAll(w http.ResponseWriter, r *http.Request) {
	produtos, err := c.service.FindAll(r.Context())

	if err != nil {
		writeError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	writeJSON(w, produtos, http.StatusOK)
}
