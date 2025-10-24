# 🍕 Fattoria Pizzeria - IDs Reference

Referência completa de todos os IDs usados no seed da Fattoria para fácil integração e testes.

## 🏢 Organização e Projeto

| Entidade | ID | Valor |
|----------|----|----|
| **Organization** | `FattoriaOrgID` | `223e4567-e89b-12d3-a456-426614174100` |
| **Project** | `FattoriaProjectID` | `223e4567-e89b-12d3-a456-426614174101` |
| **Menu** | `FattoriaMenuID` | `223e4567-e89b-12d3-a456-426614174102` |

## 🏷️ Categorias Principais

| Categoria | ID |
|-----------|-----|
| Pizzas | `223e4567-e89b-12d3-a456-426614174110` |
| Bebidas | `223e4567-e89b-12d3-a456-426614174111` |

## 📂 Subcategorias

| Subcategoria | ID | Pai |
|--------------|-----|-----|
| Entradas | `223e4567-e89b-12d3-a456-426614174120` | Pizzas |
| Pizzas | `223e4567-e89b-12d3-a456-426614174121` | Pizzas |
| Soft drinks | `223e4567-e89b-12d3-a456-426614174122` | Bebidas |
| Cervejas | `223e4567-e89b-12d3-a456-426614174123` | Bebidas |
| Cervejas artesanais | `223e4567-e89b-12d3-a456-426614174124` | Bebidas |
| Coquetéis | `223e4567-e89b-12d3-a456-426614174125` | Bebidas |

## 🏷️ Tags

| Tag | ID | Cor | Descrição |
|-----|----|----|-----------|
| Vegetariana | `223e4567-e89b-12d3-a456-426614174130` | #4CAF50 | Prato vegetariano |
| Vegana | `223e4567-e89b-12d3-a456-426614174131` | #8BC34A | Prato vegano |

## 🍕 Produtos - Pizzas

### Entradas
| Produto | ID | Preço | Tempo | Tags |
|---------|-----|-------|-------|------|
| Crostini | `223e4567-e89b-12d3-a456-426614174200` | R$ 30,00 | 15min | - |

### Pizzas
| Produto | ID | Preço | Tempo | Tags |
|---------|-----|-------|-------|------|
| Marguerita | `223e4567-e89b-12d3-a456-426614174201` | R$ 80,00 | 25min | Vegetariana |
| Marinara | `223e4567-e89b-12d3-a456-426614174202` | R$ 58,00 | 25min | Vegana |
| Parma | `223e4567-e89b-12d3-a456-426614174203` | R$ 109,00 | 25min | - |
| Vegana | `223e4567-e89b-12d3-a456-426614174204` | R$ 60,00 | 25min | Vegana |

## 🥤 Produtos - Bebidas

| Produto | ID | Preço | Volume | Tempo | Tags |
|---------|-----|-------|--------|-------|------|
| Suco de caju integral | `223e4567-e89b-12d3-a456-426614174205` | R$ 15,00 | 300ml | 2min | - |
| Baden Baden IPA | `223e4567-e89b-12d3-a456-426614174206` | R$ 23,00 | 600ml | 2min | - |
| Sônia e Zé | `223e4567-e89b-12d3-a456-426614174207` | R$ 32,00 | - | 5min | - |
| Heineken s/ álcool | `223e4567-e89b-12d3-a456-426614174208` | R$ 13,00 | 330ml | 2min | - |

## 🪑 Mesas

| Mesa | ID | Lugares | Status | Localização |
|-----|-----|---------|--------|------------|
| Mesa 1 | `223e4567-e89b-12d3-a456-426614174800` | 4 | Livre | Salão Principal - Entrada |
| Mesa 2 | `223e4567-e89b-12d3-a456-426614174801` | 2 | Livre | Salão Principal - Janela |
| Mesa 3 | `223e4567-e89b-12d3-a456-426614174802` | 6 | Livre | Salão Principal - Fundo |

## 🏢 Ambiente

| Ambiente | ID | Capacidade |
|----------|-----|-----------|
| Salão Principal | `223e4567-e89b-12d3-a456-426614174300` | 60 |

## 👥 Usuários

### Admin Fattoria
- **ID**: `223e4567-e89b-12d3-a456-426614174310`
- **Email**: admin@fattoria.com.br
- **Senha**: password
- **Permissões**: admin, products, orders, reservations, customers, tables, reports

### Relacionamentos
- **UserOrganization**: `223e4567-e89b-12d3-a456-426614174410`
- **UserProject**: `223e4567-e89b-12d3-a456-426614174510`

## 🔧 Como Usar os IDs

### Em Testes
```go
import "lep/utils"

// Usar constantes diretamente
productID := utils.FattoriaProductMargueritaID
orgID := utils.FattoriaOrgID
```

### Em Queries SQL
```sql
-- Buscar todos os produtos da Fattoria
SELECT * FROM products
WHERE organization_id = '223e4567-e89b-12d3-a456-426614174100'
  AND project_id = '223e4567-e89b-12d3-a456-426614174101';

-- Buscar apenas pizzas
SELECT * FROM products
WHERE category_id = '223e4567-e89b-12d3-a456-426614174121'
  AND organization_id = '223e4567-e89b-12d3-a456-426614174100';
```

### Em Requisições API
```bash
# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@fattoria.com.br",
    "password": "password"
  }'

# Obter produtos da Fattoria
curl -X GET http://localhost:8080/product \
  -H "Authorization: Bearer {token}" \
  -H "X-Lpe-Organization-Id: 223e4567-e89b-12d3-a456-426614174100" \
  -H "X-Lpe-Project-Id: 223e4567-e89b-12d3-a456-426614174101"

# Criar mesa
curl -X POST http://localhost:8080/table \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -H "X-Lpe-Organization-Id: 223e4567-e89b-12d3-a456-426614174100" \
  -H "X-Lpe-Project-Id: 223e4567-e89b-12d3-a456-426614174101" \
  -d '{
    "number": 4,
    "capacity": 4,
    "status": "livre",
    "environment_id": "223e4567-e89b-12d3-a456-426614174300"
  }'
```

## 📋 Padrão de IDs

Todos os IDs da Fattoria seguem um padrão:
```
223e4567-e89b-12d3-a456-426614174XXX
                          ↑
                     Intervalo Fattoria
```

Intervalos:
- **100-102**: Org, Project, Menu
- **110-111**: Categorias principais
- **120-125**: Subcategorias
- **130-131**: Tags
- **200-208**: Produtos
- **300-399**: Ambientes
- **310-399**: Usuários e relacionamentos
- **400-599**: Relacionamentos de usuários
- **800-899**: Mesas
- **900-999**: Settings, Templates

## 🔍 Como Encontrar um ID

Se precisar encontrar um ID específico:

1. **Em `utils/seed_fattoria.go`**:
   ```go
   var (
       FattoriaProductMargueritaID = uuid.MustParse("...")
   )
   ```

2. **Neste arquivo**:
   - Procure pela entidade na tabela correspondente

3. **No banco de dados**:
   ```sql
   SELECT * FROM products WHERE name = 'Marguerita';
   ```

## ✅ Checklist de Integração

- [ ] Seed executado com sucesso: `bash scripts/run_seed_fattoria.sh`
- [ ] Servidor iniciado: `go run main.go`
- [ ] API respondendo: `curl http://localhost:8080/health`
- [ ] Login funcionando com credenciais da Fattoria
- [ ] Produtos visíveis na API: `GET /product`
- [ ] Mesas listadas: `GET /table`
- [ ] Tags aplicadas aos produtos

---

**Última atualização**: 2024
**Compatibilidade**: LEP v1.0+
**Manutenido em**: `utils/seed_fattoria.go`
