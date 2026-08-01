package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/stretchr/testify/mock"
)

type MockPedidoRepository struct {
	mock.Mock
}

func (m *MockPedidoRepository) Save(ctx context.Context, pedido model.Pedido) error {
	args := m.Called(ctx, pedido)
	return args.Error(0)
}

func (m *MockPedidoRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Pedido, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Pedido), args.Error(1)
}

func (m *MockPedidoRepository) FindAll(ctx context.Context, limit, offset int) ([]model.Pedido, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]model.Pedido), args.Error(1)
}

func (m *MockPedidoRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}
