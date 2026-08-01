package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/stretchr/testify/mock"
)

// mock ProdutoRepository implementa mock ProdutoRepository
type MockProdutoRepository struct {
	mock.Mock
}

func (m *MockProdutoRepository) Save(ctx context.Context, produto model.Produto) error {
	args := m.Called(ctx, produto)
	return args.Error(0)
}

func (m *MockProdutoRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Produto, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Produto), args.Error(1)
}

func (m *MockProdutoRepository) FindAll(ctx context.Context) ([]model.Produto, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Produto), args.Error(1)
}

func (m *MockProdutoRepository) UpdateStock(ctx context.Context, id uuid.UUID, newQuantity int) error {
	args := m.Called(ctx, id, newQuantity)
	return args.Error(0)
}
