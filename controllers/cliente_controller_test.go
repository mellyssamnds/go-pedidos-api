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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClienteController_Create_Success(t *testing.T) {
	repo := new(mocks.MockClienteRepository)
	repo.On("FindByEmail", mock.Anything, "anaju@teste.com").
		Return(model.Cliente{}, model.ErrClienteNaoEncontrado)
	repo.On("Save", mock.Anything, mock.AnythingOfType("model.Cliente")).
		Return(nil)

	clienteService := service.NewClienteService(repo)
	controller := NewClienteController(clienteService)

	body := bytes.NewBufferString(`{"name":"Ana Julia","email":"anaju@teste.com","password":"senha*123"}`)
	req := httptest.NewRequest(http.MethodPost, "/clientes", body)
	rec := httptest.NewRecorder()

	controller.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resposta map[string]any
	json.Unmarshal(rec.Body.Bytes(), &resposta)
	assert.Equal(t, "Ana Julia", resposta["name"])
}

func TestClienteController_Create_EmailAlreadyExists(t *testing.T) {
	repo := new(mocks.MockClienteRepository)
	repo.On("FindByEmail", mock.Anything, "anaju@teste.com").
		Return(model.Cliente{ID: uuid.New(), Email: "anaju@teste.com"}, nil)

	clienteService := service.NewClienteService(repo)
	controller := NewClienteController(clienteService)

	body := bytes.NewBufferString(`{"name":"Ana Julia","email":"anaju@teste.com","password":"senha*123"}`)
	req := httptest.NewRequest(http.MethodPost, "/clientes", body)
	rec := httptest.NewRecorder()

	controller.Create(rec, req)

	assert.Equal(t, http.StatusConflict, rec.Code)
}

func TestClienteController_FindByID_Success(t *testing.T) {
	id := uuid.New()
	repo := new(mocks.MockClienteRepository)
	repo.On("FindByID", mock.Anything, id).
		Return(model.Cliente{ID: id, Name: "Ana Julia", Email: "anaju@teste.com"}, nil)

	clienteService := service.NewClienteService(repo)
	controller := NewClienteController(clienteService)

	req := httptest.NewRequest(http.MethodGet, "/clientes/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	rec := httptest.NewRecorder()

	controller.FindByID(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resposta map[string]any
	json.Unmarshal(rec.Body.Bytes(), &resposta)
	assert.Equal(t, id.String(), resposta["id"])
	assert.Equal(t, "Ana Julia", resposta["name"])
	assert.Equal(t, "anaju@teste.com", resposta["email"])
}

func TestClienteController_FindByID_NotFound(t *testing.T) {
	id := uuid.New()
	repo := new(mocks.MockClienteRepository)
	repo.On("FindByID", mock.Anything, id).
		Return(model.Cliente{}, model.ErrClienteNaoEncontrado)

	clienteService := service.NewClienteService(repo)
	controller := NewClienteController(clienteService)

	req := httptest.NewRequest(http.MethodGet, "/clientes/"+id.String(), nil)
	req.SetPathValue("id", id.String())
	rec := httptest.NewRecorder()

	controller.FindByID(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
