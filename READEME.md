# LEP System - Backend API

## Visão Geral

O LEP System é uma aplicação backend robusta desenvolvida em Go, utilizando arquitetura limpa e modular. O sistema foi projetado para gestão completa de operações empresariais, incluindo usuários, produtos, compras e pedidos.

### Tecnologias Utilizadas

- **Go 1.23.0** - Linguagem de programação principal
- **Gin Web Framework** - Framework HTTP para APIs RESTful
- **GORM** - ORM para manipulação de banco de dados
- **PostgreSQL** - Banco de dados principal
- **Google Cloud Platform** - Infraestrutura de nuvem (Cloud Run, Cloud SQL, Secret Manager)
- **Terraform** - Infrastructure as Code
- **JWT** - Autenticação e autorização
- **bcrypt** - Criptografia de senhas

### Características Principais

- ✅ **Arquitetura Limpa** - Separação clara entre camadas (Handler, Service, Repository)
- ✅ **API RESTful** - Endpoints padronizados seguindo convenções REST
- ✅ **Autenticação JWT** - Sistema seguro de autenticação e autorização
- ✅ **Validação de Headers** - Controle organizacional via headers obrigatórios
- ✅ **Soft Delete** - Remoção lógica de registros para auditoria
- ✅ **CRUD Completo** - Operações completas para todas as entidades
- ✅ **Logs de Auditoria** - Rastreamento completo de operações
- ✅ **Deploy Automatizado** - Configuração Terraform para GCP

---

## Estrutura do Projeto

```
lep-system/
├─ config/             # Configurações da aplicação
├─ handler/            # Camada de negócio - Interfaces e implementações
│  ├─ auth.go         # Autenticação e autorização
│  ├─ user.go         # Gestão de usuários
│  ├─ product.go      # Gestão de produtos
│  ├─ purchase.go     # Gestão de compras
│  ├─ order.go        # Gestão de pedidos
│  └─ inject.go       # Injeção de dependências dos handlers
├─ repositories/       # Camada de dados - Acesso ao banco
│  ├─ models/         # Definições de entidades/modelos
│  ├─ migrate/        # Scripts de migração
│  └─ *.go           # Implementações dos repositórios
├─ server/            # Camada de apresentação - Controladores HTTP
│  ├─ auth.go        # Endpoints de autenticação
│  ├─ user.go        # Endpoints de usuários
│  ├─ product.go     # Endpoints de produtos
│  ├─ purchase.go    # Endpoints de compras
│  ├─ order.go       # Endpoints de pedidos
│  └─ inject.go      # Injeção de dependências dos servers
├─ routes/            # Organização e configuração de rotas
│  ├─ router.go      # Configuração principal das rotas
│  └─ routes.md      # Documentação das rotas
├─ resource/          # Gerenciamento de recursos e injeção global
├─ utils/             # Funções utilitárias
├─ example.main.tf    # Configuração Terraform para GCP
└─ main.go           # Ponto de entrada da aplicação
```

---

## Arquitetura e Funcionalidades

### Padrão de Arquitetura

O sistema segue o padrão de **Arquitetura Limpa** com três camadas principais:

1. **Handler Layer** (`handler/`)
   - Contém a lógica de negócio
   - Interfaces bem definidas para cada domínio
   - Validação de regras de negócio
   - Criptografia de senhas e processamento de dados

2. **Server Layer** (`server/`)
   - Controladores HTTP (similar a Controllers no MVC)
   - Validação de headers obrigatórios
   - Processamento de requisições e respostas
   - Padronização de responses em JSON

3. **Repository Layer** (`repositories/`)
   - Acesso direto ao banco de dados via GORM
   - Implementação de operações CRUD
   - Gestão de conexões e transações

### Funcionalidades Implementadas

- **🔐 Autenticação JWT** - Login/logout seguro com validação de tokens
- **👥 Gestão de Usuários** - CRUD completo com criptografia de senhas
- **📦 Gestão de Produtos** - Controle de catálogo de produtos
- **🛒 Gestão de Compras** - Processamento de compras e pedidos
- **📋 Gestão de Pedidos** - Sistema completo de orders com status
- **🔒 Validação de Headers** - Controle organizacional via `X-Lpe-Organization-Id` e `X-Lpe-Project-Id`
- **🗑️ Soft Delete** - Remoção lógica mantendo histórico
- **📊 Logs de Auditoria** - Rastreamento completo de operações

---

## Instalação e Execução

### Pré-requisitos

- **Go 1.23.0+**
- **PostgreSQL 15+**
- **Git**

### Dependências do Projeto

- [Gin Web Framework](https://github.com/gin-gonic/gin) - Framework HTTP
- [GORM](https://gorm.io/) - ORM para Go
- [JWT-Go](https://github.com/golang-jwt/jwt) - Implementação JWT
- [bcrypt](https://golang.org/x/crypto/bcrypt) - Criptografia de senhas
- [Google UUID](https://github.com/google/uuid) - Geração de UUIDs

### Passos de Instalação

1. **Clone o repositório**:
   ```bash
   git clone <repository-url>
   cd LEP-Back
   ```

2. **Instale as dependências**:
   ```bash
   go mod tidy
   ```

3. **Configure o banco de dados**:
   - Configure as variáveis de ambiente para conexão com PostgreSQL
   - Execute as migrações necessárias

4. **Execute a aplicação**:
   ```bash
   go run main.go
   ```

5. **Teste a API**:
   ```bash
   curl http://localhost:8080/ping
   # Resposta esperada: "pong"
   ```

### Variáveis de Ambiente

```bash
DB_USER=seu_usuario_postgres
DB_PASS=sua_senha_postgres
DB_NAME=nome_do_banco
INSTANCE_UNIX_SOCKET=/caminho/para/socket (para GCP)
JWT_SECRET_PRIVATE_KEY=sua_chave_privada_jwt
JWT_SECRET_PUBLIC_KEY=sua_chave_publica_jwt
```

---

## API Endpoints

### Autenticação
```bash
POST   /login          # Login do usuário
POST   /logout         # Logout do usuário
POST   /checkToken     # Validar token JWT
```

### Usuários (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
GET    /user/:id       # Buscar usuário por ID
GET    /user/group/:id # Buscar usuários por grupo
POST   /user           # Criar usuário (público)
PUT    /user/:id       # Atualizar usuário
DELETE /user/:id       # Deletar usuário (soft delete)
```

### Produtos (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
GET    /product/:id           # Buscar produto por ID
GET    /product/purchase/:id  # Buscar produtos por compra
POST   /product              # Criar produto
PUT    /product/:id          # Atualizar produto
DELETE /product/:id          # Deletar produto
```

### Compras (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
GET    /purchase/:id       # Buscar compra por ID
GET    /purchase/group/:id # Buscar compras por grupo
POST   /purchase           # Criar compra
PUT    /purchase/:id       # Atualizar compra
DELETE /purchase/:id       # Deletar compra
```

### Pedidos (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
GET    /order/:id    # Buscar pedido por ID
GET    /orders       # Listar pedidos
POST   /order        # Criar pedido
PUT    /order/:id    # Atualizar pedido
DELETE /order/:id    # Deletar pedido
```

### Headers Obrigatórios (exceto /login e POST /user)
```bash
X-Lpe-Organization-Id: <organization-uuid>
X-Lpe-Project-Id: <project-uuid>
Authorization: Bearer <jwt-token>
```

---

## Deploy e Infraestrutura

### Deploy Local

```bash
# Build da aplicação
go build -o lep-system .

# Execução
./lep-system
```

### Deploy no Google Cloud Platform

O projeto inclui configuração completa do Terraform para deploy automatizado no GCP:

1. **Recursos provisionados**:
   - Cloud Run para a aplicação
   - Cloud SQL (PostgreSQL) para banco de dados
   - Secret Manager para chaves JWT
   - IAM roles e service accounts

2. **Deploy com Terraform**:
   ```bash
   # Inicializar Terraform
   terraform init

   # Planejar mudanças
   terraform plan

   # Aplicar infraestrutura
   terraform apply
   ```

3. **Configurações automáticas**:
   - Scaling automático (até 2 instâncias)
   - Conexão segura com Cloud SQL via Unix socket
   - Gerenciamento seguro de secrets
   - APIs necessárias habilitadas automaticamente

## Segurança e Auditoria

### Recursos de Segurança

- **🔐 Autenticação JWT** - Tokens seguros com chaves RSA
- **🔒 Criptografia bcrypt** - Senhas hashadas com salt
- **🛡️ Validação de Headers** - Controle organizacional obrigatório
- **🚫 Soft Delete** - Preservação de dados para auditoria
- **📊 Logs Detalhados** - Rastreamento completo de operações

### Auditoria

Todas as operações são registradas com:
- Timestamp da operação
- Usuário responsável
- Tipo de ação executada
- Entidade afetada
- Dados antes/depois da operação

## Documentação Adicional

- **[API Routes Documentation](routes/routes.md)** - Documentação completa dos endpoints
- **[Terraform Configuration](example.main.tf)** - Configuração de infraestrutura GCP

## Contribuição

1. Fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

