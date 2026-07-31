package model

import "errors"

var (
	ErrClienteNaoEncontrado    = errors.New("cliente não encontrado")
	ErrProdutoNaoEncontrado    = errors.New("produto não encontrado")
	ErrPedidoNaoEncontrado     = errors.New("pedido não encontrado")
	ErrItemPedidoNaoEncontrado = errors.New("item do pedido não encontrado")

	ErrEstoqueInsuficiente = errors.New("estoque insuficiente")
	ErrEmailCadastrado     = errors.New("email já cadastrado")
	ErrStatusInvalido      = errors.New("mudança de status inválida")
)
