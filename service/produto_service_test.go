package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduto_Success(t *testing.T) {
	repo := new(MockProdutoRepository)
	repo.On("Save", mock.Anything, mock.AnythingOfType("model.Produto")).
		Return(nil)

	service := NewProdutoService(repo)

	price := decimal.NewFromFloat(10.0)
	produto, err := service.CreateProduto(context.Background(), "Produto Teste", "Descrição de teste", price, 100)

	assert.NoError(t, err)
	assert.Equal(t, "Produto Teste", produto.Name)
	assert.True(t, price.Equal(produto.Price))
	assert.NotEqual(t, uuid.Nil, produto.ID)

	repo.AssertExpectations(t)
}
