# LEP Backend - Sistema de Testes

Este diretório contém testes abrangentes para o backend LEP, incluindo testes unitários, de integração e cenários de negócio.

## Estrutura dos Testes

### Arquivos Principais

- **`test_config.go`** - Configuração do ambiente de teste com banco real
- **`test_helpers.go`** - Utilitários para fazer requisições HTTP e validações
- **`test_data.go`** - Dados de exemplo para testes
- **`test_env.go`** - Configuração de variáveis de ambiente para testes

### Suites de Teste

1. **`unit_test.go`** - Testes unitários
   - Geração de dados de seed
   - Estruturas de dados de teste
   - Configuração de ambiente
   - Validação de dados

2. **`integration_test.go`** - Testes de integração básicos
   - Disponibilidade de rotas
   - Validação de autenticação
   - Verificação de headers obrigatórios

3. **`real_integration_test.go`** - Testes de integração reais
   - Operações CRUD completas
   - Interação com banco de dados real
   - Validação de entrada de dados
   - Isolamento multi-tenant

4. **`business_scenarios_test.go`** - Cenários de negócio
   - Fluxo completo de reservas
   - Criação de pedidos com múltiplos itens
   - Gerenciamento de lista de espera
   - Regras de validação de negócio

## Como Executar os Testes

### Comandos Básicos

```bash
# Executar todos os testes
go test ./tests -v

# Usar script completo
bash ./scripts/run_tests.sh

# Com cobertura de código
bash ./scripts/run_tests.sh --coverage

# Gerar relatório HTML
bash ./scripts/run_tests.sh --html

# Teste específico
bash ./scripts/run_tests.sh --test TestProductCRUD

# Com saída detalhada
bash ./scripts/run_tests.sh --verbose
```

### Executar Suites Específicas

```bash
# Apenas testes unitários
go test ./tests -run TestDataGeneration -v

# Apenas testes de integração real
go test ./tests -run RealIntegrationTestSuite -v

# Apenas cenários de negócio
go test ./tests -run BusinessScenariosTestSuite -v
```

## Configuração do Ambiente de Teste

### Pré-requisitos

1. **PostgreSQL** rodando localmente
2. **Banco de dados de teste** criado (nome: `lep_test`)
3. **Variáveis de ambiente** configuradas

### Variáveis de Ambiente Necessárias

```bash
# Banco de dados de teste
DB_NAME=lep_test
DB_USER=postgres
DB_PASS=postgres
DB_HOST=localhost
DB_PORT=5432

# Autenticação (para testes)
JWT_SECRET_PRIVATE_KEY=test-private-key
JWT_SECRET_PUBLIC_KEY=test-public-key

# Desabilitar serviços externos em testes
ENABLE_CRON_JOBS=false
ENABLE_NOTIFICATIONS=false
```

### Setup Automático

Os testes fazem setup automático do ambiente:

1. **Migração** - Auto-migra todas as tabelas necessárias
2. **Dados de Teste** - Popula com dados de exemplo para cada teste
3. **Limpeza** - Remove dados após cada teste para isolamento

## Tipos de Teste

### 1. Testes Unitários

Testam componentes isolados sem dependências externas:

- Geração de dados de seed
- Estruturas de dados
- Validações básicas
- Configuração de ambiente

### 2. Testes de Integração Básicos

Testam disponibilidade de rotas e autenticação:

- Endpoints públicos (`/ping`, `/health`)
- Validação de headers obrigatórios
- Proteção de rotas autenticadas

### 3. Testes de Integração Reais

Testam operações completas com banco de dados:

- **CRUD de Produtos** - Criar, ler, atualizar, deletar
- **CRUD de Clientes** - Operações completas
- **Gerenciamento de Mesas** - Status e capacidade
- **Validação de Dados** - Regras de negócio
- **Isolamento Multi-tenant** - Segregação por organização

### 4. Cenários de Negócio

Testam fluxos completos do sistema:

- **Fluxo de Reserva Completo**:
  1. Criar cliente
  2. Criar mesa
  3. Criar reserva
  4. Atualizar status
  5. Verificar dados

- **Criação de Pedido**:
  1. Criar produtos
  2. Criar cliente e mesa
  3. Criar pedido com múltiplos itens
  4. Atualizar status do pedido

- **Gerenciamento de Lista de Espera**:
  1. Adicionar cliente à lista
  2. Atualizar tempo estimado
  3. Marcar como atendido

## Dados de Teste

### Dados Pré-configurados

Os testes utilizam dados consistentes através do `utils.GenerateCompleteData()`:

- **1 Organização** - LEP Restaurante Demo
- **1 Projeto** - Projeto Principal
- **3 Usuários** - Admin, Garçom, Gerente
- **5 Clientes** - Com preferências diversas
- **12 Produtos** - Diferentes categorias e preços
- **8 Mesas** - Várias capacidades e localizações
- **Reservas e Pedidos** - Em diferentes estados

### IDs Consistentes

UUIDs fixos para relacionamentos:
- Organization ID: `123e4567-e89b-12d3-a456-426614174000`
- Project ID: `123e4567-e89b-12d3-a456-426614174001`
- User IDs: `123e4567-e89b-12d3-a456-426614174002-004`

## Validações Testadas

### Regras de Negócio

- **Produtos**: Preço positivo, nome obrigatório
- **Mesas**: Capacidade positiva, número único
- **Reservas**: Data futura, tamanho do grupo positivo
- **Clientes**: Email válido, telefone formatado
- **Pedidos**: Itens válidos, total calculado

### Segurança

- **Headers Obrigatórios**: Organization-Id e Project-Id
- **Autenticação**: Rotas protegidas requerem token
- **Isolamento**: Dados segregados por organização

## Solução de Problemas

### Erro: "Failed to connect to test database"

1. Verificar se PostgreSQL está rodando
2. Confirmar credenciais no `.env`
3. Criar banco `lep_test` manualmente

### Erro: "Failed to run migrations"

1. Verificar permissões do usuário do banco
2. Conferir se todas as dependências estão instaladas
3. Executar `go mod tidy`

### Testes Lentos

1. Usar `SKIP_DB_TESTS=true` para pular testes de banco
2. Executar apenas testes unitários
3. Verificar performance do PostgreSQL local

### Falhas Intermitentes

1. Verificar limpeza de dados entre testes
2. Conferir se há conflitos de concorrência
3. Executar testes individuais para isolar problemas

## Relatórios de Cobertura

```bash
# Gerar relatório de cobertura
bash ./scripts/run_tests.sh --coverage

# Relatório HTML interativo
bash ./scripts/run_tests.sh --html

# Visualizar no navegador (automático)
# Arquivo gerado: tests/reports/coverage.html
```

## Integração Contínua

Para usar em CI/CD:

```bash
# Setup de ambiente de teste
export DB_NAME=lep_test_ci
export DB_HOST=postgres-service
export SKIP_EXTERNAL_SERVICES=true

# Executar todos os testes
bash ./scripts/run_tests.sh --coverage --verbose
```

## Contribuindo

Ao adicionar novos testes:

1. **Seguir padrões** - Usar helpers existentes
2. **Isolamento** - Limpar dados após cada teste
3. **Documentar** - Explicar cenários complexos
4. **Validações** - Verificar todos os campos importantes
5. **Cobertura** - Incluir casos de erro e sucesso