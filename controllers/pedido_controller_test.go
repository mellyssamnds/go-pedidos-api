package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/mocks"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPedidoController_FindByID_Success(t *testing.T) {
	id := uuid.New()
	pedidoRepo := new(mocks.MockPedidoRepository)
	pedidoRepo.On("FindByID", mock.Anything, id).
		Return(model.Pedido{
			ID:          id,
			Status:      model.StatusPending,
			TotalAmount: decimal.NewFromFloat(100.0),
		}, nil)

	pedidoService := service.NewPedidoService(pedidoRepo, nil, nil, nil, nil)
	controller := NewPedidoController(pedidoService)

	req := httptest.NewRequest(http.MethodGet, "/pedidos/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	rec := httptest.NewRecorder()

	controller.FindByID(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resposta map[string]any
	json.Unmarshal(rec.Body.Bytes(), &resposta)
	assert.Equal(t, "PENDING", resposta["status"])
}

func TestPedidoController_FindByID_NotFound(t *testing.T) {
	id := uuid.New()
	pedidoRepo := new(mocks.MockPedidoRepository)
	pedidoRepo.On("FindByID", mock.Anything, id).
		Return(model.Pedido{}, model.ErrPedidoNaoEncontrado)

	pedidoService := service.NewPedidoService(pedidoRepo, nil, nil, nil, nil)
	controller := NewPedidoController(pedidoService)

	req := httptest.NewRequest(http.MethodGet, "/pedidos/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	rec := httptest.NewRecorder()

	controller.FindByID(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestPedidoController_FindAll_Success(t *testing.T) {
	pedidoRepo := new(mocks.MockPedidoRepository)
	pedidoRepo.On("FindAll", mock.Anything, 10, 0).
		Return([]model.Pedido{
			{ID: uuid.New(), Status: model.StatusPending, TotalAmount: decimal.NewFromFloat(50.0)},
		}, nil)

	pedidoService := service.NewPedidoService(pedidoRepo, nil, nil, nil, nil)
	controller := NewPedidoController(pedidoService)

	req := httptest.NewRequest(http.MethodGet, "/pedidos", nil)
	rec := httptest.NewRecorder()

	controller.FindAll(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
