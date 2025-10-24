# 🍕 Seed Fattoria Pizzeria

Documentação completa do seed de dados para a Fattoria Pizzeria.

## 📋 Visão Geral

O seed `seed_fattoria.go` foi criado para popular o banco de dados com dados específicos da Fattoria Pizzeria, incluindo:

- **5 Pizzas** em 2 subcategorias (Entradas e Pizzas)
- **4 Bebidas** em 4 subcategorias (Soft drinks, Cervejas, Cervejas artesanais, Coquetéis)
- **Tags** (Vegetariana, Vegana)
- **3 Mesas** com estatutos diferenciados
- **Usuário Admin** específico para Fattoria

## 🚀 Como Usar

### Opção 1: Script bash (Recomendado)

```bash
# Seed básico
bash scripts/run_seed_fattoria.sh

# Seed com limpeza de dados anteriores
bash scripts/run_seed_fattoria.sh --clear-first

# Seed com output verboso
bash scripts/run_seed_fattoria.sh --verbose

# Seed completo (limpeza + verbose)
bash scripts/run_seed_fattoria.sh --clear-first --verbose

# Seed em ambiente de teste
bash scripts/run_seed_fattoria.sh --environment=test
```

### Opção 2: Comando direto Go

```bash
# Seed básico com flag --restaurant=fattoria
go run cmd/seed/main.go --restaurant=fattoria

# Seed com opções
go run cmd/seed/main.go --restaurant=fattoria --clear-first --verbose
```

### Opção 3: Seed padrão (com flag)

```bash
# Seed padrão do LEP
go run cmd/seed/main.go                        # Padrão
go run cmd/seed/main.go --restaurant=default   # Explícito
bash scripts/run_seed.sh                       # Via script padrão
```

## 📊 Estrutura de Dados

### Organização
- **ID**: `223e4567-e89b-12d3-a456-426614174100`
- **Nome**: Fattoria Pizzeria
- **Email**: contato@fattoria.com.br
- **Website**: https://fattoria.com.br

### Projeto
- **ID**: `223e4567-e89b-12d3-a456-426614174101`
- **Nome**: Fattoria Pizzeria - Projeto Principal

### Usuário Admin
- **Email**: admin@fattoria.com.br
- **Senha**: password
- **Permissões**: admin, products, orders, reservations, customers, tables, reports

### Menu & Categorias

```
📖 Cardápio Fattoria
├── 🍕 Pizzas
│   ├── 🍞 Entradas
│   │   └── Crostini - R$ 30,00
│   │       Massa fina levemente crocante com alecrim, parmesão e azeite
│   │
│   └── 🍕 Pizzas
│       ├── Marguerita 🌱 - R$ 80,00
│       │   Molho pomodoro, mussarela de búfala, manjericão fresco, azeite de oliva e orégano
│       │
│       ├── Marinara 🌱 - R$ 58,00
│       │   Molho pomodoro, alho em lascas, azeite de oliva e orégano
│       │
│       ├── Parma - R$ 109,00
│       │   Molho pomodoro, mussarela, tomate seco, orégano, parmesão, presunto parma e rúcula
│       │
│       └── Vegana 🌱 - R$ 60,00
│           Molho pomodoro, tomate confit, alho em lascas, azeitona preta, manjericão fresco e orégano
│
└── 🥤 Bebidas
    ├── 🥤 Soft drinks
    │   └── Suco de caju integral - R$ 15,00 (300ml)
    │
    ├── 🍺 Cervejas
    │   └── Heineken s/ álcool - R$ 13,00 (330ml)
    │
    ├── 🍻 Cervejas artesanais
    │   └── Baden Baden IPA - R$ 23,00 (600ml)
    │
    └── 🍸 Coquetéis
        └── Sônia e Zé - R$ 32,00
            Suco de limão siciliano, Ramazzotti, cachaça Dom Drinks, Monin de flor de sabugueiro
```

### Tags
- 🌱 **Vegetariana** - Verde (#4CAF50)
  - Marguerita
- 🌿 **Vegana** - Verde claro (#8BC34A)
  - Marinara
  - Vegana

### Mesas
- **Mesa 1**: 4 lugares, Status: Livre, Salão Principal - Entrada
- **Mesa 2**: 2 lugares, Status: Livre, Salão Principal - Janela
- **Mesa 3**: 6 lugares, Status: Livre, Salão Principal - Fundo

### Ambiente
- **Nome**: Salão Principal
- **Capacidade**: 60 pessoas

## 🔐 Credenciais de Login

### Admin Fattoria
```
Email: admin@fattoria.com.br
Senha: password
```

## 📝 Informações de Produtos

### Pizzas

#### 1. Crostini
- **Categoria**: Entradas
- **Preço**: R$ 30,00
- **Tempo de Preparo**: 15 min
- **Descrição**: Massa fina levemente crocante com alecrim, parmesão e azeite
- **Tags**: -

#### 2. Marguerita
- **Categoria**: Pizzas
- **Preço**: R$ 80,00
- **Tempo de Preparo**: 25 min
- **Descrição**: Molho pomodoro, mussarela de búfala, manjericão fresco, azeite de oliva e orégano
- **Tags**: 🌱 Vegetariana

#### 3. Marinara
- **Categoria**: Pizzas
- **Preço**: R$ 58,00
- **Tempo de Preparo**: 25 min
- **Descrição**: Molho pomodoro, alho em lascas, azeite de oliva e orégano
- **Tags**: 🌿 Vegana

#### 4. Parma
- **Categoria**: Pizzas
- **Preço**: R$ 109,00
- **Tempo de Preparo**: 25 min
- **Descrição**: Molho pomodoro, mussarela, tomate seco, orégano, parmesão, presunto parma e rúcula
- **Tags**: -

#### 5. Vegana
- **Categoria**: Pizzas
- **Preço**: R$ 60,00
- **Tempo de Preparo**: 25 min
- **Descrição**: Molho pomodoro, tomate confit, alho em lascas, azeitona preta, manjericão fresco e orégano
- **Tags**: 🌿 Vegana

### Bebidas

#### 1. Suco de caju integral
- **Categoria**: Soft drinks
- **Preço**: R$ 15,00
- **Volume**: 300ml
- **Tempo de Preparo**: 2 min
- **Descrição**: Suco natural de caju integral
- **Tags**: -

#### 2. Heineken s/ álcool
- **Categoria**: Cervejas
- **Preço**: R$ 13,00
- **Volume**: 330ml
- **Tempo de Preparo**: 2 min
- **Descrição**: Cerveja Heineken sem álcool
- **Tags**: -

#### 3. Baden Baden IPA
- **Categoria**: Cervejas artesanais
- **Preço**: R$ 23,00
- **Volume**: 600ml
- **Tempo de Preparo**: 2 min
- **Descrição**: Cerveja artesanal Baden Baden IPA
- **Tags**: -

#### 4. Sônia e Zé
- **Categoria**: Coquetéis
- **Preço**: R$ 32,00
- **Tempo de Preparo**: 5 min
- **Descrição**: Suco de limão siciliano, Ramazzotti, cachaça Dom Drinks, Monin de flor de sabugueiro e manjericão para decorar
- **Tags**: -

## 🛠️ Estrutura Técnica

### Arquivo Principal
- **Localização**: `utils/seed_fattoria.go`
- **Função**: `GenerateFattoriaData() *SeedData`
- **Retorna**: Struct contendo todos os dados do seed

### IDs Base para Referência

```go
// Organization & Project
FattoriaOrgID      = "223e4567-e89b-12d3-a456-426614174100"
FattoriaProjectID  = "223e4567-e89b-12d3-a456-426614174101"
FattoriaMenuID     = "223e4567-e89b-12d3-a456-426614174102"

// Categories
FattoriaCategoryPizzasID      = "223e4567-e89b-12d3-a456-426614174110"
FattoriaCategoryBebidasID     = "223e4567-e89b-12d3-a456-426614174111"

// Subcategories
FattoriaSubcategoryEntradasID = "223e4567-e89b-12d3-a456-426614174120"
FattoriaSubcategoryPizzasID   = "223e4567-e89b-12d3-a456-426614174121"
FattoriaSubcategorySoftID     = "223e4567-e89b-12d3-a456-426614174122"
FattoriaSubcategoryCervejasID = "223e4567-e89b-12d3-a456-426614174123"
FattoriaSubcatCervArteID      = "223e4567-e89b-12d3-a456-426614174124"
FattoriaSubcategoryCoqueisID  = "223e4567-e89b-12d3-a456-426614174125"

// Tags
FattoriaTagVegetarianaID = "223e4567-e89b-12d3-a456-426614174130"
FattoriaTagVeganaID      = "223e4567-e89b-12d3-a456-426614174131"

// Products (Pizzas)
FattoriaProductCrostiniID   = "223e4567-e89b-12d3-a456-426614174200"
FattoriaProductMargueritaID = "223e4567-e89b-12d3-a456-426614174201"
FattoriaProductMarinaraID   = "223e4567-e89b-12d3-a456-426614174202"
FattoriaProductParmaID      = "223e4567-e89b-12d3-a456-426614174203"
FattoriaProductVeganaID     = "223e4567-e89b-12d3-a456-426614174204"

// Products (Bebidas)
FattoriaProductSucoID       = "223e4567-e89b-12d3-a456-426614174205"
FattoriaProductBadenBadenID = "223e4567-e89b-12d3-a456-426614174206"
FattoriaProductSoniaZeID    = "223e4567-e89b-12d3-a456-426614174207"
FattoriaProductHeinekeID    = "223e4567-e89b-12d3-a456-426614174208"

// Environment
FattoriaEnvironmentID = "223e4567-e89b-12d3-a456-426614174300"
```

## 🔧 Integração com Sistema

### Flags Suportadas

O seed da Fattoria segue o mesmo padrão do seed padrão:

```bash
# Flag: --restaurant=fattoria
go run cmd/seed/main.go --restaurant=fattoria

# Flag: --clear-first
go run cmd/seed/main.go --restaurant=fattoria --clear-first

# Flag: --verbose
go run cmd/seed/main.go --restaurant=fattoria --verbose

# Flag: --environment
go run cmd/seed/main.go --restaurant=fattoria --environment=test
```

### Modificações no Código

1. **`cmd/seed/main.go`**:
   - Adicionado flag `--restaurant` (default: "default")
   - Switch statement para escolher entre `GenerateCompleteData()` e `GenerateFattoriaData()`

2. **`utils/seed_fattoria.go`** (novo):
   - Função `GenerateFattoriaData()` com todos os dados da Fattoria
   - Constantes de IDs para referência

3. **`scripts/run_seed_fattoria.sh`** (novo):
   - Script bash com menu amigável
   - Suporta todos os flags do seeder
   - Mensagens informativas em português

## 📚 Uso Prático

### Cenário 1: Setup Inicial
```bash
# Instalar dependências
go mod tidy

# Seed com limpeza completa
bash scripts/run_seed_fattoria.sh --clear-first

# Iniciar servidor
go run main.go
```

### Cenário 2: Desenvolvimento
```bash
# Seed rápido (sem limpeza)
bash scripts/run_seed_fattoria.sh

# Teste a API
curl -X GET http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@fattoria.com.br","password":"password"}'
```

### Cenário 3: Testes
```bash
# Seed com verbose para debugging
bash scripts/run_seed_fattoria.sh --verbose

# Seed em ambiente de teste
bash scripts/run_seed_fattoria.sh --environment=test
```

## 🎯 Casos de Uso

### Para Desenvolvedores
- Testar funcionalidades com dados reais
- Validar integrações com API
- Desenvolver features específicas da Fattoria

### Para QA/Testes
- Testar fluxos de pedidos
- Validar cardápio completo
- Verificar sistema de preços

### Para Apresentações
- Dados realistas em demonstrações
- Interface com menu real
- Estatísticas para análise

## 🔄 Adicionar Novos Produtos

Para adicionar novos produtos ao seed da Fattoria:

1. Abra `utils/seed_fattoria.go`
2. Crie um novo UUID constante:
   ```go
   FattoriaProductNovoID = uuid.MustParse("223e4567-e89b-12d3-a456-426614174XXX")
   ```
3. Adicione o produto ao slice `Products`:
   ```go
   {
       Id:              FattoriaProductNovoID,
       OrganizationId:  FattoriaOrgID,
       ProjectId:       FattoriaProjectID,
       Name:            "Nome do Produto",
       Description:     "Descrição do produto",
       Type:            "prato" ou "bebida",
       CategoryId:      &FattoriaSubcategoryXXXID,
       PriceNormal:     99.90,
       Active:          true,
       Order:           10,
       PrepTimeMinutes: 25,
       CreatedAt:       now,
       UpdatedAt:       now,
   }
   ```
4. Se aplicável, adicione tags no slice `ProductTags`

## ❓ FAQ

**P: Como faço seed apenas da Fattoria sem limpar dados existentes?**
R: Execute `bash scripts/run_seed_fattoria.sh` (sem a flag `--clear-first`)

**P: Posso ter dados da Fattoria e do LEP Demo ao mesmo tempo?**
R: Sim, ambos têm IDs diferentes na organização. Execute o seed padrão, depois o da Fattoria.

**P: Como resetar apenas os dados da Fattoria?**
R: Use `--clear-first` para limpar tudo e fazer novo seed.

**P: Os dados da Fattoria estão sincronizados com o servidor real?**
R: Não, este é apenas um seed de desenvolvimento. Atualize `utils/seed_fattoria.go` conforme necessário.

## 📞 Suporte

Para dúvidas ou problemas:
1. Verifique se PostgreSQL está rodando
2. Confirme as credenciais em `.env`
3. Execute com `--verbose` para mais detalhes
4. Consulte a documentação do CLAUDE.md do projeto

---

**Versão**: 1.0
**Última atualização**: 2024
**Mantém compatibilidade com**: LEP v1.0+
