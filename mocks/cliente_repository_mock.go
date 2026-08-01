package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/stretchr/testify/mock"
)

type MockClienteRepository struct {
	mock.Mock
}

func (m *MockClienteRepository) Save(ctx context.Context, cliente model.Cliente) error {
	args := m.Called(ctx, cliente)
	return args.Error(0)
}

func (m *MockClienteRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Cliente, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Cliente), args.Error(1)
}

func (m *MockClienteRepository) FindByEmail(ctx context.Context, email string) (model.Cliente, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.Cliente), args.Error(1)
}

func (m *MockClienteRepository) FindAll(ctx context.Context) ([]model.Cliente, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Cliente), args.Error(1)
}
