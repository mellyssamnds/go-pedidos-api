package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ItemPedido struct {
	ID        uuid.UUID       `db:"id" json:"id"`
	PedidoID  uuid.UUID       `db:"pedido_id" json:"pedidoId"`
	ProdutoID uuid.UUID       `db:"produto_id" json:"produtoId"`
	Quantity  int             `db:"quantity" json:"quantity"`
	UnitPrice decimal.Decimal `db:"unit_price" json:"unitPrice"`
}
