# Alterações no Sistema de Reset de Theme

## Resumo das Mudanças

O endpoint `/project/settings/theme/reset` foi refatorado para **apenas apagar as cores customizadas**, deixando o frontend lidar com a cor padrão.

## Antes vs Depois

### Comportamento Anterior
- Endpoint resetava e retornava cores padrão hardcoded no backend
- Backend tinha responsabilidade de prover cor padrão
- Frontend recebia cores completas preenchidas

### Novo Comportamento
- Endpoint apaga todas as cores customizadas (seta para `nil`)
- Frontend é responsável por exibir sua cor padrão quando receber `nil`
- Separação clara de responsabilidades

## Arquivos Alterados

### 1. [`repositories/theme_customization.go`](repositories/theme_customization.go#L56-L157)
- **Função**: `ResetToDefaults(projectId uuid.UUID)`
- **Mudança**:
  - Todos os campos de cor agora são zerados para `nil`
  - Removido código de paleta padrão hardcoded
  - Ao criar novo registro: todos os campos setados como `nil`
  - Ao atualizar existente: seta todos para `nil` e chama `UpdateTheme()`

**Antes**:
```go
theme.PrimaryColorLight = &primaryLight    // "#1E3A8A"
theme.SecondaryColorLight = &secondaryLight // "#3B82F6"
// ... mais 20 linhas de cores ...
r.DeleteTheme(theme.ID)  // deletava o registro inteiro
```

**Depois**:
```go
theme.PrimaryColorLight = nil
theme.SecondaryColorLight = nil
// ... todos para nil ...
if err := r.UpdateTheme(theme); err != nil {
    return nil, err
}
```

### 2. [`handler/theme_customization.go`](handler/theme_customization.go#L146-L161)
- **Função**: `ResetToDefaults(projectId string)`
- **Mudança**:
  - Simplificado: apenas delega para repositório
  - Removida chamada para `buildDefaultTheme()`
  - Removida lógica condicional de criação de tema com defaults

**Antes**:
```go
theme, err := h.themeRepo.ResetToDefaults(projectUUID)
if err != nil {
    return nil, err
}

// Lógica complexa de verificação e criação
if theme == nil {
    org := uuid.New()
    theme = h.buildDefaultTheme(projectUUID, org)
    // ... criar tema ...
}
return theme, nil
```

**Depois**:
```go
theme, err := h.themeRepo.ResetToDefaults(projectUUID)
if err != nil {
    return nil, err
}
return theme, nil
```

## Teste da Funcionalidade

### Endpoint
```
POST /project/settings/theme/reset
Headers:
  X-Lpe-Organization-Id: {org-uuid}
  X-Lpe-Project-Id: {project-uuid}
  Authorization: Bearer {token}
```

### Response
```json
{
  "id": "uuid-do-theme",
  "project_id": "uuid-do-projeto",
  "organization_id": "uuid-da-org",
  "primary_color_light": null,
  "secondary_color_light": null,
  "background_color_light": null,
  "card_background_color_light": null,
  "text_color_light": null,
  "text_secondary_color_light": null,
  "accent_color_light": null,
  "destructive_color_light": null,
  "success_color_light": null,
  "warning_color_light": null,
  "border_color_light": null,
  "price_color_light": null,
  "focus_ring_color_light": null,
  "input_background_color_light": null,
  "primary_color_dark": null,
  "secondary_color_dark": null,
  "background_color_dark": null,
  "card_background_color_dark": null,
  "text_color_dark": null,
  "text_secondary_color_dark": null,
  "accent_color_dark": null,
  "destructive_color_dark": null,
  "success_color_dark": null,
  "warning_color_dark": null,
  "border_color_dark": null,
  "price_color_dark": null,
  "focus_ring_color_dark": null,
  "input_background_color_dark": null,
  "disabled_opacity": null,
  "shadow_intensity": null,
  "is_active": false,
  "created_at": "2025-11-09T...",
  "updated_at": "2025-11-09T..."
}
```

## Como Funciona no Frontend

Quando receber `null` em qualquer cor, o frontend deve:

1. **Usar cor padrão da aplicação**
   ```typescript
   const primaryColor = themeColors.primary_color_light || DEFAULT_PRIMARY_COLOR;
   ```

2. **Fallback para cores do Tailwind**
   ```typescript
   const bgColor = theme.background_color_light || 'bg-white';
   ```

3. **Exemplo de Implementação**
   ```typescript
   const applyTheme = (theme: ThemeCustomization) => {
     const cssVariables = {
       '--primary': theme.primary_color_light || '#000000',
       '--secondary': theme.secondary_color_light || '#666666',
       '--background': theme.background_color_light || '#ffffff',
       // ... demais variáveis ...
     };

     Object.entries(cssVariables).forEach(([key, value]) => {
       document.documentElement.style.setProperty(key, value);
     });
   };
   ```

## Benefícios

✅ **Separação de Responsabilidades**
- Backend: armazena customizações
- Frontend: responsável por cores padrão

✅ **Simplicidade**
- Código mais limpo e menor
- Menos duplicação de configuração

✅ **Flexibilidade**
- Frontend pode facilmente mudar cor padrão sem alterar backend
- Cada projeto pode ter sua própria cor padrão no frontend

✅ **Manutenibilidade**
- Paletas de cor padrão gerenciadas apenas no frontend
- Backend não precisa saber detalhes de design

## Notas para Testes

1. **GET `/project/settings/theme`** retorna tema com todos os valores `null` após reset
2. **Frontend deve ter fallbacks** para cada cor esperada
3. **Não há quebra de compatibilidade** - valores `null` são sempre permitidos
4. **Cores customizadas continuam funcionando** - apenas reset agora zera os valores

## Próximos Passos (Frontend)

- [ ] Adicionar fallbacks de cores em componentes que usam tema
- [ ] Atualizar stories do Storybook com cores padrão
- [ ] Testar visual após reset
- [ ] Documentar cores padrão na guia de Design System
