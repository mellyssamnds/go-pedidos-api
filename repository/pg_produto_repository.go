package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mellyssamnds/go-pedidos-api/model"
)

type PgProdutoRepository struct {
	db Querier
}

func NewPgProdutoRepository(db Querier) *PgProdutoRepository {
	return &PgProdutoRepository{db: db}
}

var _ ProdutoRepository = (*PgProdutoRepository)(nil)

func (r *PgProdutoRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Produto, error) {
	query := `
		SELECT id, name, description, price, stock_quantity, created_at
		FROM produtos
		WHERE id = $1
	`
	var produto model.Produto
	err := r.db.QueryRow(ctx, query, id).Scan(
		&produto.ID,
		&produto.Name,
		&produto.Description,
		&produto.Price,
		&produto.StockQuantity,
		&produto.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Produto{}, model.ErrProdutoNaoEncontrado
		}
		return model.Produto{}, fmt.Errorf("erro ao buscar produto por id: %w", err)
	}

	return produto, nil
}

func (r *PgProdutoRepository) Save(ctx context.Context, produto model.Produto) error {
	query := `
		INSERT INTO produtos (id, name, description, price, stock_quantity, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query,
		produto.ID,
		produto.Name,
		produto.Description,
		produto.Price,
		produto.StockQuantity,
		produto.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar produto: %w", err)
	}

	return nil
}

func (r *PgProdutoRepository) FindAll(ctx context.Context) ([]model.Produto, error) {
	query := `
		SELECT id, name, description, price, stock_quantity, created_at
		FROM produtos
	`
	var produtos []model.Produto
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var produto model.Produto
		err := rows.Scan(
			&produto.ID,
			&produto.Name,
			&produto.Description,
			&produto.Price,
			&produto.StockQuantity,
			&produto.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler dados do produto: %w", err)
		}

		produtos = append(produtos, produto)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %w", rows.Err())
	}

	return produtos, nil
}
