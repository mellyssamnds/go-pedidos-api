package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Produto struct {
	ID            uuid.UUID       `db:"id" json:"id"`
	Name          string          `db:"name" json:"name"`
	Description   string          `db:"description" json:"description"`
	Price         decimal.Decimal `db:"price" json:"price"`
	StockQuantity int             `db:"stock_quantity" json:"stockQuantity"`
	CreatedAt     time.Time       `db:"created_at" json:"createdAt"`
}
