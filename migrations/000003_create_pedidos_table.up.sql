CREATE TYPE order_status AS ENUM ('PENDING', 'PAID', 'CANCELED');

CREATE TABLE pedidos (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    cliente_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status order_status NOT NULL DEFAULT 'PENDING',
    total_amount DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (cliente_id) REFERENCES clientes(id)
);