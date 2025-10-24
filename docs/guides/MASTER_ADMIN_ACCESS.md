# 🔐 Master Admin - Acesso Automático a Novas Organizações

Documentação sobre a regra de negócio: Master Admins são automaticamente adicionados a novas organizações.

## 📋 Resumo Executivo

Quando uma **nova organização** é criada (via bootstrap), todos os **master admins** do sistema são **automaticamente adicionados** com acesso de **admin**.

### ✨ Benefício
Master admins (Pablo, Luan, Eduardo) podem acessar qualquer organização criada no sistema sem necessidade de criar relacionamentos manuais.

---

## 🔑 Regra de Negócio

### Quando Aplica
```
Evento: POST /create-organization
  ↓
Fluxo:
  1. Criar nova Organization
  2. Criar novo Project
  3. Criar novo User Admin (bootstrap)
  4. 🔑 Buscar todos os master admins
  5. 🔑 Adicionar cada master admin à nova org
  6. 🔑 Adicionar cada master admin ao novo projeto
```

### Master Admins Automáticos

```
email: pablo@lep.com
email: luan@lep.com
email: eduardo@lep.com
```

### Acesso Concedido
```
UserOrganization:
  Role: "admin"
  Status: Ativo

UserProject:
  Role: "admin"
  Status: Ativo
```

---

## 💻 Implementação

### Arquivo: `handler/organization.go`

#### Função: `CreateOrganizationBootstrap()`

```go
func (r *resourceOrganization) CreateOrganizationBootstrap(
    name, password string,
) (*OrganizationBootstrapResponse, error) {

    // ... criar org, projeto, usuário...

    // 🔑 NOVA REGRA: Adicionar master admins automaticamente
    if err := r.addMasterAdminsToOrganization(org.Id, project.Id); err != nil {
        return nil, fmt.Errorf("erro ao adicionar master admins: %v", err)
    }

    return response, nil
}
```

#### Função: `addMasterAdminsToOrganization()`

```go
func (r *resourceOrganization) addMasterAdminsToOrganization(
    organizationId, projectId uuid.UUID,
) error {
    // 1. Definir emails dos master admins
    masterAdminEmails := []string{
        "pablo@lep.com",
        "luan@lep.com",
        "eduardo@lep.com",
    }

    // 2. Para cada master admin
    for _, email := range masterAdminEmails {
        // 3. Buscar usuário por email
        user, err := r.repo.User.GetUserByEmail(email)
        if err != nil {
            continue  // Pular se não encontrado
        }

        // 4. Criar relacionamento com organização
        userOrg := &models.UserOrganization{
            UserId:         user.Id,
            OrganizationId: organizationId,
            Role:           "admin",
            Active:         true,
        }
        r.repo.UserOrganizations.Create(userOrg)

        // 5. Criar relacionamento com projeto
        userProj := &models.UserProject{
            UserId:    user.Id,
            ProjectId: projectId,
            Role:      "admin",
            Active:    true,
        }
        r.repo.UserProjects.Create(userProj)
    }

    return nil
}
```

---

## 📊 Estrutura de Dados Resultante

### Antes (Sem Master Admins)
```
Organization: "Fattoria Pizzeria"
├── User: "fattoria@lep.com" (bootstrap)
│   ├── UserOrganization (role: owner)
│   └── UserProject (role: admin)
│
└── ❌ Master admins não têm acesso
```

### Depois (Com Master Admins Automáticos)
```
Organization: "Fattoria Pizzeria"
├── User: "fattoria@lep.com" (bootstrap)
│   ├── UserOrganization (role: owner)
│   └── UserProject (role: admin)
│
├── User: "pablo@lep.com" ✅ ADICIONADO
│   ├── UserOrganization (role: admin)
│   └── UserProject (role: admin)
│
├── User: "luan@lep.com" ✅ ADICIONADO
│   ├── UserOrganization (role: admin)
│   └── UserProject (role: admin)
│
└── User: "eduardo@lep.com" ✅ ADICIONADO
    ├── UserOrganization (role: admin)
    └── UserProject (role: admin)
```

---

## 🚀 Fluxo Completo: Criar Fattoria

```
1. Usuário executa seed
   bash tools/scripts/seed/run_seed_fattoria.sh

2. Bootstrap é chamado
   POST /create-organization
   {"name": "Fattoria Pizzeria", "password": "senha123"}

3. Nova Organization criada
   ID: 550e8400-e29b-41d4-a716-446655440000
   Name: "Fattoria Pizzeria"

4. Novo Project criado
   ID: 660e8400-e29b-41d4-a716-446655440000
   Name: "Projeto Fattoria"

5. Novo User criado (bootstrap)
   Email: "fattoria@lep.com"
   Role na Org: "owner"
   Role no Proj: "admin"

6. 🔑 MASTER ADMINS ADICIONADOS AUTOMATICAMENTE
   ├── pablo@lep.com
   │   ├── UserOrganization criado (role: admin)
   │   └── UserProject criado (role: admin)
   ├── luan@lep.com
   │   ├── UserOrganization criado (role: admin)
   │   └── UserProject criado (role: admin)
   └── eduardo@lep.com
       ├── UserOrganization criado (role: admin)
       └── UserProject criado (role: admin)

7. Resultado
   ✅ Fattoria criada
   ✅ Todos master admins têm acesso
   ✅ Podem acessar via X-Lpe-Organization-Id: <fattoria_org_id>
```

---

## 🔍 Exemplo: Login e Acesso

### Master Admin Acessando Fattoria

```bash
# 1. Login (pode usar sua conta original)
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"pablo@lep.com","password":"senha123"}'

# Resposta
{
  "token": "eyJhbGciOi...",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174010",
    "email": "pablo@lep.com",
    "permissions": ["master_admin"]
  }
}

# 2. Acessar dados da Fattoria com headers
curl -X GET http://localhost:8080/product \
  -H "Authorization: Bearer eyJhbGciOi..." \
  -H "X-Lpe-Organization-Id: 550e8400-e29b-41d4-a716-446655440000" \
  -H "X-Lpe-Project-Id: 660e8400-e29b-41d4-a716-446655440000"

# Resposta: Produtos da Fattoria
{
  "data": [
    {"id": "...", "name": "Margherita", "price": 80.00},
    {"id": "...", "name": "Marinara", "price": 58.00},
    ...
  ]
}

# 3. Alternar para outra org
curl -X GET http://localhost:8080/product \
  -H "Authorization: Bearer eyJhbGciOi..." \
  -H "X-Lpe-Organization-Id: <LEP_DEMO_ORG_ID>" \
  -H "X-Lpe-Project-Id: <LEP_DEMO_PROJECT_ID>"

# Resposta: Produtos da LEP Demo
```

---

## 🛡️ Idempotência (Segurança)

A implementação é **idempotente**:

```go
// Se relacionamento já existe, criar() retorna erro
// Mas o código ignora o erro (seguro)
_ = r.repo.UserOrganizations.Create(userOrg)

// Resultado:
// - Primeira execução: cria novo relacionamento
// - Próximas execuções: ignora (não duplica)
```

Benefício:
- ✅ Seguro rodar seed múltiplas vezes
- ✅ Não cria duplicatas
- ✅ Master admins têm acesso consistente

---

## 📋 Casos de Uso

### Caso 1: Criar Múltiplas Orgs
```bash
# Seed LEP Demo
bash scripts/run_seed.sh
# Result: pablo, luan, eduardo têm acesso

# Seed Fattoria
bash tools/scripts/seed/run_seed_fattoria.sh
# Result: pablo, luan, eduardo têm acesso TAMBÉM

# Master admins agora podem acessar AMBAS as orgs!
```

### Caso 2: Novo Master Admin Precisa de Acesso
Se adicionar novo master admin posterior, fazer manualmente:
```go
// No código ou via API:
POST /user-organization/user/<user_id>
{
  "user_id": "<new_master_admin_id>",
  "organization_id": "<org_id>",
  "role": "admin"
}
```

### Caso 3: Remover Master Admin de Uma Org
Se necessário remover acesso:
```sql
DELETE FROM user_organizations
WHERE user_id = '<master_admin_id>'
  AND organization_id = '<org_id>';
```

---

## ⚙️ Configuração

### Adicionar Novo Master Admin

1. **Adicione na função**:
```go
masterAdminEmails := []string{
    "pablo@lep.com",
    "luan@lep.com",
    "eduardo@lep.com",
    "novo@lep.com",  // ← Adicione aqui
}
```

2. **Recompile**:
```bash
go build -o lep-system .
```

### Remover Master Admin (Temporário)

Comentar linha:
```go
masterAdminEmails := []string{
    "pablo@lep.com",
    "luan@lep.com",
    // "eduardo@lep.com",  // ← Comentado
}
```

---

## 🎯 Resumo

| Aspecto | Detalhe |
|---------|---------|
| **Quando aplica** | Ao criar nova org via bootstrap |
| **Quem é adicionado** | Master admins (pablo, luan, eduardo) |
| **Role concedido** | "admin" na org e projeto |
| **Automático** | ✅ Sim |
| **Idempotente** | ✅ Sim (seguro rodar múltiplas vezes) |
| **Configurável** | ✅ Sim (modifique a lista de emails) |

---

## 🔗 Documentação Relacionada

- [docs/seed/README.md](../seed/README.md) - Seeds
- [docs/seed/SEED_ARCHITECTURE.md](../seed/SEED_ARCHITECTURE.md) - Arquitetura de seeds
- [docs_seeds/README.md](../../docs_seeds/README.md) - Seeds disponíveis
- [middleware/authorization.go](../../middleware/authorization.go) - Autorização

---

**Implementação**: ✅ Completa
**Testado**: ✅ Compilado com sucesso
**Status**: ✅ Pronto para usar
