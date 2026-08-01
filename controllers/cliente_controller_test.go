package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mellyssamnds/go-pedidos-api/mocks"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClienteController_Create_Success(t *testing.T) {
	repo := new(mocks.MockClienteRepository)
	repo.On("FindByEmail", mock.Anything, "ana@teste.com").
		Return(model.Cliente{}, model.ErrClienteNaoEncontrado)
	repo.On("Save", mock.Anything, mock.AnythingOfType("model.Cliente")).
		Return(nil)

	clienteService := service.NewClienteService(repo)
	controller := NewClienteController(clienteService)

	body := bytes.NewBufferString(`{"name":"Ana Silva","email":"ana@teste.com","password":"senha123"}`)
	req := httptest.NewRequest(http.MethodPost, "/clientes", body)
	rec := httptest.NewRecorder()

	controller.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resposta map[string]any
	json.Unmarshal(rec.Body.Bytes(), &resposta)
	assert.Equal(t, "Ana Silva", resposta["name"])
}
