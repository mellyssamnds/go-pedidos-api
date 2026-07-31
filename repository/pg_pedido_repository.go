package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mellyssamnds/go-pedidos-api/model"
)

type PgPedidoRepository struct {
	pool *pgxpool.Pool
}

func NewPgPedidoRepository(pool *pgxpool.Pool) *PgPedidoRepository {
	return &PgPedidoRepository{pool: pool}
}

var _ PedidoRepository = (*PgPedidoRepository)(nil)

func (r *PgPedidoRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Pedido, error) {
	query := `
		SELECT id, cliente_id, created_at, status, total_amount
		FROM pedidos
		WHERE id = $1
	`

	var pedido model.Pedido
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&pedido.ID,
		&pedido.ClienteID,
		&pedido.CreatedAt,
		&pedido.Status,
		&pedido.TotalAmount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Pedido{}, model.ErrPedidoNaoEncontrado
		}
		return model.Pedido{}, fmt.Errorf("erro ao buscar pedido por id: %w", err)
	}

	return pedido, nil
}

func (r *PgPedidoRepository) Save(ctx context.Context, pedido model.Pedido) error {
	query := `
		INSERT INTO pedidos (id, cliente_id, created_at, status, total_amount)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		pedido.ID,
		pedido.ClienteID,
		pedido.CreatedAt,
		pedido.Status,
		pedido.TotalAmount,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar pedido: %w", err)
	}

	return nil
}

func (r *PgPedidoRepository) FindAll(ctx context.Context, limit, offset int) ([]model.Pedido, error) {
	query := `
		SELECT id, cliente_id, created_at, status, total_amount
		FROM pedidos
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos: %w", err)
	}
	defer rows.Close()

	var pedidos []model.Pedido
	for rows.Next() {
		var pedido model.Pedido
		err := rows.Scan(
			&pedido.ID,
			&pedido.ClienteID,
			&pedido.CreatedAt,
			&pedido.Status,
			&pedido.TotalAmount,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler dados do pedido: %w", err)
		}
		pedidos = append(pedidos, pedido)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos: %w", rows.Err())
	}

	return pedidos, nil
}

func (r *PgPedidoRepository) UpdateStatus(ctx context.Context, pedidoID uuid.UUID, status model.OrderStatus) error {
	query := `
		UPDATE pedidos
		SET status = $1
		WHERE id = $2
	`
	result, err := r.pool.Exec(ctx, query, status, pedidoID)

	if err != nil {
		return fmt.Errorf("erro ao atualizar status do pedido: %w", err)
	}

	if result.RowsAffected() == 0 {
		return model.ErrPedidoNaoEncontrado
	}

	return nil
}
