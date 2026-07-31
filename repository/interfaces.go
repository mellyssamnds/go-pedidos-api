package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
)

type ClienteRepository interface {
	Save(ctx context.Context, cliente model.Cliente) error
	FindByID(ctx context.Context, id uuid.UUID) (model.Cliente, error)
	FindByEmail(ctx context.Context, email string) (model.Cliente, error)
	FindAll(ctx context.Context) ([]model.Cliente, error)
}

type ProdutoRepository interface {
	Save(ctx context.Context, produto model.Produto) error
	FindByID(ctx context.Context, id uuid.UUID) (model.Produto, error)
	FindAll(ctx context.Context) ([]model.Produto, error)
}

type PedidoRepository interface {
	Save(ctx context.Context, pedido model.Pedido) error
	FindByID(ctx context.Context, id uuid.UUID) (model.Pedido, error)
	FindAll(ctx context.Context, limit, offset int) ([]model.Pedido, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error
}

type ItemPedidoRepository interface {
	Save(ctx context.Context, itemPedido model.ItemPedido) error
	FindByID(ctx context.Context, id uuid.UUID) (model.ItemPedido, error)
	FindByPedidoID(ctx context.Context, pedidoID uuid.UUID) ([]model.ItemPedido, error)
}
