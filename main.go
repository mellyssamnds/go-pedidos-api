package main

import (
	"context"
	"fmt"

	"github.com/mellyssamnds/api-pedidos/config"
	"github.com/mellyssamnds/api-pedidos/database"
)

func main() {
	// Carregar configuração
	cfg, err := config.Load()

	if err != nil {
		fmt.Printf("Erro ao carregar configuração: %v\n", err)
		return
	}

	// Conectar ao banco de dados
	ctx := context.Background()
	dsn := cfg.GetDSN()
	dbPool, err := database.Connect(ctx, dsn)

	if err != nil {
		fmt.Printf("Erro ao conectar ao banco de dados: %v\n", err)
		return
	}

	defer dbPool.Close()

	fmt.Println("Conexão com o banco de dados feita com sucesso!")
}
