# Sistema de Gerenciamento de Restaurante

API completa para gerenciamento de restaurantes, desenvolvida com Go, Gin, GORM e PostgreSQL.

## Tecnologias

- **Go**: Linguagem de programação backend
- **Gin**: Framework web para Go
- **GORM**: ORM (Object-Relational Mapping) para Go
- **PostgreSQL**: Banco de dados relacional
- **JWT**: Autenticação com JSON Web Tokens
- **Docker**: Containerização
- **Docker Compose**: Orquestração de containers

## Arquitetura

O projeto utiliza uma arquitetura em camadas seguindo princípios de Clean Architecture e Domain-Driven Design (DDD):

- **Domain**: Contém os modelos e interfaces de repositórios
- **Infrastructure**: Implementações concretas de repositórios e serviços externos
- **Services**: Lógica de negócios
- **API**: Manipuladores HTTP e rotas
- **Config**: Configurações da aplicação

## Funcionalidades

- **Autenticação e Autorização**
  - Login com JWT
  - Níveis de acesso (admin, manager, staff)
  - Rotas protegidas

- **Gerenciamento de Mesas**
  - Cadastro, edição e exclusão de mesas
  - Controle de status (livre, ocupada, reservada)

- **Gestão de Pedidos**
  - Criação de pedidos
  - Adição/remoção de itens
  - Acompanhamento do status do pedido
  - Associação de pedidos a mesas

- **Cardápio**
  - Cadastro de produtos
  - Categorização (comida, bebida, sobremesa)
  - Controle de estoque

- **Controle Financeiro**
  - Registro de receitas e despesas
  - Categorização de transações
  - Relatórios diários e mensais

## Inicialização

### Requisitos

- Docker e Docker Compose
- Go 1.18+

### Passos para execução

1. Clone o repositório
   ```bash
   git clone https://github.com/seu-usuario/restaurant-management-api.git
   cd api-jet-manager
   ```

2. Configure as variáveis de ambiente
   ```bash
   cp .env.example .env
   # Edite o arquivo .env com suas configurações
   ```

3. Execute com Docker Compose
   ```bash
   docker-compose up -d
   ```

4. Acesse a API em `http://localhost:8080`

### Execução sem Docker

1. Certifique-se de ter PostgreSQL instalado e rodando
2. Configure as variáveis de ambiente no arquivo `.env`
3. Execute:
   ```bash
   go run cmd/server/main.go
   ```

## Endpoints da API

### Autenticação
- `POST /api/auth/login`: Login de usuário
- `POST /api/auth/register`: Registro de usuário (apenas admin)

### Mesas
- `GET /api/tables`: Listar todas as mesas
- `GET /api/tables/:id`: Obter detalhes de uma mesa
- `POST /api/tables`: Criar uma nova mesa
- `PUT /api/tables/:id`: Atualizar uma mesa
- `DELETE /api/tables/:id`: Excluir uma mesa
- `PATCH /api/tables/:id/status`: Atualizar status de uma mesa

### Pedidos
- `GET /api/orders`: Listar todos os pedidos
- `GET /api/orders/:id`: Obter detalhes de um pedido
- `POST /api/orders`: Criar um novo pedido
- `PATCH /api/orders/:id/status`: Atualizar status de um pedido
- `POST /api/orders/:id/items`: Adicionar item a um pedido
- `DELETE /api/orders/:id/items/:item_id`: Remover item de um pedido

### Produtos
- `GET /api/products`: Listar todos os produtos
- `GET /api/products/:id`: Obter detalhes de um produto
- `POST /api/products`: Criar um novo produto
- `PUT /api/products/:id`: Atualizar um produto
- `DELETE /api/products/:id`: Excluir um produto
- `PATCH /api/products/:id/stock`: Atualizar estoque de um produto

### Finanças
- `GET /api/finance/transactions`: Listar todas as transações
- `GET /api/finance/transactions/:id`: Obter detalhes de uma transação
- `POST /api/finance/transactions`: Criar uma nova transação
- `PUT /api/finance/transactions/:id`: Atualizar uma transação
- `DELETE /api/finance/transactions/:id`: Excluir uma transação
- `GET /api/finance/summary`: Obter resumo financeiro

## Licença

Este projeto está licenciado sob a licença MIT.