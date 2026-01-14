# Sistema de Solicitação de Mudança de Plano

## Visão Geral

Este sistema permite que clientes solicitem mudanças no plano de assinatura e que administradores gerenciem essas solicitações de forma centralizada.

## Arquitetura

### Modelo de Dados

**PlanChangeRequest** ([models/plan_change_request.go](repositories/models/plan_change_request.go))
- `id`: UUID único da solicitação
- `organization_id`: ID da organização que está solicitando a mudança
- `requested_by`: ID do usuário que fez a solicitação
- `current_package_id`: ID do pacote atual (opcional)
- `requested_package_id`: ID do pacote desejado
- `current_package_name`: Nome do pacote atual (armazenado para histórico)
- `requested_package_name`: Nome do pacote desejado
- `reason`: Motivo da mudança (opcional)
- `notes`: Observações adicionais (opcional)
- `status`: Status da solicitação (`pending`, `approved`, `rejected`, `cancelled`)
- `reviewed_by`: ID do admin que revisou a solicitação
- `reviewed_at`: Data/hora da revisão
- `review_notes`: Comentários do admin sobre a decisão
- `requested_billing_cycle`: Ciclo de cobrança desejado (`monthly`, `yearly`)

### Estados da Solicitação

1. **pending**: Solicitação criada, aguardando revisão do admin
2. **approved**: Solicitação aprovada pelo admin
3. **rejected**: Solicitação rejeitada pelo admin
4. **cancelled**: Solicitação cancelada pelo próprio usuário

## Endpoints da API

### Endpoints do Cliente

#### 1. Criar Solicitação
```
POST /plan-change-request
```

**Headers necessários:**
- `Authorization`: Token do usuário
- `X-Lpe-Organization-Id`: ID da organização

**Body:**
```json
{
  "requested_package_id": "uuid-do-pacote",
  "current_package_id": "uuid-do-pacote-atual",
  "requested_package_name": "Plano Premium",
  "current_package_name": "Plano Básico",
  "reason": "Preciso de mais recursos para expandir meu negócio",
  "notes": "Gostaria de começar a mudança no próximo mês",
  "requested_billing_cycle": "yearly"
}
```

**Resposta (201 Created):**
```json
{
  "message": "Plan change request created successfully",
  "data": {
    "id": "uuid-da-solicitacao",
    "organization_id": "uuid-da-org",
    "requested_by": "uuid-do-usuario",
    "status": "pending",
    "created_at": "2024-01-15T10:30:00Z",
    ...
  }
}
```

#### 2. Listar Minhas Solicitações
```
GET /plan-change-request/my-requests
```

**Resposta (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid-da-solicitacao",
      "status": "pending",
      "requested_package_name": "Plano Premium",
      "created_at": "2024-01-15T10:30:00Z",
      ...
    }
  ]
}
```

#### 3. Buscar Solicitação por ID
```
GET /plan-change-request/:id
```

**Resposta (200 OK):**
```json
{
  "data": {
    "id": "uuid-da-solicitacao",
    "status": "pending",
    ...
  }
}
```

#### 4. Cancelar Solicitação
```
POST /plan-change-request/:id/cancel
```

**Resposta (200 OK):**
```json
{
  "message": "Request cancelled successfully"
}
```

### Endpoints do Admin

**Nota:** Todos os endpoints admin requerem permissão `master_admin`

#### 1. Listar Todas as Solicitações
```
GET /admin/plan-change-request?status=pending
```

**Query Parameters:**
- `status` (opcional): Filtrar por status (`pending`, `approved`, `rejected`, `cancelled`)

**Resposta (200 OK):**
```json
{
  "data": [
    {
      "id": "uuid-da-solicitacao",
      "organization_id": "uuid-da-org",
      "requested_by": "uuid-do-usuario",
      "status": "pending",
      "requested_package_name": "Plano Premium",
      "reason": "...",
      "created_at": "2024-01-15T10:30:00Z",
      ...
    }
  ]
}
```

#### 2. Listar Solicitações Pendentes
```
GET /admin/plan-change-request/pending
```

**Resposta (200 OK):**
```json
{
  "data": [...],
  "count": 5
}
```

#### 3. Listar Solicitações por Organização
```
GET /admin/plan-change-request/organization/:orgId?status=pending
```

**Resposta (200 OK):**
```json
{
  "data": [...]
}
```

#### 4. Aprovar Solicitação
```
POST /admin/plan-change-request/:id/approve
```

**Body:**
```json
{
  "review_notes": "Aprovado. Mudança será efetivada no próximo ciclo de cobrança."
}
```

**Resposta (200 OK):**
```json
{
  "message": "Request approved successfully",
  "data": {
    "id": "uuid-da-solicitacao",
    "status": "approved",
    "reviewed_by": "uuid-do-admin",
    "reviewed_at": "2024-01-15T11:00:00Z",
    "review_notes": "Aprovado. Mudança será efetivada no próximo ciclo de cobrança.",
    ...
  }
}
```

#### 5. Rejeitar Solicitação
```
POST /admin/plan-change-request/:id/reject
```

**Body:**
```json
{
  "review_notes": "Rejeitado. Organização possui pendências financeiras."
}
```

**Resposta (200 OK):**
```json
{
  "message": "Request rejected successfully",
  "data": {
    "id": "uuid-da-solicitacao",
    "status": "rejected",
    "reviewed_by": "uuid-do-admin",
    "reviewed_at": "2024-01-15T11:00:00Z",
    "review_notes": "Rejeitado. Organização possui pendências financeiras.",
    ...
  }
}
```

## Fluxo de Trabalho

### Fluxo do Cliente

1. Cliente acessa a página de configurações/planos
2. Cliente visualiza os planos disponíveis
3. Cliente seleciona o novo plano desejado
4. Cliente preenche o formulário de solicitação (motivo, observações, etc.)
5. Sistema cria a solicitação com status `pending`
6. Cliente pode visualizar o status da solicitação em "Minhas Solicitações"
7. Cliente pode cancelar a solicitação enquanto estiver `pending`

### Fluxo do Admin

1. Admin acessa o painel administrativo
2. Admin visualiza lista de solicitações pendentes
3. Admin revisa os detalhes da solicitação:
   - Informações da organização
   - Plano atual vs. plano solicitado
   - Motivo da mudança
   - Histórico de solicitações anteriores
4. Admin toma uma decisão:
   - **Aprovar**: Adiciona notas sobre próximos passos
   - **Rejeitar**: Adiciona motivo da rejeição
5. Sistema atualiza o status e notifica o cliente

## Regras de Negócio

1. **Apenas usuários autenticados** podem criar solicitações
2. **Apenas o usuário que criou** pode visualizar e cancelar sua própria solicitação
3. **Apenas solicitações pendentes** podem ser canceladas, aprovadas ou rejeitadas
4. **Apenas Master Admins** podem aprovar/rejeitar solicitações
5. **Histórico completo** é mantido com timestamps e informações do revisor

## Permissões

### Cliente
- Criar solicitações para sua organização
- Visualizar suas próprias solicitações
- Cancelar suas próprias solicitações pendentes

### Admin (master_admin)
- Visualizar todas as solicitações
- Filtrar solicitações por status e organização
- Aprovar ou rejeitar solicitações
- Adicionar notas de revisão

## Integração com o Frontend

### Exemplo de Integração - Criar Solicitação

```typescript
const createPlanChangeRequest = async (data: {
  requestedPackageId: string;
  currentPackageId?: string;
  requestedPackageName: string;
  currentPackageName?: string;
  reason?: string;
  notes?: string;
  requestedBillingCycle?: 'monthly' | 'yearly';
}) => {
  const response = await fetch('/plan-change-request', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
      'X-Lpe-Organization-Id': organizationId,
    },
    body: JSON.stringify({
      requested_package_id: data.requestedPackageId,
      current_package_id: data.currentPackageId,
      requested_package_name: data.requestedPackageName,
      current_package_name: data.currentPackageName,
      reason: data.reason,
      notes: data.notes,
      requested_billing_cycle: data.requestedBillingCycle,
    }),
  });

  return response.json();
};
```

### Exemplo de Integração - Admin Aprovar Solicitação

```typescript
const approvePlanChangeRequest = async (requestId: string, reviewNotes: string) => {
  const response = await fetch(`/admin/plan-change-request/${requestId}/approve`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${adminToken}`,
    },
    body: JSON.stringify({
      review_notes: reviewNotes,
    }),
  });

  return response.json();
};
```

## Próximos Passos Sugeridos

1. **Notificações**: Integrar com o sistema de notificações existente para alertar:
   - Admin quando nova solicitação é criada
   - Cliente quando solicitação é aprovada/rejeitada

2. **Workflow de Aprovação**: Implementar aprovação automática para planos de menor valor

3. **Histórico**: Criar endpoint para visualizar histórico completo de mudanças de plano

4. **Dashboard**: Criar dashboard admin com estatísticas de solicitações

5. **Automação**: Após aprovação, aplicar automaticamente a mudança de plano na tabela `organization_packages`

## Estrutura de Arquivos

```
LEP-Back/
├── repositories/
│   ├── models/
│   │   └── plan_change_request.go         # Modelo de dados
│   └── plan_change_request.go             # Repository (acesso ao DB)
├── handler/
│   └── plan_change_request.go             # Business logic
├── server/
│   └── plan_change_request.go             # Endpoints HTTP
└── routes/
    └── routes.go                          # Configuração de rotas
```

## Testando o Sistema

### 1. Criar uma solicitação (Cliente)
```bash
curl -X POST http://localhost:8080/plan-change-request \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "X-Lpe-Organization-Id: YOUR_ORG_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "requested_package_id": "package-uuid",
    "requested_package_name": "Plano Premium",
    "reason": "Preciso de mais recursos"
  }'
```

### 2. Listar solicitações pendentes (Admin)
```bash
curl -X GET http://localhost:8080/admin/plan-change-request/pending \
  -H "Authorization: Bearer ADMIN_TOKEN"
```

### 3. Aprovar solicitação (Admin)
```bash
curl -X POST http://localhost:8080/admin/plan-change-request/REQUEST_ID/approve \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "review_notes": "Aprovado"
  }'
```
