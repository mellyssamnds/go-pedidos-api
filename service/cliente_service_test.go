package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/mocks"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCliente_Success(t *testing.T) {
	repo := new(mocks.MockClienteRepository)
	repo.On("FindByEmail", mock.Anything, "anaju@teste.com").
		Return(model.Cliente{}, model.ErrClienteNaoEncontrado)
	repo.On("Save", mock.Anything, mock.AnythingOfType("model.Cliente")).
		Return(nil)

	service := NewClienteService(repo)

	cliente, err := service.CreateCliente(context.Background(), "Ana Julia", "anaju@teste.com", "senha*123")

	assert.NoError(t, err)
	assert.Equal(t, "Ana Julia", cliente.Name)
	assert.Equal(t, "anaju@teste.com", cliente.Email)
	assert.NotEqual(t, uuid.Nil, cliente.ID)
	assert.NotEmpty(t, cliente.PasswordHash)
	assert.NotEqual(t, "senha*123", cliente.PasswordHash) // garante que não salvou em texto puro

	repo.AssertExpectations(t)
}

func TestCreateCliente_EmailJaCadastrado(t *testing.T) {
	clienteExistente := model.Cliente{ID: uuid.New(), Email: "anaju@teste.com"}

	repo := new(mocks.MockClienteRepository)
	repo.On("FindByEmail", mock.Anything, "anaju@teste.com").
		Return(clienteExistente, nil)

	service := NewClienteService(repo)

	_, err := service.CreateCliente(context.Background(), "Ana Julia", "anaju@teste.com", "senha*123")

	assert.ErrorIs(t, err, model.ErrEmailCadastrado)
	repo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
}
