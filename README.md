
# Desafio Clean Architecture

Este projeto é um serviço de pedidos desenvolvido em Go, que expõe APIs REST, GraphQL e gRPC para manipulação de dados de pedidos. 
A aplicação segue os princípios de arquitetura limpa.

## Requisitos

- Go 1.24+
- Docker e Docker Compose
- PostgreSQL

## Instalação

1. Clone o repositório:

```bash
git clone <url-do-repositorio>
cd clean-architecture
```

2. Baixe as dependências Go:

```bash
go mod tidy
```

## Execução com Docker Compose

1. Configure as variáveis de ambiente:

```env
POSTGRES_USER=user
POSTGRES_PASSWORD=password
POSTGRES_DB=orders
```

2. Suba os containers:

```bash
docker compose up --build
```

## Testes

### Verificar se o banco de dados foi criado

Dentro do container do Postgres:

```bash
docker exec -it clean-architecture-postgres-1 psql -U {seu_user} -d {nome_do_seu_db}
```

Saída em caso de sucesso:

```bash
psql (15.13)
Type "help" for help.
{nome_do_seu_db}=#
```

Ainda no Postgres, execute:

```sql
INSERT INTO orders (product_name, quantity, created_at) VALUES ('Iphone 18', 2, NOW());
```

### Testar API REST

Teste via url: [http://localhost:8080/orders](http://localhost:8080/orders)

Teste via terminal:

```bash
curl http://localhost:8080/orders
```

Saída impressa:

```bash
[
  {
    "id": 1,
    "product_name": "Iphone 18",
    "quantity": 2,
    "created_at": "2025-05-21T14:27:04.752741Z"
  },
  {...}
]
```


### Testar GraphQL

GraphQL disponível na porta `8081`.
Teste via navegador acessando a url: [http://localhost:8081](http://localhost:8081)
Utilize a seguinte quwey para teste:

```graphql
query {
  orders {
    id
    productName
    quantity
    createdAt
  }
}
```

Ou via terminal com o comando "**curl**":

```bash
curl -X POST http://localhost:8081/graphql \
  -H "Content-Type: application/json" \
  -d '{"query":"query { orders { id productName quantity createdAt } }"}'
```

Repare que a saída será o conteúdo do arquivo html apresentando pelo GraphQL

### Testar gRPC

gRPC disponível na porta `50051`.

Acessando com **grpcurl**:

Entrada no terminal:

```bash
grpcurl -plaintext localhost:50051 list
```

Saída no terminal:

```bash
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
order.OrderService
```

Entrada no terminal:

```bash
grpcurl -plaintext localhost:50051 order.OrderService/ListOrders
```

Saída no terminal:

```bash
{
  "orders": [
    {
      "id": 1,
      "productName": "Iphone",
      "quantity": 2,
      "createdAt": "2025-05-21T14:27:04Z"
    },
    {...}
  ]
}
```

Via **grpcui**:

```bash
grpcui -plaintext localhost:50051
```