package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mellyssamnds/go-pedidos-api/config"
	"github.com/mellyssamnds/go-pedidos-api/controllers"
	"github.com/mellyssamnds/go-pedidos-api/database"
	"github.com/mellyssamnds/go-pedidos-api/repository"
	"github.com/mellyssamnds/go-pedidos-api/routes"
	"github.com/mellyssamnds/go-pedidos-api/service"
)

func main() {
	//carrega a configuração
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("erro ao carregar configuração: %v", err)
	}

	//conecta banco de dados
	pool, err := database.Connect(context.Background(), cfg.GetDSN())
	if err != nil {
		log.Fatalf("erro ao conectar ao banco de dados: %v", err)
	}
	defer pool.Close()

	log.Println("conexão com o banco de dados feita com sucesso!")

	//repositories
	clienteRepo := repository.NewPgClienteRepository(pool)
	produtoRepo := repository.NewPgProdutoRepository(pool)
	pedidoRepo := repository.NewPgPedidoRepository(pool)
	itemPedidoRepo := repository.NewPgItemPedidoRepository(pool)

	//services
	clienteService := service.NewClienteService(clienteRepo)
	produtoService := service.NewProdutoService(produtoRepo)
	pedidoService := service.NewPedidoService(pedidoRepo, clienteRepo, produtoRepo, itemPedidoRepo, pool)

	//controllers
	clienteController := controllers.NewClienteController(clienteService)
	produtoController := controllers.NewProdutoController(produtoService)
	pedidoController := controllers.NewPedidoController(pedidoService)

	//rotas
	mux := routes.NewRouter(clienteController, produtoController, pedidoController)

	log.Println("servidor rodando na porta 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
