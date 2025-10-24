# 🏗️ Arquitetura de Seeds - Como Funciona

Explicação detalhada de como o sistema de seeds funciona, especialmente a relação entre seeds e acesso de usuários.

## 🔑 Resposta Rápida

**Sua Pergunta**: "Ao rodar o seed de Fattoria isso será cadastrado como uma nova organização? O usuário master admin continua com acesso para transitar nessa org?"

**Resposta Curta**:
- ✅ **SIM**, Fattoria é cadastrada como uma **NOVA organização**
- ✅ **SIM**, você pode transitar entre organizações
- ⚠️ **MAS**: Cada seed cria um **novo usuário admin específico** para aquela organização
- ⚠️ Os **master admins antigos NÃO têm acesso automático** à nova organização

---

## 🔍 Entendendo o Fluxo

### 1. Quando você roda um Seed

```bash
bash scripts/run_seed.sh           # LEP Demo (nova org)
# OU
bash tools/scripts/seed/run_seed_fattoria.sh  # Fattoria (nova org)
```

### 2. O que acontece no Bootstrap

```
1. Cria NOVA Organização
   ↓
2. Cria Projeto padrão dentro dessa Org
   ↓
3. Cria um NOVO usuário admin específico para essa Org
   ↓
4. Cria relacionamento: usuário → organização
   ↓
5. Cria relacionamento: usuário → projeto
```

### 3. Estrutura de Dados Criada

```
Organization (Nova!)
├── name: "Fattoria Pizzeria"
├── email: "fattoria@lep.com"
├── id: uuid_gerado_dinamicamente
│
└── Project
    ├── name: "Projeto Fattoria"
    ├── id: uuid_gerado_dinamicamente
    │
    └── User (NOVO ADMIN)
        ├── name: "Fattoria Pizzeria"
        ├── email: "fattoria@lep.com"
        ├── permissions: ["admin"]
        │
        ├── UserOrganization (owner)
        │
        └── UserProject (admin)
```

---

## 🚨 Ponto Crítico: IDs Dinâmicos

**IMPORTANTE**: O seed gera IDs **dinamicamente** via bootstrap!

```go
// No bootstrap (handler/organization.go)
org := &models.Organization{
    Id: uuid.New(),  // ← ID gerado ALEATORIAMENTE!
    Name: name,
}
```

Isso significa:
- ❌ Os IDs NO seed_fattoria.go (`FattoriaOrgID`, etc.) são **IGNORADOS**!
- ✅ O bootstrap cria IDs completamente novos cada vez
- ✅ Você pode rodar o seed múltiplas vezes e terá múltiplas orgs

---

## 👥 Acesso de Usuários

### Usuário Master Admin Antigo
```
Exemplo: pablo@lep.com (criado antes)
├── Acesso: Organização LEP Demo (original)
└── ❌ NÃO tem acesso a Fattoria (nova org)
```

### Novo Admin Fattoria
```
Criado pelo seed:
├── Email: fattoria@lep.com
├── Senha: senha123
├── Permissões: ["admin"]
├── Acesso: Organização Fattoria
└── ✅ Pode fazer login e transitar nessa org
```

---

## 📊 Cenário 1: Duas Orgs na Mesma DB

```
Banco de Dados
│
├── Organization: "LEP Demo"
│   ├── User: admin@lep-demo.com
│   │   ├── UserOrganization (admin)
│   │   └── UserProject (admin)
│   └── Project: "Projeto Principal"
│
└── Organization: "Fattoria Pizzeria"
    ├── User: fattoria@lep.com
    │   ├── UserOrganization (owner)
    │   └── UserProject (admin)
    └── Project: "Projeto Fattoria"
```

**Acesso via Headers**:
```bash
# Login LEP Demo
POST /login → email: admin@lep-demo.com

# Requisição com dados LEP
GET /product
  X-Lpe-Organization-Id: <LEP_ORG_ID>
  X-Lpe-Project-Id: <LEP_PROJECT_ID>

# Login Fattoria
POST /login → email: fattoria@lep.com

# Requisição com dados Fattoria
GET /product
  X-Lpe-Organization-Id: <FATTORIA_ORG_ID>
  X-Lpe-Project-Id: <FATTORIA_PROJECT_ID>
```

---

## 🔐 Sistema de Permissões Multi-Tenant

### Headers Obrigatórios (exceto /login)
```
X-Lpe-Organization-Id: <org_id>
X-Lpe-Project-Id: <project_id>
Authorization: Bearer <token>
```

### Validação do Backend
```go
// middleware/headers.go
func HeaderValidationMiddleware() gin.HandlerFunc {
    // 1. Valida se organização existe
    // 2. Valida se projeto existe
    // 3. Valida se usuário tem acesso a essa org
    // 4. Retorna dados apenas dessa org
}
```

---

## 📋 Código: Como o Bootstrap Funciona

### Arquivo: `handler/organization.go` - Função: `CreateOrganizationBootstrap()`

```go
func (r *resourceOrganization) CreateOrganizationBootstrap(
    name, password string,
) (*OrganizationBootstrapResponse, error) {

    // 1. Validar senha hard-coded
    if password != "senha123" {
        return nil, errors.New("senha inválida")
    }

    // 2. Criar Organização com ID dinâmico
    org := &models.Organization{
        Id:    uuid.New(),  // ← ID NOVO!
        Name:  name,
        Email: fmt.Sprintf("%s@lep.com", name),
    }
    r.repo.Organizations.CreateOrganization(org)

    // 3. Criar Projeto
    project := &models.Project{
        Id:             uuid.New(),  // ← ID NOVO!
        OrganizationId: org.Id,
        Name:           fmt.Sprintf("Projeto %s", name),
    }
    r.repo.Projects.CreateProject(project)

    // 4. Criar Usuário Admin
    user := &models.User{
        Id:          uuid.New(),  // ← ID NOVO!
        Name:        name,
        Email:       fmt.Sprintf("%s@lep.com", name),
        Password:    bcrypt.HashPassword("senha123"),
        Permissions: []string{"admin"},
    }
    r.repo.User.CreateUser(user)

    // 5. Criar Relacionamentos
    userOrg := &models.UserOrganization{
        UserId:         user.Id,
        OrganizationId: org.Id,
        Role:           "owner",  // ← DONO da organização!
    }
    r.repo.UserOrganizations.Create(userOrg)

    userProj := &models.UserProject{
        UserId:    user.Id,
        ProjectId: project.Id,
        Role:      "admin",
    }
    r.repo.UserProjects.Create(userProj)

    return &OrganizationBootstrapResponse{
        Organization: org,
        Project:      project,
        User:         user,
    }, nil
}
```

### Arquivo: `cmd/seed/main.go` - Função: `seedDatabaseViaServer()`

```go
func seedDatabaseViaServer(router *gin.Engine, data *utils.SeedData) error {
    // 1. BOOTSTRAP: Criar nova Org via POST /create-organization
    orgId, projectId, adminEmail, err := createOrganizationBootstrap(
        router,
        data.Organizations[0].Name,  // "Fattoria Pizzeria"
    )
    // Resposta:
    // - orgId: uuid novo
    // - projectId: uuid novo
    // - adminEmail: "fattoria@lep.com"

    // 2. LOGIN: Fazer login com o novo admin
    adminToken, err := loginUser(router, adminEmail, "senha123")

    // 3. AUTORIZAR: Configurar headers
    headers := map[string]string{
        "Authorization":         "Bearer " + adminToken,
        "X-Lpe-Organization-Id": orgId.String(),    // ← ID da NOVA org
        "X-Lpe-Project-Id":      projectId.String(), // ← ID do NOVO project
    }

    // 4. POPULAR: Criar dados (produtos, mesas, etc)
    // Todos os endpoints usam os headers acima
    // Logo, todos os dados são criados nessa ORG específica

    for _, product := range data.Products {
        createProduct(router, product, headers)  // ← headers com org_id da Fattoria
    }
}
```

---

## 🔄 Fluxo Completo de um Seed

```
1. Usuario roda: bash run_seed_fattoria.sh
                        ↓
2. Script chama: go run cmd/seed/main.go --restaurant=fattoria
                        ↓
3. main.go executa: seedDatabaseViaServer(router, GenerateFattoriaData())
                        ↓
4. Bootstrap:
   POST /create-organization
   {"name": "Fattoria Pizzeria", "password": "senha123"}

   ↓ Resposta:
   {
     "data": {
       "organization": {
         "id": "550e8400-e29b-41d4-a716-446655440000",  ← NOVO ID!
         "name": "Fattoria Pizzeria"
       },
       "project": {
         "id": "660e8400-e29b-41d4-a716-446655440000"   ← NOVO ID!
       },
       "user": {
         "email": "fattoria@lep.com"
       }
     }
   }
                        ↓
5. Login:
   POST /login
   {"email": "fattoria@lep.com", "password": "senha123"}

   ↓ Resposta:
   {"token": "eyJhbGciOi..."}
                        ↓
6. Criar dados com headers:
   X-Lpe-Organization-Id: 550e8400-e29b-41d4-a716-446655440000
   X-Lpe-Project-Id: 660e8400-e29b-41d4-a716-446655440000
   Authorization: Bearer eyJhbGciOi...

   POST /product, POST /table, etc.
                        ↓
7. Resultado:
   Nova organização criada
   Todos os dados da Fattoria nessa nova org
   Novo admin com permissões nessa org
```

---

## 🎯 Casos de Uso

### Caso 1: Testar LEP + Fattoria na Mesma DB

```bash
# 1. Seed LEP Demo
bash scripts/run_seed.sh
# Cria: Organization "LEP Demo", User "admin@lep-demo.com"

# 2. Seed Fattoria
bash tools/scripts/seed/run_seed_fattoria.sh
# Cria: Organization "Fattoria", User "fattoria@lep.com"

# Resultado: 2 Organizações na mesma DB
# - LEP Demo: Acesso com admin@lep-demo.com
# - Fattoria: Acesso com fattoria@lep.com
```

### Caso 2: Resetar Apenas Fattoria

```bash
# Limpar dados
bash tools/scripts/seed/run_seed_fattoria.sh --clear-first

# Resultado:
# - LEP Demo: Mantida (se não usou --clear-first no LEP)
# - Fattoria: Resetada com NOVO org_id
```

### Caso 3: Usar Master Admin em Múltiplas Orgs

**Atualmente**: NÃO é possível com o seed atual

Se precisar que um usuário tenha acesso a múltiplas orgs:
```go
// Seria necessário criar relacionamentos adicionais:
// UserOrganization: pablo@lep.com → Fattoria Org

// Isso seria feito programaticamente após o seed
POST /user-organization/user/<user_id>
{
  "user_id": "<pablo_id>",
  "organization_id": "<fattoria_org_id>",
  "role": "admin"
}
```

---

## 📝 Resumo Importante

| Aspecto | Detalhe |
|---------|---------|
| **Nova Org** | ✅ SIM, cada seed cria nova |
| **IDs Dinâmicos** | ✅ SIM, não usa IDs do code |
| **Novo Admin** | ✅ SIM, específico da org |
| **Master Acesso** | ⚠️ NÃO automático (precisa relacionamento) |
| **Múltiplas Orgs** | ✅ SIM, possível na mesma DB |
| **Trocar de Org** | ✅ SIM, via diferentes headers |

---

## 🔗 Documentação Relacionada

- [docs_seeds/README.md](../../docs_seeds/README.md) - Seeds disponíveis
- [docs_seeds/fattoria/SEED_FATTORIA.md](../../docs_seeds/fattoria/SEED_FATTORIA.md) - Detalhes Fattoria
- [docs/SETUP.md](../SETUP.md) - Setup do ambiente

---

**Última Atualização**: 2024
**Clareza**: ✅ Respondidas todas as perguntas
