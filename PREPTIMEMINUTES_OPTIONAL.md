# PrepTimeMinutes Optional Change

**Data**: 2025-11-09
**Status**: ✅ Implementado e compilado com sucesso
**Objetivo**: Tornar o campo `PrepTimeMinutes` do Product opcional

---

## Resumo da Mudança

O campo `PrepTimeMinutes` agora é **opcional** (nullable). Produtos podem ser criados/atualizados sem especificar tempo de preparo.

---

## Arquivos Modificados

### 1. `repositories/models/product.go`

**Antes**:
```go
PrepTimeMinutes int       `json:"prep_time_minutes,omitempty"`
```

**Depois**:
```go
PrepTimeMinutes *int      `json:"prep_time_minutes,omitempty"`
```

**O que mudou**: Convertido de `int` (tipo primitivo) para `*int` (pointer), permitindo null/nil.

---

### 2. `utils/order_time_calculator.go`

**Função**: `CalculateOrderPrepTime()`

**Antes**:
```go
itemTime := product.PrepTimeMinutes * item.Quantity
```

**Depois**:
```go
// Se PrepTimeMinutes for nil, usa 0
prepTime := 0
if product.PrepTimeMinutes != nil {
    prepTime = *product.PrepTimeMinutes
}
itemTime := prepTime * item.Quantity
```

**O que mudou**: Agora verifica se `PrepTimeMinutes` é nil antes de usar. Se for nil, usa 0 como default.

---

## Impacto na API

### Endpoint: POST /product (Create)
- **Antes**: `prep_time_minutes` era ignorado (tinha tag `omitempty`)
- **Depois**: Continua opcional na request, mas agora pode ser nil no banco

### Endpoint: PUT /product/{id} (Update)
- **Antes**: Poderia zerar o valor
- **Depois**: Continua opcional, permite deixar como nil

### Endpoint: GET /product/{id} (Read)
- **Antes**: Retornava 0 se não tivesse valor
- **Depois**: Retorna `null` se não tivesse valor

---

## Comportamento

### Cenário 1: Criar produto SEM PrepTimeMinutes ✅
```json
POST /product
{
  "name": "Água",
  "price_normal": 2.50,
  "type": "bebida"
  // sem prep_time_minutes
}

Response:
{
  "id": "...",
  "prep_time_minutes": null
}
```

### Cenário 2: Criar produto COM PrepTimeMinutes ✅
```json
POST /product
{
  "name": "Brigadeiro",
  "price_normal": 5.00,
  "type": "prato",
  "prep_time_minutes": 15
}

Response:
{
  "id": "...",
  "prep_time_minutes": 15
}
```

### Cenário 3: Calcular tempo de pedido
- Produtos com `PrepTimeMinutes` null usam **0 minutos** no cálculo
- Produtos com valor específico usam esse valor
- Exemplo: pedido com [produto null + produto 15min] = tempo máximo = 15 minutos

---

## Validação

O campo **NÃO está em nenhuma validação**, então:
- ✅ Não é obrigatório na criação
- ✅ Não é obrigatório na atualização
- ✅ Pode ser null no banco de dados

---

## Status de Compilação

```
✅ Compilado com sucesso (0 erros, 0 warnings)
✅ Binary: lep-system (atualizado)
```

---

## Backwards Compatibility

- ✅ Produtos antigos com `prep_time_minutes` = 0 funcionam normalmente
- ✅ Calculadora de tempo trata 0 e null de forma equivalente
- ✅ Sem quebra de compatibilidade com dados existentes

---

## Conclusão

O campo `PrepTimeMinutes` é agora totalmente opcional, permitindo criar produtos sem especificar tempo de preparo. A calculadora de tempo de pedidos funciona corretamente mesmo com produtos sem tempo definido.
