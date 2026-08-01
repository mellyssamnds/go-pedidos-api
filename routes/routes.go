package routes

import (
	"net/http"

	"github.com/mellyssamnds/go-pedidos-api/controllers"
)

func NewRouter(
	clienteController *controllers.ClienteController,
	produtoController *controllers.ProdutoController,
	pedidoController *controllers.PedidoController,
) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /clientes", clienteController.Create)
	mux.HandleFunc("GET /clientes", clienteController.FindAll)
	mux.HandleFunc("GET /clientes/{id}", clienteController.FindByID)

	mux.HandleFunc("POST /produtos", produtoController.Create)
	mux.HandleFunc("GET /produtos", produtoController.FindAll)
	mux.HandleFunc("GET /produtos/{id}", produtoController.FindByID)

	mux.HandleFunc("POST /pedidos", pedidoController.Create)
	mux.HandleFunc("GET /pedidos", pedidoController.FindAll)
	mux.HandleFunc("GET /pedidos/{id}", pedidoController.FindByID)

	return mux
}
