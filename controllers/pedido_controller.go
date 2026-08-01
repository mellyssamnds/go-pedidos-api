package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/service"
)

type PedidoController struct {
	pedidoService *service.PedidoService
}

func NewPedidoController(pedidoService *service.PedidoService) *PedidoController {
	return &PedidoController{
		pedidoService: pedidoService,
	}
}

type itemPedidoRequest struct {
	ProdutoID uuid.UUID `json:"produtoId"`
	Quantity  int       `json:"quantity"`
}

type createPedidoRequest struct {
	ClienteID uuid.UUID           `json:"clienteId"`
	Itens     []itemPedidoRequest `json:"itens"`
}

func (pc *PedidoController) Create(w http.ResponseWriter, r *http.Request) {
	var req createPedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	itens := make([]service.ItemPedidoInput, len(req.Itens))
	for i, item := range req.Itens {
		itens[i] = service.ItemPedidoInput{
			ProductID: item.ProdutoID,
			Quantity:  item.Quantity,
		}
	}

	//chama o serviço para criar o pedido
	pedido, err := pc.pedidoService.CreatePedido(r.Context(), req.ClienteID, itens)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrClienteNaoEncontrado),
			errors.Is(err, model.ErrProdutoNaoEncontrado):
			writeError(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, model.ErrEstoqueInsuficiente):
			writeError(w, err.Error(), http.StatusConflict)
		default:
			writeError(w, "erro interno", http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, pedido, http.StatusCreated)
}

func (pc *PedidoController) FindAll(w http.ResponseWriter, r *http.Request) {
	//obtém os parâmetros de paginação da query string
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	//obtém o parâmetro de offset da query string
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}

	//chama o serviço para buscar todos os pedidos com paginação
	pedidos, err := pc.pedidoService.FindAll(r.Context(), limit, offset)
	if err != nil {
		writeError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	writeJSON(w, pedidos, http.StatusOK)
}

func (pc *PedidoController) FindByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	//converte o id de string para uuid
	id, err := uuid.Parse(idStr)
	if err != nil {
		writeError(w, "id inválido", http.StatusBadRequest)
		return
	}

	pedido, err := pc.pedidoService.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrPedidoNaoEncontrado) {
			writeError(w, err.Error(), http.StatusNotFound)
			return
		}
		writeError(w, "erro interno", http.StatusInternalServerError)
		return
	}

	writeJSON(w, pedido, http.StatusOK)
}
