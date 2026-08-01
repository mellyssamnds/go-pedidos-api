package controllers

import (
	"bytes"
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

func TestProdutoController_Create_Success(t *testing.T) {
	repo := new(mocks.MockProdutoRepository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("model.Produto")).
		Return(nil)

	produtoService := service.NewProdutoService(repo)
	controller := NewProdutoController(produtoService)

	body := bytes.NewBufferString(`{"name":"Produto Teste","price":10.5}`)
	req := httptest.NewRequest(http.MethodPost, "/produtos", body)
	rec := httptest.NewRecorder()

	controller.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resposta map[string]any
	json.Unmarshal(rec.Body.Bytes(), &resposta)
	assert.Equal(t, "Produto Teste", resposta["name"])
	assert.Equal(t, "10.5", resposta["price"])
}

func TestProdutoController_FindByID_Success(t *testing.T) {
	id := uuid.New()
	repo := new(mocks.MockProdutoRepository)
	repo.On("FindByID", mock.Anything, id).
		Return(model.Produto{ID: id, Name: "Produto Teste", Price: decimal.NewFromFloat(10.5)}, nil)

	produtoService := service.NewProdutoService(repo)
	controller := NewProdutoController(produtoService)

	req := httptest.NewRequest(http.MethodGet, "/produtos/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	rec := httptest.NewRecorder()

	controller.FindByID(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestProdutoController_FindByID_NotFound(t *testing.T) {
	id := uuid.New()
	repo := new(mocks.MockProdutoRepository)
	repo.On("FindByID", mock.Anything, id).
		Return(model.Produto{}, model.ErrProdutoNaoEncontrado)

	produtoService := service.NewProdutoService(repo)
	controller := NewProdutoController(produtoService)

	req := httptest.NewRequest(http.MethodGet, "/produtos/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	rec := httptest.NewRecorder()

	controller.FindByID(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
