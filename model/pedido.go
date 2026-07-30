package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Pedido struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	ClienteID   uuid.UUID       `db:"cliente_id" json:"clienteId"`
	CreatedAt   time.Time       `db:"created_at" json:"createdAt"`
	Status      OrderStatus     `db:"status" json:"status"`
	TotalAmount decimal.Decimal `db:"total_amount" json:"totalAmount"`
}
