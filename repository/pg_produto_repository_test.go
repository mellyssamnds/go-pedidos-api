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

func TestPgProdutoRepository_SaveAndFindByID(t *testing.T) {
	pool := setupTestDB(t)
	repo := repository.NewPgProdutoRepository(pool)
	ctx := context.Background()

	produto := model.Produto{
		ID:            uuid.New(),
		Name:          "Notebook Teste",
		Description:   "Descrição de teste",
		Price:         decimal.NewFromFloat(3500.00),
		StockQuantity: 5,
		CreatedAt:     time.Now(),
	}

	err := repo.Save(ctx, produto)
	assert.NoError(t, err)

	encontrado, err := repo.FindByID(ctx, produto.ID)
	assert.NoError(t, err)
	assert.Equal(t, produto.Name, encontrado.Name)
	assert.True(t, produto.Price.Equal(encontrado.Price))
	assert.Equal(t, produto.StockQuantity, encontrado.StockQuantity)

	//limpa o banco após o teste
	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM produtos WHERE id = $1", produto.ID)
	})
}

func TestPgProdutoRepository_FindByID_NaoEncontrado(t *testing.T) {
	pool := setupTestDB(t)
	repo := repository.NewPgProdutoRepository(pool)

	//tenta buscar um produto que não existe
	_, err := repo.FindByID(context.Background(), uuid.New())

	assert.ErrorIs(t, err, model.ErrProdutoNaoEncontrado)
}
