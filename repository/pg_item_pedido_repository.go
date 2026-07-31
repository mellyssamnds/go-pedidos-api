package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mellyssamnds/go-pedidos-api/model"
)

type PgItemPedidoRepository struct {
	db Querier
}

func NewPgItemPedidoRepository(db Querier) *PgItemPedidoRepository {
	return &PgItemPedidoRepository{db: db}
}

var _ ItemPedidoRepository = (*PgItemPedidoRepository)(nil)

func (r *PgItemPedidoRepository) Save(ctx context.Context, itemPedido model.ItemPedido) error {
	query := `
		INSERT INTO itens_pedido (id, pedido_id, produto_id, quantity, unit_price)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query, itemPedido.ID, itemPedido.PedidoID, itemPedido.ProdutoID, itemPedido.Quantity, itemPedido.UnitPrice)
	if err != nil {
		return fmt.Errorf("falha ao criar itemPedido: %w", err)
	}

	return nil
}

func (r *PgItemPedidoRepository) FindByID(ctx context.Context, id uuid.UUID) (model.ItemPedido, error) {
	query := `
		SELECT id, pedido_id, produto_id, quantity, unit_price
		FROM itens_pedido
		WHERE id = $1
	`

	row := r.db.QueryRow(ctx, query, id)
	var itemPedido model.ItemPedido
	err := row.Scan(&itemPedido.ID, &itemPedido.PedidoID, &itemPedido.ProdutoID, &itemPedido.Quantity, &itemPedido.UnitPrice)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ItemPedido{}, model.ErrItemPedidoNaoEncontrado
		}

		return model.ItemPedido{}, fmt.Errorf("falha ao obter itemPedido por ID: %w", err)
	}

	return itemPedido, nil
}

func (r *PgItemPedidoRepository) FindByPedidoID(ctx context.Context, pedidoID uuid.UUID) ([]model.ItemPedido, error) {
	query := `
		SELECT id, pedido_id, produto_id, quantity, unit_price
		FROM itens_pedido
		WHERE pedido_id = $1
	`

	rows, err := r.db.Query(ctx, query, pedidoID)

	if err != nil {
		return nil, fmt.Errorf("falha ao obter itemPedido por pedido_id: %w", err)
	}

	defer rows.Close()

	var itemPedidos []model.ItemPedido
	for rows.Next() {
		var itemPedido model.ItemPedido
		err := rows.Scan(
			&itemPedido.ID,
			&itemPedido.PedidoID,
			&itemPedido.ProdutoID,
			&itemPedido.Quantity,
			&itemPedido.UnitPrice,
		)

		if err != nil {
			return nil, fmt.Errorf("falha ao ler dados de itemPedidos: %w", err)
		}

		itemPedidos = append(itemPedidos, itemPedido)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("falha ao iterar itemPedidos: %w", err)
	}

	return itemPedidos, nil
}
