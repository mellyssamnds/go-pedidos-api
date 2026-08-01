//go:build integration

package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/repository"
	"github.com/mellyssamnds/go-pedidos-api/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreatePedido_FluxoCompleto(t *testing.T) {
	pool := setupTestDB(t)
	ctx := context.Background()

	clienteRepo := repository.NewPgClienteRepository(pool)
	produtoRepo := repository.NewPgProdutoRepository(pool)
	pedidoRepo := repository.NewPgPedidoRepository(pool)
	itemPedidoRepo := repository.NewPgItemPedidoRepository(pool)

	//cliente e produto precisam existir antes da efetivação pedido
	cliente := model.Cliente{
		ID:           uuid.New(),
		Name:         "Cliente Teste",
		Email:        "clienteteste@teste.com",
		PasswordHash: "hash-fake",
		CreatedAt:    time.Now(),
	}
	assert.NoError(t, clienteRepo.Save(ctx, cliente))

	produto := model.Produto{
		ID:            uuid.New(),
		Name:          "Produto Teste",
		Description:   "Descrição",
		Price:         decimal.NewFromFloat(100.00),
		StockQuantity: 10,
		CreatedAt:     time.Now(),
	}
	assert.NoError(t, produtoRepo.Save(ctx, produto))

	pedidoService := service.NewPedidoService(pedidoRepo, clienteRepo, produtoRepo, itemPedidoRepo, pool)

	itens := []service.ItemPedidoInput{
		{ProductID: produto.ID, Quantity: 3},
	}

	pedido, err := pedidoService.CreatePedido(ctx, cliente.ID, itens)
	assert.NoError(t, err)
	assert.Equal(t, model.StatusPending, pedido.Status)
	assert.True(t, decimal.NewFromFloat(300.00).Equal(pedido.TotalAmount))

	produtoAtualizado, err := produtoRepo.FindByID(ctx, produto.ID)
	assert.NoError(t, err)
	assert.Equal(t, 7, produtoAtualizado.StockQuantity) //garante que o estoque foi atualizado corretamente

	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM itens_pedido WHERE pedido_id = $1", pedido.ID)
		pool.Exec(context.Background(), "DELETE FROM pedidos WHERE id = $1", pedido.ID)
		pool.Exec(context.Background(), "DELETE FROM produtos WHERE id = $1", produto.ID)
		pool.Exec(context.Background(), "DELETE FROM clientes WHERE id = $1", cliente.ID)
	})
}

func TestCreatePedido_EstoqueInsuficiente_Rollback(t *testing.T) {
	pool := setupTestDB(t)
	ctx := context.Background()

	clienteRepo := repository.NewPgClienteRepository(pool)
	produtoRepo := repository.NewPgProdutoRepository(pool)
	pedidoRepo := repository.NewPgPedidoRepository(pool)
	itemPedidoRepo := repository.NewPgItemPedidoRepository(pool)

	cliente := model.Cliente{
		ID:           uuid.New(),
		Name:         "Cliente Teste 2",
		Email:        "clienteteste2@teste.com",
		PasswordHash: "hash-fake",
		CreatedAt:    time.Now(),
	}
	assert.NoError(t, clienteRepo.Save(ctx, cliente))

	produto := model.Produto{
		ID:            uuid.New(),
		Name:          "Produto Estoque Baixo",
		Description:   "Descrição",
		Price:         decimal.NewFromFloat(50.00),
		StockQuantity: 2, // só 2 em estoque
		CreatedAt:     time.Now(),
	}
	assert.NoError(t, produtoRepo.Save(ctx, produto))

	pedidoService := service.NewPedidoService(pedidoRepo, clienteRepo, produtoRepo, itemPedidoRepo, pool)

	itens := []service.ItemPedidoInput{
		{ProductID: produto.ID, Quantity: 5}, //verifica que não é possível criar pedido com quantidade maior que o estoque
	}
	_, err := pedidoService.CreatePedido(ctx, cliente.ID, itens)
	assert.ErrorIs(t, err, model.ErrEstoqueInsuficiente)

	//confirma que o rollback aconteceu
	//o estoque não deve ter sido alterado, pois a transação falhou
	produtoAposFalha, err := produtoRepo.FindByID(ctx, produto.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2, produtoAposFalha.StockQuantity)

	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM produtos WHERE id = $1", produto.ID)
		pool.Exec(context.Background(), "DELETE FROM clientes WHERE id = $1", cliente.ID)
	})
}
