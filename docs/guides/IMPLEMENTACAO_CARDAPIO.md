# Implementação do Sistema de Admin Cardápio

## ✅ Já Implementado

### Modelos (PostgresLEP.go)
- [x] Menu (Cardápio)
- [x] Category (Categoria)
- [x] Subcategory (Subcategoria)
- [x] SubcategoryCategory (relacionamento N:N)
- [x] Product (refatorado com novos campos)

### Repositories
- [x] MenuRepository (`repositories/menu.go`)
- [x] CategoryRepository (`repositories/category.go`)
- [x] SubcategoryRepository (`repositories/subcategory.go`)
- [x] ProductRepository (atualizado)

### Handlers
- [x] MenuHandler (`handler/menu.go`)
- [x] CategoryHandler (`handler/category.go`)
- [x] SubcategoryHandler (`handler/subcategory.go`)
- [x] ProductHandler (atualizado)

## 🔨 Falta Implementar

### 1. Servers (Controllers)
Criar arquivos:
- [ ] `server/menu.go`
- [ ] `server/category.go`
- [ ] `server/subcategory.go`
- [ ] Atualizar `server/product.go` com novos endpoints

### 2. Validações
Criar/atualizar arquivos em `resource/validation/`:
- [ ] `menu.go`
- [ ] `category.go`
- [ ] `subcategory.go`
- [ ] Atualizar `product.go`

### 3. Injeção de Dependências
Atualizar arquivos:
- [ ] `repositories/inject.go` - adicionar Menus, Categories, Subcategories
- [ ] `handler/inject.go` - adicionar HandlerMenu, HandlerCategory, HandlerSubcategory
- [ ] `server/inject.go` - adicionar SourceMenu, SourceCategory, SourceSubcategory

### 4. Rotas da API
Atualizar `routes/routes.go`:
- [ ] setupMenuRoutes()
- [ ] setupCategoryRoutes()
- [ ] setupSubcategoryRoutes()
- [ ] Atualizar setupProductRoutes()

## 📋 Endpoints a Criar

### Menu
```
GET    /menu                  - Listar
GET    /menu/:id              - Buscar
GET    /menu/active           - Listar ativos
POST   /menu                  - Criar
PUT    /menu/:id              - Atualizar
PUT    /menu/:id/order        - Atualizar ordem
PUT    /menu/:id/status       - Play/Pause
DELETE /menu/:id              - Soft delete
```

### Category
```
GET    /category                - Listar
GET    /category/:id            - Buscar
GET    /category/menu/:menuId   - Por cardápio
GET    /category/active         - Listar ativos
POST   /category                - Criar
PUT    /category/:id            - Atualizar
PUT    /category/:id/order      - Atualizar ordem
PUT    /category/:id/status     - Play/Pause
DELETE /category/:id            - Soft delete
```

### Subcategory
```
GET    /subcategory                    - Listar
GET    /subcategory/:id                - Buscar
GET    /subcategory/category/:catId    - Por categoria
GET    /subcategory/active             - Listar ativos
POST   /subcategory                    - Criar
PUT    /subcategory/:id                - Atualizar
PUT    /subcategory/:id/categories     - Vincular categorias
PUT    /subcategory/:id/order          - Atualizar ordem
PUT    /subcategory/:id/status         - Play/Pause
DELETE /subcategory/:id                - Soft delete
```

### Product (novos endpoints)
```
GET    /product?type=prato         - Filtrar por tipo
GET    /product?category_id=xxx    - Filtrar por categoria
GET    /product?subcategory_id=xxx - Filtrar por subcategoria
PUT    /product/:id/order          - Atualizar ordem
PUT    /product/:id/status         - Play/Pause
```

## 🔄 Campos Novos do Product

### Campos Gerais
- `type` (string): "prato" | "bebida" | "vinho"
- `order` (int): ordem de exibição
- `active` (bool): status ativo/inativo
- `pdv_code` (string): código PDV
- `category_id` (UUID): FK para Category
- `subcategory_id` (UUID): FK para Subcategory

### Preços
- `price_normal` (float64): preço normal
- `price_promo` (float64): preço promocional

### Bebida/Vinho
- `volume` (int): volume em ml
- `alcohol_content` (float64): teor alcoólico %

### Vinho Específico
- `vintage` (string): safra
- `country` (string): país de origem
- `region` (string): região
- `winery` (string): vinícola
- `wine_type` (string): tipo do vinho
- `grapes` ([]string): uvas (array)
- `price_bottle` (float64): preço garrafa
- `price_half_bottle` (float64): preço meia garrafa
- `price_glass` (float64): preço taça

## 🗺️ Hierarquia

```
Cardápio (Menu)
  └── Categoria (Category)
        └── Subcategoria (Subcategory) [N:N]
              └── Produto (Product)
                    ├── Prato (type="prato")
                    ├── Bebida (type="bebida")
                    └── Vinho (type="vinho")
```

## ⚡ Próximos Passos

1. Criar servers (menu.go, category.go, subcategory.go)
2. Criar validações
3. Configurar injeção de dependências
4. Criar rotas
5. Testar endpoints
6. Atualizar frontend
