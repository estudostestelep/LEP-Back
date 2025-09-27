# LEP System - Master Interactive Script

O **Master Interactive Script** é uma interface unificada que consolida todos os scripts e funcionalidades do LEP System em um único ponto de controle intuitivo e fácil de usar.

## 🚀 Início Rápido

```bash
# Modo interativo (recomendado)
./scripts/master-interactive.sh

# Modo batch (comandos diretos)
./scripts/master-interactive.sh --help
./scripts/master-interactive.sh --setup
./scripts/master-interactive.sh --status
```

## 🎛️ Funcionalidades

### 1. 🏠 Desenvolvimento Local
- **Iniciar servidor**: `go run main.go` com verificações automáticas
- **Build da aplicação**: Compilação com validações
- **Docker local**: Build e execução de containers
- **Health check**: Verificação de status da API
- **Limpeza de artifacts**: Remoção de builds e cache

### 2. ⚙️ Setup & Configuração
- **Setup completo**: Configuração automática do ambiente
- **Verificação de dependências**: Go, Docker, gcloud, Terraform
- **Geração de chaves JWT**: Chaves criptográficas seguras
- **Configuração Google Cloud**: Autenticação e projeto
- **Validação completa**: Verificação de todo o ambiente

### 3. 🌱 Database & Seeding
- **Popular com dados demo**: Dados realistas para desenvolvimento
- **Limpeza e repopulação**: Reset completo da database
- **Ambientes específicos**: dev, test, staging
- **Status da database**: Verificação de conectividade
- **Usuários demo**: admin@lep-demo.com, garcom@lep-demo.com, gerente@lep-demo.com

### 4. 🧪 Testes
- **Execução completa**: Todos os testes do projeto
- **Cobertura de código**: Relatórios detalhados de coverage
- **Relatório HTML**: Visualização interativa da cobertura
- **Testes específicos**: Execução de testes individuais
- **Modo verbose**: Saída detalhada para debugging

### 5. ☁️ Deploy GCP
- **Deploy interativo**: Multi-ambiente (dev, stage, prod)
- **Deploy rápido**: Processo otimizado e híbrido
- **Bootstrap GCP**: Criação inicial de recursos
- **Deploy infraestrutura**: Apenas Terraform
- **Deploy aplicação**: Apenas Cloud Run
- **Gerenciamento de segredos**: Atualização via Secret Manager

### 6. 🛠️ Utilitários
- **Verificação de dependências**: Status completo do sistema
- **Limpeza completa**: Reset de todos os caches e builds
- **Backup de configurações**: Preservação de arquivos importantes
- **Status do projeto**: Visão geral completa
- **Troubleshooting**: Diagnóstico automático de problemas
- **Relatório de ambiente**: Documentação técnica completa

### 7. ❓ Ajuda
- **Guia de primeiros passos**: Tutorial completo para iniciantes
- **Comandos úteis**: Referência rápida de comandos
- **Solução de problemas**: Troubleshooting comum
- **Links úteis**: Documentação oficial das ferramentas

## 🎯 Modo Batch (Command Line)

```bash
# Comandos diretos sem interface interativa
./scripts/master-interactive.sh --setup          # Setup completo
./scripts/master-interactive.sh --seed           # Popular database
./scripts/master-interactive.sh --test           # Executar testes
./scripts/master-interactive.sh --quick-deploy   # Deploy rápido
./scripts/master-interactive.sh --status         # Status do projeto
./scripts/master-interactive.sh --clean          # Limpeza completa
./scripts/master-interactive.sh --help           # Ajuda
```

## 📋 Scripts Consolidados

O Master Script integra e substitui todos os scripts individuais:

| Script Original | Funcionalidade | Localização no Master |
|----------------|----------------|---------------------|
| `setup.sh` | Setup completo do ambiente | Menu 2 → 1 |
| `local-dev.sh` | Desenvolvimento local | Menu 1 → * |
| `run_seed.sh` | Database seeding | Menu 3 → * |
| `run_tests.sh` | Execução de testes | Menu 4 → * |
| `deploy-interactive.sh` | Deploy multi-ambiente | Menu 5 → 1 |
| `quick-deploy.sh` | Deploy rápido | Menu 5 → 2 |
| `bootstrap-gcp.sh` | Bootstrap GCP | Menu 5 → 3 |

## 🔧 Primeiro Uso

### 1. Configuração Inicial (primeira vez)
1. Execute: `./scripts/master-interactive.sh`
2. Selecione **Menu 2 → 1** (Setup completo)
3. Selecione **Menu 2 → 3** (Gerar chaves JWT)
4. Selecione **Menu 2 → 6** (Validar configuração)

### 2. Desenvolvimento Diário
1. Selecione **Menu 3 → 1** (Popular database)
2. Selecione **Menu 1 → 1** (Iniciar servidor)
3. Acesse: http://localhost:8080/health

### 3. Deploy para Produção
1. Selecione **Menu 5 → 3** (Bootstrap GCP)
2. Selecione **Menu 5 → 1** (Deploy interativo)
3. Escolha o ambiente desejado

## 🔑 Credenciais Padrão

Após executar o seeding da database:

- **Admin**: admin@lep-demo.com / password
- **Garçom**: garcom@lep-demo.com / password
- **Gerente**: gerente@lep-demo.com / password

## 📊 Verificações Automáticas

O script inclui verificações automáticas para:

- ✅ Dependências do sistema (Go, Docker, gcloud, etc.)
- ✅ Arquivos de configuração (.env, terraform.tfvars)
- ✅ Chaves JWT e certificados
- ✅ Conectividade de rede e serviços
- ✅ Status da database e aplicação
- ✅ Build e testes do código

## 🩺 Troubleshooting Automático

**Menu 6 → 5** executa diagnóstico completo:

- Verifica instalação de ferramentas
- Testa conectividade de rede
- Valida configurações
- Identifica conflitos de porta
- Corrige permissões de arquivos
- Sugere soluções para problemas encontrados

## 🎨 Interface e Experiência

- **Cores consistentes**: Códigos de cor padronizados para melhor UX
- **Navegação intuitiva**: Menus numerados e hierárquicos
- **Feedback em tempo real**: Status e progresso das operações
- **Tratamento de erros**: Mensagens claras e sugestões de solução
- **Confirmações de segurança**: Para operações destrutivas
- **Interrupção graceful**: Ctrl+C tratado adequadamente

## 🔐 Segurança

- Não armazena credenciais em texto plano
- Backup automático antes de operações destrutivas
- Validações de entrada do usuário
- Confirmações para operações críticas
- Logs detalhados de operações

## 🚀 Performance

- **Execução paralela**: Comandos independentes executados em paralelo
- **Cache inteligente**: Reutilização de builds e downloads
- **Verificações otimizadas**: Skip de verificações desnecessárias
- **Feedback imediato**: Respostas rápidas para ações do usuário

## 🛠️ Personalização

O script pode ser customizado editando as variáveis no topo do arquivo:

```bash
# Global Configuration
PROJECT_ID="leps-472702"
PROJECT_NAME="leps"
REGION="us-central1"
```

## 📞 Suporte

Para problemas ou sugestões:

1. Execute **Menu 6 → 5** (Troubleshooting)
2. Execute **Menu 6 → 6** (Gerar relatório do ambiente)
3. Consulte **Menu 7** (Ajuda completa)

---

**LEP System Master Interactive Script v1.0.0**
*Uma ferramenta unificada para gerenciar todo o ciclo de vida do LEP System*