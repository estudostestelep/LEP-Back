# Solução: Seed de Produtos não está criando os produtos

## Problema Identificado

O seed da Fattoria Pizzeria (e potencialmente do seed padrão) **NÃO ESTÁ CRIANDO OS PRODUTOS** porque há um problema estrutural no código de seed:

### Causa Raiz

A estrutura `SeedData` em `utils/seed_data.go` estava **FALTANDO** dois campos:
- `Subcategories []models.Subcategory`
- `SubcategoryCategories []models.SubcategoryCategory`

No arquivo `utils/seed_fattoria.go`, os produtos estão sendo criados com `CategoryId` apontando para **SubcategoryId** (linhas 285, 301, 316, etc):

```go
CategoryId: &FattoriaSubcategoryEntradasID  // ← Subcategoria, não categoria!
```

Porém, as subcategorias estavam sendo criadas como `Category` (com `MenuId`), causando:
1. **Conflito de tipos**: Subcategorias com `MenuId` quando deveriam ser standalone
2. **Falta de relacionamentos**: Sem `SubcategoryCategory` linking, as subcategorias ficariam órfãs
3. **Falha na validação**: O handler de produtos recusa produtos sem categoria válida

### Por que o Script não mostra erro?

O script `run_seed_fattoria.sh` apenas reporta "sucesso" porque:
- O programa `go run cmd/seed/main.go` executa sem mostrar stdout
- A validação silencia erros de request (status 400-500)
- Não há log visível dos produtos sendo rejeitados

## Solução Implementada

### Passo 1: Atualizar SeedData struct

No arquivo `utils/seed_data.go`, adicione após `Categories`:

```go
type SeedData struct {
    // ... campos existentes ...
    Categories         []models.Category
    Subcategories          []models.Subcategory                // ← ADICIONADO
    SubcategoryCategories  []models.SubcategoryCategory        // ← ADICIONADO
    Tags               []models.Tag
    // ... resto dos campos ...
}
```

**Status**: ✅ JÁ IMPLEMENTADO (linhas 23-24 de seed_data.go)

### Passo 2: Adicionar campos ao seed padrão

No arquivo `utils/seed_data.go`, na função `GenerateCompleteData()`, adicione após `ProductTags`:

```go
Subcategories:        []models.Subcategory{},           // Vazio - seed padrão não usa
SubcategoryCategories: []models.SubcategoryCategory{},  // Vazio
```

**Status**: ✅ JÁ IMPLEMENTADO (linhas 654-655 de seed_data.go)

### Passo 3: Reestruturar seed_fattoria.go

No arquivo `utils/seed_fattoria.go`:

#### 3.1 Remover subcategorias de `Categories`

Mude de:
```go
Categories: []models.Category{
    // Pizzas
    { Id: FattoriaCategoryPizzasID, ... },
    // Bebidas  
    { Id: FattoriaCategoryBebidasID, ... },
    // Entradas (como category)
    { Id: FattoriaSubcategoryEntradasID, ... },  // ← REMOVER
    // Pizzas (como category)
    { Id: FattoriaSubcategoryPizzasID, ... },    // ← REMOVER
    // ... etc
}
```

Para:
```go
Categories: []models.Category{
    // Pizzas
    { Id: FattoriaCategoryPizzasID, ... },
    // Bebidas  
    { Id: FattoriaCategoryBebidasID, ... },
}
```

#### 3.2 Adicionar `Subcategories` separadamente

Após `Tags`, adicione:

```go
Subcategories: []models.Subcategory{
    {
        Id:             FattoriaSubcategoryEntradasID,
        OrganizationId: FattoriaOrgID,
        ProjectId:      FattoriaProjectID,
        Name:           "Entradas",
        Order:          1,
        Active:         true,
        CreatedAt:      now,
        UpdatedAt:      now,
    },
    // ... e assim para: Pizzas, Soft drinks, Cervejas, Cervejas artesanais, Coquetéis
}
```

#### 3.3 Adicionar `SubcategoryCategories` para linking

Após `Subcategories`, adicione:

```go
SubcategoryCategories: []models.SubcategoryCategory{
    {
        Id:            uuid.MustParse("223e4567-e89b-12d3-a456-426614174150"),
        SubcategoryId: FattoriaSubcategoryEntradasID,
        CategoryId:    FattoriaCategoryPizzasID,
    },
    {
        Id:            uuid.MustParse("223e4567-e89b-12d3-a456-426614174151"),
        SubcategoryId: FattoriaSubcategoryPizzasID,
        CategoryId:    FattoriaCategoryPizzasID,
    },
    {
        Id:            uuid.MustParse("223e4567-e89b-12d3-a456-426614174152"),
        SubcategoryId: FattoriaSubcategorySoftID,
        CategoryId:    FattoriaCategoryBebidasID,
    },
    // ... e assim para os outros relacionamentos
}
```

**Status**: 🔄 PARCIALMENTE IMPLEMENTADO
- ✅ Adicionado `Subcategories` (linhas 276-337)
- ✅ Adicionado `SubcategoryCategories` (linhas 339-370)
- ❌ **FALTA**: Remover subcategorias da seção `Categories`

### Passo 4: Adicionar handlers para criar subcategorias no seed

No arquivo `cmd/seed/main.go`, na função `seedDatabaseViaServer()`, adicione ANTES de criar produtos:

```go
// 9.5 Criar subcategories
if len(data.Subcategories) > 0 {
    fmt.Println("  📑 Criando subcategories...")
    for _, subcategory := range data.Subcategories {
        subcategory.OrganizationId = orgId
        subcategory.ProjectId = projectId
        if err := createSubcategory(router, subcategory, headers); err != nil {
            return fmt.Errorf("failed to create subcategory %s: %v", subcategory.Name, err)
        }
    }
}

// 9.6 Criar subcategory-category relationships
if len(data.SubcategoryCategories) > 0 {
    fmt.Println("  🔗 Criando subcategory-category relationships...")
    for _, sc := range data.SubcategoryCategories {
        if err := createSubcategoryCategory(router, sc, headers); err != nil {
            return fmt.Errorf("failed to create subcategory-category relationship: %v", err)
        }
    }
}
```

E adicione as funções:

```go
func createSubcategory(router *gin.Engine, subcategory models.Subcategory, headers map[string]string) error {
    body, _ := json.Marshal(subcategory)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/subcategory", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    router.ServeHTTP(w, req)
    
    if w.Code != 201 && w.Code != 409 {
        return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
    }
    if verbose {
        fmt.Printf("    ✓ %s\n", subcategory.Name)
    }
    return nil
}

func createSubcategoryCategory(router *gin.Engine, sc models.SubcategoryCategory, headers map[string]string) error {
    body, _ := json.Marshal(sc)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/subcategory-category", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    router.ServeHTTP(w, req)
    
    if w.Code != 201 && w.Code != 409 {
        return fmt.Errorf("status %d - %s", w.Code, w.Body.String())
    }
    return nil
}
```

**Status**: ❌ NÃO IMPLEMENTADO - Necessário implementar

## Checklist de Aplicação

- [x] Atualizar `SeedData` struct em `utils/seed_data.go`
- [x] Adicionar campos vazios ao seed padrão
- [x] Adicionar `Subcategories` e `SubcategoryCategories` ao seed_fattoria
- [ ] **FALTA**: Remover subcategorias da seção `Categories` em seed_fattoria.go (linhas 180-246)
- [ ] Adicionar funções `createSubcategory()` e `createSubcategoryCategory()` em `cmd/seed/main.go`
- [ ] Adicionar chamadas para criar subcategorias em `seedDatabaseViaServer()` (ANTES dos produtos)
- [ ] Testar o seed completo com `bash scripts/run_seed_fattoria.sh --verbose`

## Como Testar Depois

```bash
# 1. Limpar e fazer seed com verbose
cd LEP-Back
bash scripts/run_seed_fattoria.sh --clear-first --verbose

# 2. Conectar ao banco e verificar
PGPASSWORD=lep_password psql -h localhost -U lep_user -d lep_database -c "SELECT COUNT(*) FROM products;"

# 3. Esperado: 9 produtos (Crostini, Marguerita, Marinara, Parma, Vegana, Suco, Heineken, Baden Baden, Sônia e Zé)
```

## Notas

1. As subcategorias são criadas como entidades separadas de categories
2. O `SubcategoryCategory` cria o relacionamento N:N entre elas
3. Os produtos referenciam `CategoryId` para subcategorias (não para categorias pai)
4. Isso permite: Categoria Pai > Subcategoria > Produtos

