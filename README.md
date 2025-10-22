# LEP System - Backend API

## Visão Geral

O LEP System é uma aplicação backend robusta desenvolvida em Go, utilizando arquitetura limpa e modular. O sistema foi projetado para gestão completa de operações empresariais, incluindo usuários, produtos, compras e pedidos.

### Tecnologias Utilizadas

- **Go 1.24.0** - Linguagem de programação principal
- **Gin Web Framework** - Framework HTTP para APIs RESTful
- **GORM** - ORM para manipulação de banco de dados
- **PostgreSQL** - Banco de dados principal
- **Google Cloud Platform** - Infraestrutura de nuvem (Cloud Run, Cloud SQL, Secret Manager)
- **Terraform** - Infrastructure as Code
- **JWT** - Autenticação e autorização
- **bcrypt** - Criptografia de senhas
- **Twilio** - SMS e WhatsApp (API de notificações)
- **SMTP** - Email (sistema de notificações)

### Características Principais

- ✅ **Arquitetura Limpa** - Separação clara entre camadas (Handler, Service, Repository)
- ✅ **API RESTful** - Endpoints padronizados seguindo convenções REST
- ✅ **Autenticação JWT** - Sistema seguro de autenticação e autorização
- ✅ **Validação de Headers** - Controle organizacional via headers obrigatórios
- ✅ **Soft Delete** - Remoção lógica de registros para auditoria
- ✅ **CRUD Completo** - Operações completas para todas as entidades
- ✅ **Logs de Auditoria** - Rastreamento completo de operações
- ✅ **Deploy Automatizado** - Configuração Terraform para GCP
- ✅ **Sistema de Notificações** - SMS, WhatsApp e Email automatizados
- ✅ **Cron Jobs** - Confirmações 24h e processamento de eventos
- ✅ **Relatórios Avançados** - Analytics de ocupação, reservas e waitlist

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

- **🔐 Autenticação JWT** - Login/logout seguro com validação de tokens e blacklist
- **👥 Gestão de Usuários** - CRUD completo com criptografia bcrypt
- **📦 Gestão de Produtos** - Controle de catálogo de produtos
- **🛒 Gestão de Compras** - Processamento de compras e pedidos
- **📋 Gestão de Pedidos** - Sistema completo de orders com status
- **🏠 Gestão de Mesas** - Controle de mesas e disponibilidade
- **⏳ Lista de Espera** - Sistema de fila para mesas ocupadas
- **📅 Reservas** - Agendamento de mesas com controle de horários
- **👤 Gestão de Clientes** - CRUD completo de clientes
- **🔒 Validação de Headers** - Controle organizacional via `X-Lpe-Organization-Id` e `X-Lpe-Project-Id`
- **🗑️ Soft Delete** - Remoção lógica mantendo histórico
- **📊 Logs de Auditoria** - Rastreamento completo de operações
- **📱 Notificações Automatizadas** - SMS, WhatsApp e Email com templates dinâmicos
- **⏰ Cron Jobs** - Confirmações 24h antes das reservas e processamento de eventos
- **📈 Sistema de Relatórios** - Analytics de ocupação, reservas, waitlist e leads

---

## Instalação e Execução

### Pré-requisitos

- **Go 1.24.0+**
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
# Database
DB_USER=seu_usuario_postgres
DB_PASS=sua_senha_postgres
DB_NAME=nome_do_banco
INSTANCE_UNIX_SOCKET=/caminho/para/socket # Para GCP Cloud SQL

# Autenticação
JWT_SECRET_PRIVATE_KEY=sua_chave_privada_jwt
JWT_SECRET_PUBLIC_KEY=sua_chave_publica_jwt

# Twilio (SMS/WhatsApp)
TWILIO_ACCOUNT_SID=seu_account_sid_twilio
TWILIO_AUTH_TOKEN=seu_auth_token_twilio
TWILIO_PHONE_NUMBER=+5511999999999

# SMTP (Email)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=seu_email@gmail.com
SMTP_PASSWORD=sua_senha_app

# Cron Jobs (opcional)
ENABLE_CRON_JOBS=true
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

### Mesas (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
GET    /table/:id    # Buscar mesa por ID
GET    /table        # Listar mesas
POST   /table        # Criar mesa
PUT    /table/:id    # Atualizar mesa
DELETE /table/:id    # Deletar mesa
```

### Lista de Espera (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
GET    /waitlist/:id # Buscar entrada na lista por ID
GET    /waitlist     # Listar entradas da lista de espera
POST   /waitlist     # Criar entrada na lista de espera
PUT    /waitlist/:id # Atualizar entrada na lista de espera
DELETE /waitlist/:id # Remover da lista de espera
```

### Reservas (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
GET    /reservation/:id # Buscar reserva por ID
GET    /reservation     # Listar reservas
POST   /reservation     # Criar reserva
PUT    /reservation/:id # Atualizar reserva
DELETE /reservation/:id # Cancelar reserva
```

### Clientes (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
GET    /customer/:id # Buscar cliente por ID
GET    /customer     # Listar clientes
POST   /customer     # Criar cliente
PUT    /customer/:id # Atualizar cliente
DELETE /customer/:id # Deletar cliente
```

### Notificações (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
# Configuração de Notificações
POST   /notification/config          # Criar/atualizar configuração de evento
GET    /notification/config/:event   # Buscar configuração por evento

# Templates de Notificação
POST   /notification/template        # Criar template
PUT    /notification/template/:id    # Atualizar template
GET    /notification/templates       # Listar templates

# Envio Manual de Notificações
POST   /notification/send           # Enviar notificação manual

# Logs e Histórico
GET    /notification/logs           # Buscar logs de notificações
GET    /notification/logs/:id       # Buscar log específico

# Webhooks (para integração com Twilio)
POST   /notification/webhook/twilio/status    # Status de entrega SMS/WhatsApp
POST   /notification/webhook/twilio/inbound   # Mensagens recebidas
```

### Relatórios (Headers obrigatórios: X-Lpe-Organization-Id, X-Lpe-Project-Id)
```bash
# Relatórios Analíticos
GET    /reports/occupancy          # Relatório de ocupação de mesas
GET    /reports/reservations       # Relatório de reservas
GET    /reports/waitlist           # Relatório de lista de espera
GET    /reports/leads              # Relatório de leads (futuro)

# Exportação
GET    /reports/export/csv         # Exportar relatório em CSV
```

### Headers Obrigatórios (exceto /login e POST /user)
```bash
X-Lpe-Organization-Id: <organization-uuid>
X-Lpe-Project-Id: <project-uuid>
Authorization: Bearer <jwt-token>
```

---

## Sistema de Notificações

### Visão Geral

O LEP System inclui um sistema completo de notificações automatizadas que suporta:
- **SMS** via Twilio
- **WhatsApp Business** via Twilio API
- **Email** via SMTP (Gmail, Outlook, etc.)

### Configuração de Notificações

#### 1. Configuração de Eventos

Para configurar quais eventos irão disparar notificações:

```bash
POST /notification/config
```

**Payload:**
```json
{
  "event_type": "reservation_create",
  "enabled": true,
  "channels": ["sms", "whatsapp", "email"],
  "delay_minutes": 0
}
```

**Eventos Disponíveis:**
- `reservation_create` - Nova reserva criada
- `reservation_update` - Reserva atualizada
- `reservation_cancel` - Reserva cancelada
- `table_available` - Mesa disponível (waitlist)
- `confirmation_24h` - Confirmação 24h antes (automático)

#### 2. Criação de Templates

Para criar templates personalizados para cada canal:

```bash
POST /notification/template
```

**Payload:**
```json
{
  "channel": "sms",
  "event_type": "reservation_create",
  "subject": "Reserva Confirmada",
  "body": "Olá {{nome}}! Sua reserva para {{pessoas}} pessoas na mesa {{mesa}} está confirmada para {{data_hora}}. Até breve!"
}
```

**Variáveis Disponíveis:**
- `{{nome}}` ou `{{cliente}}` - Nome do cliente
- `{{mesa}}` ou `{{numero_mesa}}` - Número da mesa
- `{{data}}` - Data (DD/MM/YYYY)
- `{{hora}}` - Hora (HH:MM)
- `{{data_hora}}` - Data e hora completa
- `{{pessoas}}` - Quantidade de pessoas
- `{{tempo_espera}}` - Tempo estimado de espera
- `{{status}}` - Status da reserva

#### 3. Envio Manual de Notificações

Para enviar notificações pontuais:

```bash
POST /notification/send
```

**Payload:**
```json
{
  "event_type": "reservation_create",
  "entity_type": "reservation",
  "entity_id": "uuid-da-reserva",
  "recipient": "+5511999999999",
  "channel": "sms",
  "variables": {
    "nome": "João Silva",
    "mesa": "5",
    "data_hora": "25/12/2023 às 19:30"
  }
}
```

### Configuração de Webhooks

#### Twilio Webhooks

Para receber atualizações de status e mensagens inbound, configure os webhooks no Twilio:

**Status de Entrega:**
```
URL: https://seu-dominio.com/notification/webhook/twilio/status
Método: POST
```

**Mensagens Recebidas:**
```
URL: https://seu-dominio.com/notification/webhook/twilio/inbound
Método: POST
```

### Configuração do Projeto

Para habilitar notificações em um projeto específico, utilize as configurações:

```json
{
  "notify_reservation_create": true,
  "notify_reservation_update": true,
  "notify_reservation_cancel": true,
  "notify_table_available": true,
  "notify_confirmation_24h": true
}
```

### Logs e Monitoramento

Para acompanhar o envio de notificações:

```bash
GET /notification/logs?limit=50
```

**Resposta:**
```json
{
  "logs": [
    {
      "id": "uuid",
      "event_type": "reservation_create",
      "channel": "sms",
      "recipient": "+5511999999999",
      "status": "sent",
      "external_id": "twilio-message-id",
      "created_at": "2023-12-25T10:00:00Z",
      "delivered_at": "2023-12-25T10:00:05Z"
    }
  ]
}
```

**Status Possíveis:**
- `sent` - Enviado com sucesso
- `delivered` - Entregue ao destinatário
- `failed` - Falha no envio
- `pending` - Aguardando processamento

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

- **🔐 Autenticação JWT** - Tokens seguros com HS256 (24h expiração)
- **🔒 Criptografia bcrypt** - Senhas hashadas com salt automático
- **🚫 Token Blacklist** - Sistema de revogação via BannedLists
- **✅ Token Whitelist** - Controle de sessões ativas via LoggedLists
- **🛡️ Validação de Headers** - Controle organizacional obrigatório
- **🗑️ Soft Delete** - Preservação de dados para auditoria
- **📊 Logs Detalhados** - Rastreamento completo de operações

### ⚠️ Melhorias de Segurança Recomendadas

**Problemas Identificados:**
- Inconsistência nas chaves JWT (uso de chave pública/privada diferentes)
- Limpeza de tokens expirados apenas no logout
- Formato de data em string ao invés de timestamp

**Recomendações:**
- Padronizar uso de uma única chave secreta para HS256
- Implementar middleware de autenticação centralizado
- Adicionar refresh tokens para melhor UX
- Implementar rate limiting para login
- Job periódico para limpeza de tokens expirados

### Auditoria

Todas as operações são registradas com:
- Timestamp da operação
- Usuário responsável
- Tipo de ação executada
- Entidade afetada
- Dados antes/depois da operação

## Seeding de Dados

### Seed Local (Desenvolvimento)

Para popular o banco de dados local com dados de exemplo:

```bash
# Executar seed padrão
bash ./scripts/run_seed.sh

# Com verbose
bash ./scripts/run_seed.sh --verbose

# Limpar dados antes de popular
bash ./scripts/run_seed.sh --clear-first

# Especificar ambiente
bash ./scripts/run_seed.sh --environment=test
```

### Seed Remoto (Staging/Produção)

Para popular o banco de dados remoto via HTTP API (sem acesso direto ao PostgreSQL):

**Windows:**
```bash
# Seed para staging (padrão)
scripts\run_seed_remote.bat

# Com verbose
scripts\run_seed_remote.bat --verbose

# URL customizada
scripts\run_seed_remote.bat --url https://api.example.com --environment prod
```

**Linux/Mac:**
```bash
# Seed para staging (padrão)
bash ./scripts/run_seed_remote.sh

# Com verbose
bash ./scripts/run_seed_remote.sh --verbose

# URL customizada
bash ./scripts/run_seed_remote.sh --url https://api.example.com --environment prod
```

**Execução direta:**
```bash
# Build do binário
go build -o lep-seed-remote.exe cmd/seed-remote/main.go

# Executar
./lep-seed-remote.exe --url https://lep-system-516622888070.us-central1.run.app --verbose
```

### Credenciais após Seeding

**Master Admins (Acesso Total):**
- pablo@lep.com / senha123
- luan@lep.com / senha123
- eduardo@lep.com / senha123

**Demo Users:**
- teste@gmail.com / password (Admin)
- garcom1@gmail.com / password (Waiter)
- gerente1@gmail.com / password (Manager)

### Dados Gerados

O seed cria automaticamente:
- 1 organização e 1 projeto
- 6 usuários (3 master admins + 3 demo users)
- 12 produtos em 3 categorias
- 8 mesas com diferentes status
- 4 pedidos ativos
- 6 reservas (passadas, presentes, futuras)
- 3 entradas na lista de espera
- 5 clientes com preferências
- 5 templates de notificação

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

