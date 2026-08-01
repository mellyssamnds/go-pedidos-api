//go:build integration

package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mellyssamnds/go-pedidos-api/controllers"
	"github.com/mellyssamnds/go-pedidos-api/model"
	"github.com/mellyssamnds/go-pedidos-api/repository"
	"github.com/mellyssamnds/go-pedidos-api/service"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestPedidoController_Create_Integration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@127.0.0.1:5433/api_pedidos_test?sslmode=disable"
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	assert.NoError(t, err)
	t.Cleanup(func() { pool.Close() })

	//cria cliente e produto para teste
	ctx := context.Background()
	clienteRepo := repository.NewPgClienteRepository(pool)
	produtoRepo := repository.NewPgProdutoRepository(pool)
	pedidoRepo := repository.NewPgPedidoRepository(pool)
	itemPedidoRepo := repository.NewPgItemPedidoRepository(pool)

	cliente := model.Cliente{
		ID: uuid.New(), Name: "Cliente Controller", Email: "controller@teste.com",
		PasswordHash: "hash-fake", CreatedAt: time.Now(),
	}
	assert.NoError(t, clienteRepo.Save(ctx, cliente))

	produto := model.Produto{
		ID: uuid.New(), Name: "Produto Controller", Description: "desc",
		Price: decimal.NewFromFloat(20.0), StockQuantity: 10, CreatedAt: time.Now(),
	}
	assert.NoError(t, produtoRepo.Save(ctx, produto))

	//cria pedido via controller
	pedidoService := service.NewPedidoService(pedidoRepo, clienteRepo, produtoRepo, itemPedidoRepo, pool)
	controller := controllers.NewPedidoController(pedidoService)

	payload := `{"clienteId":"` + cliente.ID.String() + `","itens":[{"produtoId":"` + produto.ID.String() + `","quantity":2}]}`
	req := httptest.NewRequest(http.MethodPost, "/pedidos", bytes.NewBufferString(payload))
	rec := httptest.NewRecorder()

	controller.Create(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var resposta map[string]any
	json.Unmarshal(rec.Body.Bytes(), &resposta)
	assert.Equal(t, "PENDING", resposta["status"])

	//limpa banco de teste
	t.Cleanup(func() {
		pool.Exec(context.Background(), "DELETE FROM itens_pedido WHERE produto_id = $1", produto.ID)
		pool.Exec(context.Background(), "DELETE FROM pedidos WHERE cliente_id = $1", cliente.ID)
		pool.Exec(context.Background(), "DELETE FROM produtos WHERE id = $1", produto.ID)
		pool.Exec(context.Background(), "DELETE FROM clientes WHERE id = $1", cliente.ID)
	})
}
