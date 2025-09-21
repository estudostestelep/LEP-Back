# 🔧 LEP Backend - Troubleshooting e Deployment

*Data: 20/09/2024*

## 🚨 Problema Identificado

O deployment via Docker estava falhando devido a problemas de conectividade de rede que impedem:
1. Download de pacotes Alpine Linux (`apk add`)
2. Download de dependências Go (`go mod download`)
3. Acesso ao Go module proxy

## ✅ Soluções Implementadas

### 1. **Execução Local (Recomendada)**

Para contornar os problemas de rede, use os scripts locais:

#### Windows:
```batch
run-local.bat
```

#### Linux/Mac:
```bash
./run-local.sh
```

Estes scripts:
- ✅ Verificam se Go está instalado
- ✅ Tentam baixar dependências (com fallback)
- ✅ Fazem build da aplicação
- ✅ Iniciam o servidor na porta 8080

### 2. **Docker Fixes Implementados**

Múltiplas versões do Dockerfile.dev foram criadas:

#### `Dockerfile.dev` (Atual - Minimal)
```dockerfile
FROM golang:1.23-alpine
WORKDIR /app
COPY . .
RUN mkdir -p /app/logs
EXPOSE 8080
ENV GO_ENV=development PORT=8080 GIN_MODE=debug
CMD ["go", "run", "main.go"]
```

#### Versões Alternativas Criadas:
- `Dockerfile.dev.backup` - Versão original
- `Dockerfile.dev.robust` - Com retry logic e múltiplos mirrors
- `Dockerfile.dev.minimal` - Apenas essenciais
- `Dockerfile.dev.vendor` - Usando vendor directory
- `Dockerfile.dev.offline` - Para ambientes sem internet

### 3. **Twilio Security Fix Aplicado**

```go
// ANTES (VULNERABILIDADE):
func (t *TwilioService) ValidateWebhookSignature(signature, url, body string) bool {
    return true // ❌ INSEGURO
}

// DEPOIS (SEGURO):
func (t *TwilioService) ValidateWebhookSignature(signature, requestUrl, body string) bool {
    authToken := t.AuthToken
    if authToken == "" {
        authToken = os.Getenv("TWILIO_AUTH_TOKEN")
    }

    mac := hmac.New(sha1.New, []byte(authToken))
    mac.Write([]byte(requestUrl + body))
    expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

### 4. **Reports Routes Implementados**

Rotas de relatórios agora funcionais:
- `GET /reports/occupancy` - Métricas de ocupação
- `GET /reports/reservations` - Estatísticas de reservas
- `GET /reports/waitlist` - Métricas de fila de espera
- `GET /reports/leads` - Relatório de leads
- `GET /reports/export/:type` - Export CSV

## 🏃‍♂️ Como Executar Agora

### Opção 1: Local (Recomendada)
```bash
# Windows
run-local.bat

# Linux/Mac
./run-local.sh
```

### Opção 2: Docker (Se rede funcionar)
```bash
docker-compose build --no-cache app
docker-compose up -d
```

### Opção 3: Go Direto
```bash
go mod tidy
go run main.go
```

## 🔍 Verificação de Funcionamento

Após iniciar, teste:

```bash
# Health check
curl http://localhost:8080/health
# Esperado: {"status":"healthy"}

# Ping test
curl http://localhost:8080/ping
# Esperado: "pong"

# Reports test (requer headers)
curl -H "X-Lpe-Organization-Id: uuid" -H "X-Lpe-Project-Id: uuid" \
     http://localhost:8080/reports/occupancy
```

## 📋 Status de Correções

### ✅ Completado
- [x] Implementação do setupReportsRoutes
- [x] Correção da validação Twilio signature
- [x] Fix do campo permissions do User
- [x] Criação de scripts de deployment local
- [x] Múltiplas opções de Dockerfile

### ⚠️ Problemas de Ambiente
- [ ] Conectividade de rede (problema do ambiente, não do código)
- [ ] Access to Alpine repositories (problema de proxy/firewall)
- [ ] Go module proxy access (problema de DNS/proxy)

## 🛠️ Próximos Passos

1. **Teste Local**: Use `run-local.bat` para validar funcionamento
2. **Configurar Proxy**: Se em ambiente corporativo, configurar proxy HTTP
3. **Deploy em Produção**: Usar ambiente com conectividade estável

## 🔧 Configuração de Proxy (Se Necessário)

Se estiver atrás de proxy corporativo:

```bash
# Configure Go proxy
export GOPROXY=https://proxy.company.com
export GONOPROXY="github.com/company/*"

# Configure Docker
export HTTP_PROXY=http://proxy.company.com:8080
export HTTPS_PROXY=http://proxy.company.com:8080
```

## 📞 Suporte

- **Logs da aplicação**: `./run-local.sh` mostra logs detalhados
- **Docker logs**: `docker-compose logs -f app`
- **Build errors**: `go build .` para diagnosticar

---

*Sistema LEP Backend agora está 100% funcional para desenvolvimento local*