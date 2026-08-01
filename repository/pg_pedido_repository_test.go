//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/repository"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPgPedidoRepository_UpdateStatus(t *testing.T) {
	pool := setupTestDB(t)
	ctx := context.Background()

	clienteRepo := repository.NewPgClienteRepository(pool)
	pedidoRepo := repository.NewPgPedidoRepository(pool)

	cliente := model.Cliente{
		ID:           uuid.New(),
		Name:         "Cliente Pedido Teste",
		Email:        "clienteteste@teste.com",
		PasswordHash: "hash-fake",
		CreatedAt:    time.Now(),
	}
	assert.NoError(t, clienteRepo.Save(ctx, cliente))

	pedido := model.Pedido{
		ID:          uuid.New(),
		ClienteID:   cliente.ID,
		CreatedAt:   time.Now(),
		Status:      model.StatusPending,
		TotalAmount: decimal.NewFromFloat(150.0),
	}
	assert.NoError(t, pedidoRepo.Save(ctx, pedido))

	err := pedidoRepo.UpdateStatus(ctx, pedido.ID, model.StatusPaid)
	assert.NoError(t, err)

	atualizado, err := pedidoRepo.FindByID(ctx, pedido.ID)
	assert.NoError(t, err)
	assert.Equal(t, model.StatusPaid, atualizado.Status)

	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM pedidos WHERE id = $1", pedido.ID)
		pool.Exec(context.Background(), "DELETE FROM clientes WHERE id = $1", cliente.ID)
	})
}

func TestPgPedidoRepository_UpdateStatus_NaoEncontrado(t *testing.T) {
	pool := setupTestDB(t)
	repo := repository.NewPgPedidoRepository(pool)

	err := repo.UpdateStatus(context.Background(), uuid.New(), model.StatusPaid)

	assert.ErrorIs(t, err, model.ErrPedidoNaoEncontrado)
}
