# Prompt para Implementação do Sistema de Cardápio no Frontend

## Visão Geral
Este documento fornece todas as interfaces TypeScript, endpoints da API e exemplos de implementação necessários para construir o sistema de cardápio (Admin Menu) no frontend React + TypeScript.

## Arquitetura de Dados

### Hierarquia do Sistema
```
Menu (Cardápio)
  └── Category (Categoria)
       └── Subcategory (Subcategoria)
            └── Product (Produto: Prato | Bebida | Vinho)
                 └── Tags (Características)
```

**Relacionamentos importantes:**
- Menu 1:N Category
- Subcategory N:N Category (uma subcategoria pode estar em várias categorias)
- Category/Subcategory 1:N Product
- Product N:N Tags

---

## Interfaces TypeScript

### 1. Menu (Cardápio)
```typescript
interface Menu {
  id: string; // UUID
  organization_id: string; // UUID
  project_id: string; // UUID
  name: string; // Nome do cardápio
  styling?: {
    colors?: {
      primary?: string;
      secondary?: string;
      accent?: string;
    };
    fonts?: {
      title?: string;
      body?: string;
    };
    layout?: 'grid' | 'list' | 'cards';
  }; // Configurações de estilização (JSON)
  order: number; // Ordem de exibição
  active: boolean; // Status ativo/pausado (play/pause)
  created_at: string; // ISO timestamp
  updated_at: string; // ISO timestamp
  deleted_at?: string; // ISO timestamp (soft delete)
}

// DTO para criação
interface CreateMenuDTO {
  name: string;
  styling?: Menu['styling'];
  order?: number;
  active?: boolean;
}

// DTO para atualização
interface UpdateMenuDTO {
  id: string;
  name?: string;
  styling?: Menu['styling'];
  order?: number;
  active?: boolean;
}
```

### 2. Category (Categoria)
```typescript
interface Category {
  id: string; // UUID
  organization_id: string; // UUID
  project_id: string; // UUID
  menu_id: string; // UUID (FK para Menu)
  name: string; // Nome da categoria
  photo?: string; // URL da foto
  notes?: string; // Observações
  order: number; // Ordem de exibição
  active: boolean; // Status ativo/pausado
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

interface CreateCategoryDTO {
  menu_id: string;
  name: string;
  photo?: string;
  notes?: string;
  order?: number;
  active?: boolean;
}

interface UpdateCategoryDTO {
  id: string;
  menu_id?: string;
  name?: string;
  photo?: string;
  notes?: string;
  order?: number;
  active?: boolean;
}
```

### 3. Subcategory (Subcategoria)
```typescript
interface Subcategory {
  id: string; // UUID
  organization_id: string; // UUID
  project_id: string; // UUID
  name: string; // Nome da subcategoria
  photo?: string; // URL da foto
  notes?: string; // Observações
  order: number; // Ordem de exibição
  active: boolean; // Status ativo/pausado
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

interface CreateSubcategoryDTO {
  name: string;
  photo?: string;
  notes?: string;
  order?: number;
  active?: boolean;
}

interface UpdateSubcategoryDTO {
  id: string;
  name?: string;
  photo?: string;
  notes?: string;
  order?: number;
  active?: boolean;
}
```

### 4. Product (Produto - Refatorado)
```typescript
type ProductType = 'prato' | 'bebida' | 'vinho';

interface Product {
  // Campos base
  id: string; // UUID
  organization_id: string; // UUID
  project_id: string; // UUID
  name: string;
  description?: string;
  image_url?: string;

  // Tipo e organização
  type: ProductType; // "prato" | "bebida" | "vinho"
  order: number; // Ordem de exibição
  active: boolean; // Status ativo/pausado (play/pause)
  pdv_code?: string; // Código PDV

  // Relacionamentos com estrutura de cardápio
  category_id?: string; // UUID (FK para Category)
  subcategory_id?: string; // UUID (FK para Subcategory)

  // Campos de preço
  price_normal: number; // Preço normal (obrigatório)
  price_promo?: number; // Preço promocional

  // Campos para Bebida/Vinho
  volume?: number; // Volume em ml
  alcohol_content?: number; // Teor alcoólico em %

  // Campos específicos para Vinho
  vintage?: string; // Safra (ano)
  country?: string; // País
  region?: string; // Região
  winery?: string; // Vinícola
  wine_type?: string; // Tipo de vinho (tinto, branco, rosé, etc)
  grapes?: string[]; // Array de uvas
  price_bottle?: number; // Preço garrafa
  price_half_bottle?: number; // Preço meia garrafa
  price_glass?: number; // Preço taça

  // Outros campos
  stock?: number;
  prep_time_minutes?: number;

  // Campos deprecados (manter para compatibilidade)
  category?: string; // DEPRECATED - usar category_id
  available?: boolean; // DEPRECATED - usar active
  price?: number; // DEPRECATED - usar price_normal

  // Timestamps
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

interface CreateProductDTO {
  name: string;
  description?: string;
  image_url?: string;
  type: ProductType;
  order?: number;
  active?: boolean;
  pdv_code?: string;
  category_id?: string;
  subcategory_id?: string;
  price_normal: number;
  price_promo?: number;
  volume?: number;
  alcohol_content?: number;
  vintage?: string;
  country?: string;
  region?: string;
  winery?: string;
  wine_type?: string;
  grapes?: string[];
  price_bottle?: number;
  price_half_bottle?: number;
  price_glass?: number;
  stock?: number;
  prep_time_minutes?: number;
}

interface UpdateProductDTO {
  id: string;
  name?: string;
  description?: string;
  image_url?: string;
  type?: ProductType;
  order?: number;
  active?: boolean;
  pdv_code?: string;
  category_id?: string;
  subcategory_id?: string;
  price_normal?: number;
  price_promo?: number;
  volume?: number;
  alcohol_content?: number;
  vintage?: string;
  country?: string;
  region?: string;
  winery?: string;
  wine_type?: string;
  grapes?: string[];
  price_bottle?: number;
  price_half_bottle?: number;
  price_glass?: number;
  stock?: number;
  prep_time_minutes?: number;
}
```

### 5. Tag (Característica)
```typescript
interface Tag {
  id: string; // UUID
  organization_id: string; // UUID
  project_id: string; // UUID
  name: string; // Nome da tag (ex: "Sem Glúten", "Vegetariano")
  color?: string; // Cor em hex (ex: "#FF5733")
  description?: string; // Descrição
  entity_type?: 'product' | 'customer' | 'table' | 'reservation' | 'order'; // Tipo de entidade
  active: boolean; // Status ativo/pausado
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

interface CreateTagDTO {
  name: string;
  color?: string; // Hex color (ex: "#FF5733")
  description?: string;
  entity_type?: string;
  active?: boolean;
}

interface UpdateTagDTO {
  id: string;
  name?: string;
  color?: string;
  description?: string;
  entity_type?: string;
  active?: boolean;
}
```

---

## Endpoints da API

### Base URL
```
http://localhost:8080
```

### Headers Obrigatórios
```typescript
{
  'Authorization': `Bearer ${token}`,
  'X-Lpe-Organization-Id': organizationId,
  'X-Lpe-Project-Id': projectId,
  'Content-Type': 'application/json'
}
```

---

## 1. Menu (Cardápio) - `/menu`

### Listar todos os menus
```typescript
GET /menu
Response: Menu[]
```

### Listar apenas menus ativos
```typescript
GET /menu/active
Response: Menu[]
```

### Obter menu por ID
```typescript
GET /menu/:id
Response: Menu
```

### Criar menu
```typescript
POST /menu
Body: CreateMenuDTO
Response: Menu
```

### Atualizar menu
```typescript
PUT /menu/:id
Body: UpdateMenuDTO
Response: Menu
```

### Atualizar ordem do menu
```typescript
PUT /menu/:id/order
Body: { order: number }
Response: { message: string }
```

### Atualizar status do menu (play/pause)
```typescript
PUT /menu/:id/status
Body: { active: boolean }
Response: { message: string }
```

### Deletar menu (soft delete)
```typescript
DELETE /menu/:id
Response: { message: string }
```

---

## 2. Category (Categoria) - `/category`

### Listar todas as categorias
```typescript
GET /category
Response: Category[]
```

### Listar apenas categorias ativas
```typescript
GET /category/active
Response: Category[]
```

### Obter categoria por ID
```typescript
GET /category/:id
Response: Category
```

### Obter categorias por menu
```typescript
GET /category/menu/:menuId
Response: Category[]
```

### Criar categoria
```typescript
POST /category
Body: CreateCategoryDTO
Response: Category
```

### Atualizar categoria
```typescript
PUT /category/:id
Body: UpdateCategoryDTO
Response: Category
```

### Atualizar ordem da categoria
```typescript
PUT /category/:id/order
Body: { order: number }
Response: { message: string }
```

### Atualizar status da categoria (play/pause)
```typescript
PUT /category/:id/status
Body: { active: boolean }
Response: { message: string }
```

### Deletar categoria (soft delete)
```typescript
DELETE /category/:id
Response: { message: string }
```

---

## 3. Subcategory (Subcategoria) - `/subcategory`

### Listar todas as subcategorias
```typescript
GET /subcategory
Response: Subcategory[]
```

### Listar apenas subcategorias ativas
```typescript
GET /subcategory/active
Response: Subcategory[]
```

### Obter subcategoria por ID
```typescript
GET /subcategory/:id
Response: Subcategory
```

### Obter subcategorias por categoria
```typescript
GET /subcategory/category/:categoryId
Response: Subcategory[]
```

### Criar subcategoria
```typescript
POST /subcategory
Body: CreateSubcategoryDTO
Response: Subcategory
```

### Atualizar subcategoria
```typescript
PUT /subcategory/:id
Body: UpdateSubcategoryDTO
Response: Subcategory
```

### Atualizar ordem da subcategoria
```typescript
PUT /subcategory/:id/order
Body: { order: number }
Response: { message: string }
```

### Atualizar status da subcategoria (play/pause)
```typescript
PUT /subcategory/:id/status
Body: { active: boolean }
Response: { message: string }
```

### Deletar subcategoria (soft delete)
```typescript
DELETE /subcategory/:id
Response: { message: string }
```

### Adicionar categoria à subcategoria (N:N)
```typescript
POST /subcategory/:id/category/:categoryId
Response: { message: string }
```

### Remover categoria da subcategoria (N:N)
```typescript
DELETE /subcategory/:id/category/:categoryId
Response: { message: string }
```

### Obter todas as categorias de uma subcategoria
```typescript
GET /subcategory/:id/categories
Response: Category[]
```

---

## 4. Product (Produto) - `/product` - ATUALIZADO

### Listar todos os produtos
```typescript
GET /product
Response: Product[]
```

### Obter produto por ID
```typescript
GET /product/:id
Response: Product
```

### Criar produto
```typescript
POST /product
Body: CreateProductDTO
Response: Product
```

### Atualizar produto
```typescript
PUT /product/:id
Body: UpdateProductDTO
Response: Product
```

### Atualizar apenas imagem do produto
```typescript
PUT /product/:id/image
Body: { image_url: string }
Response: Product
```

### Deletar produto (soft delete)
```typescript
DELETE /product/:id
Response: { message: string }
```

### **NOVOS ENDPOINTS:**

### Atualizar ordem do produto
```typescript
PUT /product/:id/order
Body: { order: number }
Response: { message: string }
```

### Atualizar status do produto (play/pause)
```typescript
PUT /product/:id/status
Body: { active: boolean }
Response: { message: string }
```

### Filtrar produtos por tipo
```typescript
GET /product/type/:type
// :type = "prato" | "bebida" | "vinho"
Response: Product[]
```

### Filtrar produtos por categoria
```typescript
GET /product/category/:categoryId
Response: Product[]
```

### Filtrar produtos por subcategoria
```typescript
GET /product/subcategory/:subcategoryId
Response: Product[]
```

### Gerenciamento de Tags do Produto

### Obter todas as tags de um produto
```typescript
GET /product/:id/tags
Response: Tag[]
```

### Adicionar tag a um produto
```typescript
POST /product/:id/tags
Body: { tag_id: string }
Response: { message: string }
```

### Remover tag de um produto
```typescript
DELETE /product/:id/tags/:tagId
Response: { message: string }
```

### Buscar produtos por tag
```typescript
GET /product/by-tag?tag_id=<uuid>
Response: Product[]
```

---

## 5. Tag (Característica) - `/tag`

### Listar todas as tags
```typescript
GET /tag
Response: Tag[]
```

### Listar apenas tags ativas
```typescript
GET /tag/active
Response: Tag[]
```

### Obter tag por ID
```typescript
GET /tag/:id
Response: Tag
```

### Obter tags por tipo de entidade
```typescript
GET /tag/entity/:entityType
// :entityType = "product" | "customer" | "table" | "reservation" | "order"
Response: Tag[]
```

### Criar tag
```typescript
POST /tag
Body: CreateTagDTO
Response: Tag
```

### Atualizar tag
```typescript
PUT /tag/:id
Body: UpdateTagDTO
Response: Tag
```

### Deletar tag (soft delete)
```typescript
DELETE /tag/:id
Response: { message: string }
```

---

## Implementação dos Services

### 1. menuService.ts
```typescript
import api from './api';
import { Menu, CreateMenuDTO, UpdateMenuDTO } from '../types';

export const menuService = {
  // Listar todos
  getAll: async (): Promise<Menu[]> => {
    const { data } = await api.get<Menu[]>('/menu');
    return data;
  },

  // Listar ativos
  getActive: async (): Promise<Menu[]> => {
    const { data } = await api.get<Menu[]>('/menu/active');
    return data;
  },

  // Obter por ID
  getById: async (id: string): Promise<Menu> => {
    const { data } = await api.get<Menu>(`/menu/${id}`);
    return data;
  },

  // Criar
  create: async (menu: CreateMenuDTO): Promise<Menu> => {
    const { data } = await api.post<Menu>('/menu', menu);
    return data;
  },

  // Atualizar
  update: async (id: string, menu: UpdateMenuDTO): Promise<Menu> => {
    const { data } = await api.put<Menu>(`/menu/${id}`, menu);
    return data;
  },

  // Atualizar ordem
  updateOrder: async (id: string, order: number): Promise<void> => {
    await api.put(`/menu/${id}/order`, { order });
  },

  // Atualizar status (play/pause)
  updateStatus: async (id: string, active: boolean): Promise<void> => {
    await api.put(`/menu/${id}/status`, { active });
  },

  // Deletar
  delete: async (id: string): Promise<void> => {
    await api.delete(`/menu/${id}`);
  },
};
```

### 2. categoryService.ts
```typescript
import api from './api';
import { Category, CreateCategoryDTO, UpdateCategoryDTO } from '../types';

export const categoryService = {
  // Listar todos
  getAll: async (): Promise<Category[]> => {
    const { data } = await api.get<Category[]>('/category');
    return data;
  },

  // Listar ativos
  getActive: async (): Promise<Category[]> => {
    const { data } = await api.get<Category[]>('/category/active');
    return data;
  },

  // Obter por ID
  getById: async (id: string): Promise<Category> => {
    const { data } = await api.get<Category>(`/category/${id}`);
    return data;
  },

  // Obter por menu
  getByMenu: async (menuId: string): Promise<Category[]> => {
    const { data } = await api.get<Category[]>(`/category/menu/${menuId}`);
    return data;
  },

  // Criar
  create: async (category: CreateCategoryDTO): Promise<Category> => {
    const { data } = await api.post<Category>('/category', category);
    return data;
  },

  // Atualizar
  update: async (id: string, category: UpdateCategoryDTO): Promise<Category> => {
    const { data } = await api.put<Category>(`/category/${id}`, category);
    return data;
  },

  // Atualizar ordem
  updateOrder: async (id: string, order: number): Promise<void> => {
    await api.put(`/category/${id}/order`, { order });
  },

  // Atualizar status
  updateStatus: async (id: string, active: boolean): Promise<void> => {
    await api.put(`/category/${id}/status`, { active });
  },

  // Deletar
  delete: async (id: string): Promise<void> => {
    await api.delete(`/category/${id}`);
  },
};
```

### 3. subcategoryService.ts
```typescript
import api from './api';
import { Subcategory, Category, CreateSubcategoryDTO, UpdateSubcategoryDTO } from '../types';

export const subcategoryService = {
  // Listar todos
  getAll: async (): Promise<Subcategory[]> => {
    const { data } = await api.get<Subcategory[]>('/subcategory');
    return data;
  },

  // Listar ativos
  getActive: async (): Promise<Subcategory[]> => {
    const { data } = await api.get<Subcategory[]>('/subcategory/active');
    return data;
  },

  // Obter por ID
  getById: async (id: string): Promise<Subcategory> => {
    const { data } = await api.get<Subcategory>(`/subcategory/${id}`);
    return data;
  },

  // Obter por categoria
  getByCategory: async (categoryId: string): Promise<Subcategory[]> => {
    const { data } = await api.get<Subcategory[]>(`/subcategory/category/${categoryId}`);
    return data;
  },

  // Criar
  create: async (subcategory: CreateSubcategoryDTO): Promise<Subcategory> => {
    const { data } = await api.post<Subcategory>('/subcategory', subcategory);
    return data;
  },

  // Atualizar
  update: async (id: string, subcategory: UpdateSubcategoryDTO): Promise<Subcategory> => {
    const { data } = await api.put<Subcategory>(`/subcategory/${id}`, subcategory);
    return data;
  },

  // Atualizar ordem
  updateOrder: async (id: string, order: number): Promise<void> => {
    await api.put(`/subcategory/${id}/order`, { order });
  },

  // Atualizar status
  updateStatus: async (id: string, active: boolean): Promise<void> => {
    await api.put(`/subcategory/${id}/status`, { active });
  },

  // Deletar
  delete: async (id: string): Promise<void> => {
    await api.delete(`/subcategory/${id}`);
  },

  // Adicionar categoria à subcategoria (N:N)
  addCategory: async (subcategoryId: string, categoryId: string): Promise<void> => {
    await api.post(`/subcategory/${subcategoryId}/category/${categoryId}`);
  },

  // Remover categoria da subcategoria (N:N)
  removeCategory: async (subcategoryId: string, categoryId: string): Promise<void> => {
    await api.delete(`/subcategory/${subcategoryId}/category/${categoryId}`);
  },

  // Obter categorias da subcategoria
  getCategories: async (subcategoryId: string): Promise<Category[]> => {
    const { data } = await api.get<Category[]>(`/subcategory/${subcategoryId}/categories`);
    return data;
  },
};
```

### 4. productService.ts (ATUALIZADO)
```typescript
import api from './api';
import { Product, Tag, CreateProductDTO, UpdateProductDTO, ProductType } from '../types';

export const productService = {
  // CRUD básico (existente)
  getAll: async (): Promise<Product[]> => {
    const { data } = await api.get<Product[]>('/product');
    return data;
  },

  getById: async (id: string): Promise<Product> => {
    const { data } = await api.get<Product>(`/product/${id}`);
    return data;
  },

  create: async (product: CreateProductDTO): Promise<Product> => {
    const { data } = await api.post<Product>('/product', product);
    return data;
  },

  update: async (id: string, product: UpdateProductDTO): Promise<Product> => {
    const { data } = await api.put<Product>(`/product/${id}`, product);
    return data;
  },

  updateImage: async (id: string, imageUrl: string): Promise<Product> => {
    const { data } = await api.put<Product>(`/product/${id}/image`, { image_url: imageUrl });
    return data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/product/${id}`);
  },

  // NOVOS MÉTODOS para cardápio

  // Atualizar ordem
  updateOrder: async (id: string, order: number): Promise<void> => {
    await api.put(`/product/${id}/order`, { order });
  },

  // Atualizar status (play/pause)
  updateStatus: async (id: string, active: boolean): Promise<void> => {
    await api.put(`/product/${id}/status`, { active });
  },

  // Filtros por tipo
  getByType: async (type: ProductType): Promise<Product[]> => {
    const { data } = await api.get<Product[]>(`/product/type/${type}`);
    return data;
  },

  // Filtros por categoria
  getByCategory: async (categoryId: string): Promise<Product[]> => {
    const { data } = await api.get<Product[]>(`/product/category/${categoryId}`);
    return data;
  },

  // Filtros por subcategoria
  getBySubcategory: async (subcategoryId: string): Promise<Product[]> => {
    const { data } = await api.get<Product[]>(`/product/subcategory/${subcategoryId}`);
    return data;
  },

  // Gerenciamento de tags (existente)
  getTags: async (productId: string): Promise<Tag[]> => {
    const { data } = await api.get<Tag[]>(`/product/${productId}/tags`);
    return data;
  },

  addTag: async (productId: string, tagId: string): Promise<void> => {
    await api.post(`/product/${productId}/tags`, { tag_id: tagId });
  },

  removeTag: async (productId: string, tagId: string): Promise<void> => {
    await api.delete(`/product/${productId}/tags/${tagId}`);
  },

  getByTag: async (tagId: string): Promise<Product[]> => {
    const { data } = await api.get<Product[]>(`/product/by-tag?tag_id=${tagId}`);
    return data;
  },
};
```

### 5. tagService.ts
```typescript
import api from './api';
import { Tag, CreateTagDTO, UpdateTagDTO } from '../types';

export const tagService = {
  // Listar todos
  getAll: async (): Promise<Tag[]> => {
    const { data } = await api.get<Tag[]>('/tag');
    return data;
  },

  // Listar ativos
  getActive: async (): Promise<Tag[]> => {
    const { data } = await api.get<Tag[]>('/tag/active');
    return data;
  },

  // Obter por ID
  getById: async (id: string): Promise<Tag> => {
    const { data } = await api.get<Tag>(`/tag/${id}`);
    return data;
  },

  // Obter por tipo de entidade
  getByEntityType: async (entityType: string): Promise<Tag[]> => {
    const { data } = await api.get<Tag[]>(`/tag/entity/${entityType}`);
    return data;
  },

  // Criar
  create: async (tag: CreateTagDTO): Promise<Tag> => {
    const { data } = await api.post<Tag>('/tag', tag);
    return data;
  },

  // Atualizar
  update: async (id: string, tag: UpdateTagDTO): Promise<Tag> => {
    const { data } = await api.put<Tag>(`/tag/${id}`, tag);
    return data;
  },

  // Deletar
  delete: async (id: string): Promise<void> => {
    await api.delete(`/tag/${id}`);
  },
};
```

---

## Exemplos de Uso

### 1. Criar estrutura completa de cardápio
```typescript
// 1. Criar Menu
const menu = await menuService.create({
  name: "Cardápio Executivo",
  styling: {
    colors: {
      primary: "#FF5733",
      secondary: "#C70039",
      accent: "#900C3F"
    },
    fonts: {
      title: "Playfair Display",
      body: "Roboto"
    },
    layout: "grid"
  },
  order: 1,
  active: true
});

// 2. Criar Categoria
const category = await categoryService.create({
  menu_id: menu.id,
  name: "Pratos Principais",
  photo: "https://example.com/category.jpg",
  notes: "Pratos quentes servidos diariamente",
  order: 1,
  active: true
});

// 3. Criar Subcategoria
const subcategory = await subcategoryService.create({
  name: "Massas",
  photo: "https://example.com/subcategory.jpg",
  notes: "Massas artesanais",
  order: 1,
  active: true
});

// 4. Associar Subcategoria à Categoria (N:N)
await subcategoryService.addCategory(subcategory.id, category.id);

// 5. Criar Produto (Prato)
const product = await productService.create({
  name: "Lasanha Bolonhesa",
  description: "Lasanha tradicional com molho bolonhesa e queijo gratinado",
  type: "prato",
  category_id: category.id,
  subcategory_id: subcategory.id,
  price_normal: 45.00,
  price_promo: 39.90,
  prep_time_minutes: 30,
  order: 1,
  active: true
});

// 6. Criar Tags
const tagVegetarian = await tagService.create({
  name: "Vegetariano",
  color: "#4CAF50",
  entity_type: "product",
  active: true
});

const tagGlutenFree = await tagService.create({
  name: "Sem Glúten",
  color: "#FFC107",
  entity_type: "product",
  active: true
});

// 7. Adicionar tags ao produto
await productService.addTag(product.id, tagVegetarian.id);
await productService.addTag(product.id, tagGlutenFree.id);
```

### 2. Criar produto tipo Vinho
```typescript
const wine = await productService.create({
  name: "Cabernet Sauvignon Reserva",
  description: "Vinho tinto encorpado com notas de frutas vermelhas",
  type: "vinho",
  category_id: categoryId,
  subcategory_id: subcategoryId,

  // Campos de vinho
  vintage: "2019",
  country: "Chile",
  region: "Vale do Maipo",
  winery: "Concha y Toro",
  wine_type: "Tinto Seco",
  grapes: ["Cabernet Sauvignon", "Merlot"],
  volume: 750,
  alcohol_content: 13.5,

  // Preços múltiplos
  price_normal: 120.00, // Preço padrão
  price_bottle: 120.00,
  price_half_bottle: 65.00,
  price_glass: 25.00,

  order: 1,
  active: true
});
```

### 3. Filtrar produtos por estrutura de cardápio
```typescript
// Buscar todos os pratos
const dishes = await productService.getByType('prato');

// Buscar produtos de uma categoria
const categoryProducts = await productService.getByCategory(categoryId);

// Buscar produtos de uma subcategoria
const subcategoryProducts = await productService.getBySubcategory(subcategoryId);

// Buscar produtos com tag específica
const vegetarianProducts = await productService.getByTag(tagVegetarianId);
```

### 4. Gerenciar ordenação (Drag & Drop)
```typescript
// Atualizar ordem ao arrastar item
const handleDragEnd = async (result: DropResult) => {
  if (!result.destination) return;

  const items = Array.from(products);
  const [reorderedItem] = items.splice(result.source.index, 1);
  items.splice(result.destination.index, 0, reorderedItem);

  // Atualizar ordem no backend
  for (let i = 0; i < items.length; i++) {
    await productService.updateOrder(items[i].id, i);
  }
};
```

### 5. Toggle play/pause
```typescript
// Pausar/ativar produto
const handleToggleStatus = async (productId: string, currentStatus: boolean) => {
  await productService.updateStatus(productId, !currentStatus);
  // Atualizar lista
  fetchProducts();
};
```

---

## Validações Frontend

### Validação de Menu
```typescript
const validateMenu = (menu: CreateMenuDTO): string[] => {
  const errors: string[] = [];

  if (!menu.name || menu.name.trim().length === 0) {
    errors.push('Nome do menu é obrigatório');
  }

  if (menu.name && menu.name.length > 100) {
    errors.push('Nome do menu deve ter no máximo 100 caracteres');
  }

  if (menu.order !== undefined && menu.order < 0) {
    errors.push('Ordem deve ser maior ou igual a 0');
  }

  return errors;
};
```

### Validação de Produto
```typescript
const validateProduct = (product: CreateProductDTO): string[] => {
  const errors: string[] = [];

  if (!product.name || product.name.trim().length === 0) {
    errors.push('Nome do produto é obrigatório');
  }

  if (!product.type || !['prato', 'bebida', 'vinho'].includes(product.type)) {
    errors.push('Tipo de produto inválido');
  }

  if (!product.price_normal || product.price_normal <= 0) {
    errors.push('Preço normal é obrigatório e deve ser maior que 0');
  }

  // Validações específicas para vinho
  if (product.type === 'vinho') {
    if (!product.vintage) {
      errors.push('Safra é obrigatória para vinhos');
    }
    if (!product.country) {
      errors.push('País é obrigatório para vinhos');
    }
    if (!product.grapes || product.grapes.length === 0) {
      errors.push('Uvas são obrigatórias para vinhos');
    }
  }

  return errors;
};
```

### Validação de Tag
```typescript
const validateTag = (tag: CreateTagDTO): string[] => {
  const errors: string[] = [];

  if (!tag.name || tag.name.trim().length === 0) {
    errors.push('Nome da tag é obrigatório');
  }

  if (tag.name && tag.name.length > 50) {
    errors.push('Nome da tag deve ter no máximo 50 caracteres');
  }

  if (tag.color && !/^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/.test(tag.color)) {
    errors.push('Cor deve ser um código hexadecimal válido (ex: #FF5733)');
  }

  return errors;
};
```

---

## Sugestões de Componentes

### 1. MenuManager (Gerenciador de Cardápios)
```typescript
interface MenuManagerProps {
  onSelectMenu?: (menu: Menu) => void;
}

const MenuManager: React.FC<MenuManagerProps> = ({ onSelectMenu }) => {
  const [menus, setMenus] = useState<Menu[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchMenus();
  }, []);

  const fetchMenus = async () => {
    setLoading(true);
    try {
      const data = await menuService.getAll();
      setMenus(data);
    } catch (error) {
      console.error('Error fetching menus:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleToggleStatus = async (id: string, active: boolean) => {
    await menuService.updateStatus(id, !active);
    fetchMenus();
  };

  return (
    <div>
      {/* Implementação da UI */}
    </div>
  );
};
```

### 2. CategoryList (Lista de Categorias com Drag & Drop)
```typescript
import { DragDropContext, Droppable, Draggable, DropResult } from '@hello-pangea/dnd';

interface CategoryListProps {
  menuId: string;
}

const CategoryList: React.FC<CategoryListProps> = ({ menuId }) => {
  const [categories, setCategories] = useState<Category[]>([]);

  const handleDragEnd = async (result: DropResult) => {
    if (!result.destination) return;

    const items = Array.from(categories);
    const [reorderedItem] = items.splice(result.source.index, 1);
    items.splice(result.destination.index, 0, reorderedItem);

    setCategories(items);

    // Atualizar ordem no backend
    for (let i = 0; i < items.length; i++) {
      await categoryService.updateOrder(items[i].id, i);
    }
  };

  return (
    <DragDropContext onDragEnd={handleDragEnd}>
      <Droppable droppableId="categories">
        {(provided) => (
          <div {...provided.droppableProps} ref={provided.innerRef}>
            {categories.map((category, index) => (
              <Draggable key={category.id} draggableId={category.id} index={index}>
                {(provided) => (
                  <div
                    ref={provided.innerRef}
                    {...provided.draggableProps}
                    {...provided.dragHandleProps}
                  >
                    {/* Categoria UI */}
                  </div>
                )}
              </Draggable>
            ))}
            {provided.placeholder}
          </div>
        )}
      </Droppable>
    </DragDropContext>
  );
};
```

### 3. ProductForm (Formulário de Produto com tipo dinâmico)
```typescript
const ProductForm: React.FC = () => {
  const [productType, setProductType] = useState<ProductType>('prato');
  const [formData, setFormData] = useState<CreateProductDTO>({
    name: '',
    type: 'prato',
    price_normal: 0,
  });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const errors = validateProduct(formData);
    if (errors.length > 0) {
      // Mostrar erros
      return;
    }
    await productService.create(formData);
  };

  return (
    <form onSubmit={handleSubmit}>
      {/* Campos base */}
      <input name="name" value={formData.name} onChange={...} />
      <select value={productType} onChange={(e) => setProductType(e.target.value as ProductType)}>
        <option value="prato">Prato</option>
        <option value="bebida">Bebida</option>
        <option value="vinho">Vinho</option>
      </select>

      {/* Campos condicionais baseado no tipo */}
      {productType === 'bebida' && (
        <>
          <input name="volume" type="number" placeholder="Volume (ml)" />
          <input name="alcohol_content" type="number" placeholder="Teor alcoólico (%)" />
        </>
      )}

      {productType === 'vinho' && (
        <>
          <input name="vintage" placeholder="Safra" />
          <input name="country" placeholder="País" />
          <input name="region" placeholder="Região" />
          <input name="winery" placeholder="Vinícola" />
          <select name="wine_type">
            <option value="Tinto Seco">Tinto Seco</option>
            <option value="Branco Seco">Branco Seco</option>
            <option value="Rosé">Rosé</option>
            <option value="Espumante">Espumante</option>
          </select>
          {/* Multi-select para uvas */}
          <input name="price_bottle" type="number" placeholder="Preço garrafa" />
          <input name="price_half_bottle" type="number" placeholder="Preço meia garrafa" />
          <input name="price_glass" type="number" placeholder="Preço taça" />
        </>
      )}

      <button type="submit">Salvar Produto</button>
    </form>
  );
};
```

### 4. TagSelector (Seletor de Tags)
```typescript
interface TagSelectorProps {
  productId: string;
  selectedTags: Tag[];
  onTagsChange: (tags: Tag[]) => void;
}

const TagSelector: React.FC<TagSelectorProps> = ({ productId, selectedTags, onTagsChange }) => {
  const [availableTags, setAvailableTags] = useState<Tag[]>([]);

  useEffect(() => {
    fetchAvailableTags();
  }, []);

  const fetchAvailableTags = async () => {
    const tags = await tagService.getByEntityType('product');
    setAvailableTags(tags.filter(tag => tag.active));
  };

  const handleAddTag = async (tagId: string) => {
    await productService.addTag(productId, tagId);
    const updatedTags = await productService.getTags(productId);
    onTagsChange(updatedTags);
  };

  const handleRemoveTag = async (tagId: string) => {
    await productService.removeTag(productId, tagId);
    const updatedTags = await productService.getTags(productId);
    onTagsChange(updatedTags);
  };

  return (
    <div>
      {/* UI de seleção de tags */}
    </div>
  );
};
```

---

## Padrões de Estado

### Context para Cardápio
```typescript
interface MenuContextType {
  currentMenu: Menu | null;
  currentCategory: Category | null;
  currentSubcategory: Subcategory | null;
  setCurrentMenu: (menu: Menu | null) => void;
  setCurrentCategory: (category: Category | null) => void;
  setCurrentSubcategory: (subcategory: Subcategory | null) => void;
}

const MenuContext = createContext<MenuContextType | undefined>(undefined);

export const MenuProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [currentMenu, setCurrentMenu] = useState<Menu | null>(null);
  const [currentCategory, setCurrentCategory] = useState<Category | null>(null);
  const [currentSubcategory, setCurrentSubcategory] = useState<Subcategory | null>(null);

  return (
    <MenuContext.Provider value={{
      currentMenu,
      currentCategory,
      currentSubcategory,
      setCurrentMenu,
      setCurrentCategory,
      setCurrentSubcategory,
    }}>
      {children}
    </MenuContext.Provider>
  );
};

export const useMenu = () => {
  const context = useContext(MenuContext);
  if (!context) throw new Error('useMenu must be used within MenuProvider');
  return context;
};
```

---

## Casos de Uso Importantes

### 1. Navegação Hierárquica
```typescript
// Breadcrumb de navegação: Menu > Categoria > Subcategoria > Produto
const MenuBreadcrumb: React.FC = () => {
  const { currentMenu, currentCategory, currentSubcategory } = useMenu();

  return (
    <nav>
      {currentMenu && <span>{currentMenu.name}</span>}
      {currentCategory && <span> > {currentCategory.name}</span>}
      {currentSubcategory && <span> > {currentSubcategory.name}</span>}
    </nav>
  );
};
```

### 2. Busca e Filtros Combinados
```typescript
const ProductSearch: React.FC = () => {
  const [filters, setFilters] = useState({
    type: '' as ProductType | '',
    categoryId: '',
    subcategoryId: '',
    tagId: '',
    searchTerm: '',
  });

  const fetchProducts = async () => {
    let products: Product[] = [];

    if (filters.type) {
      products = await productService.getByType(filters.type);
    } else if (filters.categoryId) {
      products = await productService.getByCategory(filters.categoryId);
    } else if (filters.subcategoryId) {
      products = await productService.getBySubcategory(filters.subcategoryId);
    } else if (filters.tagId) {
      products = await productService.getByTag(filters.tagId);
    } else {
      products = await productService.getAll();
    }

    // Filtro local por termo de busca
    if (filters.searchTerm) {
      products = products.filter(p =>
        p.name.toLowerCase().includes(filters.searchTerm.toLowerCase())
      );
    }

    return products;
  };

  // ...
};
```

### 3. Multi-select de Categorias para Subcategoria
```typescript
const SubcategoryCategorySelector: React.FC<{ subcategoryId: string }> = ({ subcategoryId }) => {
  const [availableCategories, setAvailableCategories] = useState<Category[]>([]);
  const [selectedCategories, setSelectedCategories] = useState<Category[]>([]);

  useEffect(() => {
    fetchCategories();
    fetchSelectedCategories();
  }, [subcategoryId]);

  const fetchCategories = async () => {
    const categories = await categoryService.getActive();
    setAvailableCategories(categories);
  };

  const fetchSelectedCategories = async () => {
    const categories = await subcategoryService.getCategories(subcategoryId);
    setSelectedCategories(categories);
  };

  const handleToggleCategory = async (categoryId: string) => {
    const isSelected = selectedCategories.some(c => c.id === categoryId);

    if (isSelected) {
      await subcategoryService.removeCategory(subcategoryId, categoryId);
    } else {
      await subcategoryService.addCategory(subcategoryId, categoryId);
    }

    fetchSelectedCategories();
  };

  return (
    <div>
      {availableCategories.map(category => (
        <label key={category.id}>
          <input
            type="checkbox"
            checked={selectedCategories.some(c => c.id === category.id)}
            onChange={() => handleToggleCategory(category.id)}
          />
          {category.name}
        </label>
      ))}
    </div>
  );
};
```

---

## Observações Finais

### Soft Delete
Todos os endpoints de DELETE fazem soft delete (define `deleted_at`). Itens deletados não aparecem nas listagens normais.

### Ordenação
A ordenação é manual via campo `order`. Implementar drag & drop para melhor UX.

### Play/Pause
O campo `active` controla se o item está visível/disponível. Usar toggle switch na UI.

### Multi-tenant
Todos os requests incluem `organization_id` e `project_id` nos headers automaticamente via interceptors do axios.

### Relacionamento N:N Subcategory-Category
Uma subcategoria pode estar em múltiplas categorias. Use multi-select ou checkboxes para gerenciar.

### Tipos de Produto
- **Prato**: Campos básicos + prep_time_minutes
- **Bebida**: Campos básicos + volume + alcohol_content
- **Vinho**: Todos os campos específicos de vinho (vintage, country, region, winery, wine_type, grapes, preços múltiplos)

### Upload de Imagens
Use o endpoint existente `/upload/product/image` para fazer upload de imagens e obter a URL, depois atualize o produto com `PUT /product/:id/image`.

---

## Checklist de Implementação

- [ ] Criar interfaces TypeScript (Menu, Category, Subcategory, Product atualizado, Tag)
- [ ] Implementar services (menuService, categoryService, subcategoryService, atualizar productService, tagService)
- [ ] Criar componente MenuManager (CRUD de menus)
- [ ] Criar componente CategoryList com drag & drop
- [ ] Criar componente SubcategoryList com drag & drop
- [ ] Criar componente ProductForm com campos dinâmicos baseado em tipo
- [ ] Criar componente TagManager (CRUD de tags)
- [ ] Criar componente TagSelector (adicionar/remover tags de produtos)
- [ ] Implementar SubcategoryCategorySelector (N:N)
- [ ] Implementar filtros de produtos (tipo, categoria, subcategoria, tag)
- [ ] Implementar breadcrumb de navegação hierárquica
- [ ] Adicionar validações frontend
- [ ] Implementar toggle play/pause em todos os níveis
- [ ] Testar ordenação com drag & drop
- [ ] Testar fluxo completo de criação: Menu > Category > Subcategory > Product > Tags

---

Pronto! Use este prompt como guia completo para implementar o sistema de cardápio no frontend. 🚀
