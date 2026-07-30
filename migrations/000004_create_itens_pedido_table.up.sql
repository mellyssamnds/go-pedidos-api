CREATE TABLE itens_pedido (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    pedido_id UUID NOT NULL,
    produto_id UUID NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (pedido_id) REFERENCES pedidos(id),
    FOREIGN KEY (produto_id) REFERENCES produtos(id)
);