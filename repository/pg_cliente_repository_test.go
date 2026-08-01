//go:build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/repository"
	"github.com/stretchr/testify/assert"
)

func TestPgClienteRepository_FindByEmail(t *testing.T) {
	//configura o banco de teste
	pool := setupTestDB(t)
	repo := repository.NewPgClienteRepository(pool)
	ctx := context.Background()

	cliente := model.Cliente{
		ID:           uuid.New(),
		Name:         "Cliente Email Teste",
		Email:        "email@teste.com",
		PasswordHash: "hash-fake",
		CreatedAt:    time.Now(),
	}
	assert.NoError(t, repo.Save(ctx, cliente))

	encontrado, err := repo.FindByEmail(ctx, cliente.Email)
	assert.NoError(t, err)
	assert.Equal(t, cliente.Name, encontrado.Name)
	assert.Equal(t, cliente.Email, encontrado.Email)

	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM clientes WHERE id = $1", cliente.ID)
	})
}

func TestPgClienteRepository_FindByEmail_NaoEncontrado(t *testing.T) {
	pool := setupTestDB(t)
	repo := repository.NewPgClienteRepository(pool)

	_, err := repo.FindByEmail(context.Background(), "naoexiste@teste.com")

	assert.ErrorIs(t, err, model.ErrClienteNaoEncontrado)
}

func TestPgClienteRepository_FindAll(t *testing.T) {
	pool := setupTestDB(t)
	repo := repository.NewPgClienteRepository(pool)
	ctx := context.Background()

	cliente := model.Cliente{
		ID:           uuid.New(),
		Name:         "Cliente FindAll Teste",
		Email:        "findalltest@teste.com",
		PasswordHash: "hash-fake",
		CreatedAt:    time.Now(),
	}
	assert.NoError(t, repo.Save(ctx, cliente))

	clientes, err := repo.FindAll(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, clientes)

	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM clientes WHERE id = $1", cliente.ID)
	})
}
