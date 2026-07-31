package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mellyssamnds/go-pedidos-api/model"
)

type PgClienteRepository struct {
	db Querier
}

func NewPgClienteRepository(db Querier) *PgClienteRepository {
	return &PgClienteRepository{db: db}
}

var _ ClienteRepository = (*PgClienteRepository)(nil)

func (r *PgClienteRepository) FindByID(ctx context.Context, id uuid.UUID) (model.Cliente, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM clientes
		WHERE id = $1
	`

	var cliente model.Cliente
	err := r.db.QueryRow(ctx, query, id).Scan(
		&cliente.ID,
		&cliente.Name,
		&cliente.Email,
		&cliente.PasswordHash,
		&cliente.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Cliente{}, model.ErrClienteNaoEncontrado
		}
		return model.Cliente{}, fmt.Errorf("erro ao buscar cliente por id: %w", err)
	}

	return cliente, nil
}

func (r *PgClienteRepository) Save(ctx context.Context, cliente model.Cliente) error {
	query := `
		INSERT INTO clientes (id, name, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(ctx, query,
		cliente.ID,
		cliente.Name,
		cliente.Email,
		cliente.PasswordHash,
		cliente.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("erro ao salvar cliente: %w", err)
	}

	return nil
}

func (r *PgClienteRepository) FindByEmail(ctx context.Context, email string) (model.Cliente, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM clientes
		WHERE email = $1
	`

	var cliente model.Cliente
	err := r.db.QueryRow(ctx, query, email).Scan(
		&cliente.ID,
		&cliente.Name,
		&cliente.Email,
		&cliente.PasswordHash,
		&cliente.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Cliente{}, model.ErrClienteNaoEncontrado
		}
		return model.Cliente{}, fmt.Errorf("erro ao buscar cliente por email: %w", err)
	}

	return cliente, nil
}

func (r *PgClienteRepository) FindAll(ctx context.Context) ([]model.Cliente, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM clientes
	`

	var clientes []model.Cliente
	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("erro ao buscar clientes: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var cliente model.Cliente
		err := rows.Scan(
			&cliente.ID,
			&cliente.Name,
			&cliente.Email,
			&cliente.PasswordHash,
			&cliente.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear cliente: %w", err)
		}
		clientes = append(clientes, cliente)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os clientes: %w", rows.Err())
	}

	return clientes, nil
}
