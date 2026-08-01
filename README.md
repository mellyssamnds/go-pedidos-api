# API de Pedidos

![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18-4169E1?logo=postgresql)
![License](https://img.shields.io/badge/License-MIT-green)

API REST desenvolvida em **Go** para gerenciamento de **clientes, produtos e pedidos**, utilizando **PostgreSQL** como banco de dados e garantindo a criação de pedidos de forma **100% transacional**.

---

## Índice

- [Funcionalidades](#funcionalidades)
- [Tecnologias](#tecnologias)
- [Arquitetura](#arquitetura)
- [Estrutura do projeto](#estrutura-do-projeto)
- [Como executar](#como-executar)
- [Endpoints](#endpoints)
- [Testes](#testes)
- [Decisões técnicas](#decisões-técnicas)
- [Funcionalidades planejadas](#funcionalidades-planejadas)

## Funcionalidades

- Cadastro de clientes
- Consulta de clientes
- Cadastro de produtos
- Consulta de produtos
- Criação de pedidos com transação
- Controle de estoque durante a criação do pedido
- Persistência em PostgreSQL
- Testes unitários e de integração

## Tecnologias

| Tecnologia | Finalidade |
|------------|------------|
| Go 1.26 | Linguagem principal |
| net/http (biblioteca padrão) | Servidor HTTP |
| PostgreSQL 18 | Banco de dados |
| pgx/v5 + pgxpool | Driver e pool de conexões |
| golang-migrate | Versionamento do banco |
| testify | Testes e mocks |
| bcrypt | Hash de senhas |
| shopspring/decimal | Precisão para valores monetários |
| Docker Compose | Ambiente de desenvolvimento |

## Arquitetura

O projeto adota uma **arquitetura em camadas**, onde cada camada possui uma responsabilidade específica e depende apenas das abstrações da camada inferior.

```text
config/
├── carregamento das variáveis de ambiente

database/
├── conexão com PostgreSQL

model/
├── entidades de domínio
└── erros conhecidos

repository/
├── interfaces
├── acesso aos dados
└── implementações PostgreSQL

service/
└── regras de negócio

controllers/
└── camada HTTP

routes/
└── registro das rotas

migrations/
└── versionamento do banco
```

Fluxo de dependências:

```text
HTTP
  │
Controllers
  │
Services
  │
Repositories
  │
Database
  │
PostgreSQL
```

Essa organização reduz o acoplamento entre as camadas, facilita a manutenção e permite testar a regra de negócio sem depender de um banco de dados real.

## Estrutura do projeto

```text
.
.
├── config/
├── controllers/
├── database/
├── migrations/
├── mocks/
├── model/
├── repository/
├── routes/
├── scripts/
├── service/
├── .env
├── .env.example
├── .gitignore
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Como executar

### Pré-requisitos

- Go 1.26
- Docker Desktop
- golang-migrate CLI

### 1. Suba o banco de dados

```bash
docker compose up -d
```

Serão criados automaticamente os bancos:

- `api_pedidos`
- `api_pedidos_test`

O PostgreSQL será exposto na porta **5433**, evitando conflitos com instalações locais.

### 2. Configure as variáveis de ambiente

Copie o arquivo de exemplo:

```bash
cp .env.example .env
```

Caso necessário, ajuste os valores conforme o seu ambiente. O arquivo `.env.example` já contém a configuração padrão utilizada pelo projeto.

### 3. Execute as migrations

```bash
migrate \
-database "postgres://postgres:postgres@127.0.0.1:5433/api_pedidos?sslmode=disable" \
-path migrations up
```

### 4. Execute a aplicação

```bash
go run main.go
```

A API ficará disponível em:

```text
http://localhost:8080
```

## Endpoints

| Recurso | Método | Endpoint | Descrição |
|----------|---------|----------|-----------|
| Clientes | POST | `/clientes` | Cadastra um cliente |
| Clientes | GET | `/clientes` | Lista clientes |
| Clientes | GET | `/clientes/{id}` | Busca um cliente |
| Produtos | POST | `/produtos` | Cadastra um produto |
| Produtos | GET | `/produtos` | Lista produtos |
| Produtos | GET | `/produtos/{id}` | Busca um produto |
| Pedidos | POST | `/pedidos` | Cria um pedido (transacional) |
| Pedidos | GET | `/pedidos` | Lista pedidos |
| Pedidos | GET | `/pedidos/{id}` | Busca um pedido |

A listagem de pedidos suporta paginação pelos parâmetros:

| Parâmetro | Descrição |
|-----------|-----------|
| `limit` | Quantidade de registros |
| `offset` | Registro inicial |

Exemplo:

```http
GET /pedidos?limit=10&offset=0
```

## Testes

O projeto possui dois tipos de testes.

### Testes unitários

Executam apenas a lógica da aplicação, utilizando mocks dos repositórios.

```bash
go test ./...
```

### Testes de integração

Executam contra um banco PostgreSQL real.

Aplicar as migrations:

```bash
migrate \
-database "postgres://postgres:postgres@127.0.0.1:5433/api_pedidos_test?sslmode=disable" \
-path migrations up
```

Executar os testes:

```bash
go test -tags=integration ./...
```

### Cobertura

```bash
go test -tags=integration -coverprofile=coverage.out ./...

go tool cover -func=coverage.out
```

---

## Decisões técnicas

### PostgreSQL na porta 5433

O banco é exposto na porta **5433** em vez da porta padrão **5432** para evitar conflitos com instalações locais do PostgreSQL. Dessa forma, o ambiente pode ser executado sem alterar a configuração existente na máquina do desenvolvedor.

### Criação de pedidos com transação

A criação de um pedido é realizada dentro de uma única transação do banco de dados. O processo engloba a validação do estoque, atualização das quantidades disponíveis, criação do pedido e persistência de seus itens.

Caso qualquer etapa falhe, toda a operação é revertida, garantindo consistência entre pedidos e estoque.

### Interface `Querier`

Os repositórios dependem da interface `Querier`, implementada tanto por `*pgxpool.Pool` quanto por `pgx.Tx`. Essa abstração permite reutilizar a mesma implementação em operações comuns e transacionais, evitando duplicação de código e reduzindo o acoplamento com a camada de persistência.

### Persistência do preço do produto

O preço do produto é copiado para o item do pedido no momento da compra. Dessa forma, alterações futuras no cadastro do produto não modificam o histórico de pedidos já registrados.

### Services sem interfaces

A camada de **Service** não possui interfaces próprias, pois cada serviço possui apenas uma implementação. O desacoplamento necessário para testes é obtido pelas interfaces dos repositórios, mantendo a arquitetura simples e evitando abstrações desnecessárias.

## Funcionalidades planejadas

As seguintes funcionalidades ainda não foram implementadas:

- `POST /pedidos/{id}/pagar`
- `POST /pedidos/{id}/cancelar`

Os pontos de extensão encontram-se documentados em:

```text
service/pedido_service.go
```