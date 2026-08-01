package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/repository"
	"github.com/shopspring/decimal"
)

type ItemPedidoInput struct {
	ProductID uuid.UUID
	Quantity  int
}

type PedidoService struct {
	pedidoRepo     repository.PedidoRepository
	clienteRepo    repository.ClienteRepository
	produtoRepo    repository.ProdutoRepository
	itemPedidoRepo repository.ItemPedidoRepository
	pool           *pgxpool.Pool
}

func NewPedidoService(
	pedidoRepo repository.PedidoRepository,
	clienteRepo repository.ClienteRepository,
	produtoRepo repository.ProdutoRepository,
	itemPedidoRepo repository.ItemPedidoRepository,
	pool *pgxpool.Pool,
) *PedidoService {
	return &PedidoService{
		pedidoRepo:     pedidoRepo,
		clienteRepo:    clienteRepo,
		produtoRepo:    produtoRepo,
		itemPedidoRepo: itemPedidoRepo,
		pool:           pool,
	}
}

func (s *PedidoService) CreatePedido(ctx context.Context, clienteID uuid.UUID, itens []ItemPedidoInput) (model.Pedido, error) {
	//valida itens do pedido
	if len(itens) == 0 {
		return model.Pedido{}, errors.New("pedido precisa ter pelo menos um item")
	}

	//valida se o cliente existe
	if _, err := s.clienteRepo.FindByID(ctx, clienteID); err != nil {
		return model.Pedido{}, err
	}

	//inicia transação
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return model.Pedido{}, fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer tx.Rollback(ctx)

	//repositories com as dependências do contexto da transação
	pedidoRepoTx := repository.NewPgPedidoRepository(tx)
	produtoRepoTx := repository.NewPgProdutoRepository(tx)
	itemPedidoRepoTx := repository.NewPgItemPedidoRepository(tx)

	pedido := model.Pedido{
		ID:        uuid.New(),
		ClienteID: clienteID,
		CreatedAt: time.Now(),
		Status:    model.StatusPending,
	}

	total := decimal.NewFromInt(0)
	var itensParaSalvar []model.ItemPedido

	for _, item := range itens {
		produto, err := produtoRepoTx.FindByID(ctx, item.ProductID)

		if err != nil {
			return model.Pedido{}, err
		}

		if produto.StockQuantity < item.Quantity {
			return model.Pedido{}, model.ErrEstoqueInsuficiente
		}

		//atualiza o estoque do produto
		estoqueAtualizado := produto.StockQuantity - item.Quantity
		if err := produtoRepoTx.UpdateStock(ctx, produto.ID, estoqueAtualizado); err != nil {
			return model.Pedido{}, fmt.Errorf("erro ao atualizar estoque: %w", err)
		}

		itemPedido := model.ItemPedido{
			ID:        uuid.New(),
			PedidoID:  pedido.ID,
			ProdutoID: produto.ID,
			Quantity:  item.Quantity,
			UnitPrice: produto.Price,
		}

		//adiciona item à lista de itens para salvar
		itensParaSalvar = append(itensParaSalvar, itemPedido)
		total = total.Add(produto.Price.Mul(decimal.NewFromInt(int64(item.Quantity))))
	}

	//atualiza o total do pedido
	pedido.TotalAmount = total

	if err := pedidoRepoTx.Save(ctx, pedido); err != nil {
		return model.Pedido{}, err
	}

	//salva os itens do pedido
	for _, item := range itensParaSalvar {
		if err := itemPedidoRepoTx.Save(ctx, item); err != nil {
			return model.Pedido{}, err
		}
	}
	//confirma a transação
	if err := tx.Commit(ctx); err != nil {
		return model.Pedido{}, fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return pedido, nil
}

func (s *PedidoService) FindByID(ctx context.Context, id uuid.UUID) (model.Pedido, error) {
	return s.pedidoRepo.FindByID(ctx, id)
}

func (s *PedidoService) FindAll(ctx context.Context, limit, offset int) ([]model.Pedido, error) {
	return s.pedidoRepo.FindAll(ctx, limit, offset)
}

// TODO: implementar PagarPedido e CancelarPedido (com devolução de estoque via transação)
