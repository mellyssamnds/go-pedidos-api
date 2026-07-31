package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/repository"
	"github.com/shopspring/decimal"
)

type ProdutoService struct {
	produtoRepo repository.ProdutoRepository
}

func NewProdutoService(produtoRepo repository.ProdutoRepository) *ProdutoService {
	return &ProdutoService{
		produtoRepo: produtoRepo,
	}
}

func (s *ProdutoService) CreateProduto(ctx context.Context, name, description string, price decimal.Decimal, stockQuantity int) (model.Produto, error) {
	produto := model.Produto{
		ID:            uuid.New(),
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		CreatedAt:     time.Now(),
	}

	if err := s.produtoRepo.Save(ctx, produto); err != nil {
		return model.Produto{}, err
	}

	return produto, nil
}

func (s *ProdutoService) FindByID(ctx context.Context, id uuid.UUID) (model.Produto, error) {
	return s.produtoRepo.FindByID(ctx, id)
}

func (s *ProdutoService) FindAll(ctx context.Context) ([]model.Produto, error) {
	return s.produtoRepo.FindAll(ctx)
}
