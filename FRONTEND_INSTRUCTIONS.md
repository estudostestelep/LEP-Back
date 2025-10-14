# Instruções de Implementação - Frontend Multi-Tenant

## 📋 Índice
1. [Visão Geral](#visão-geral)
2. [Mudanças no Login](#mudanças-no-login)
3. [Implementação Passo a Passo](#implementação-passo-a-passo)
4. [Componentes Prontos](#componentes-prontos)
5. [Testes e Validação](#testes-e-validação)

---

## Visão Geral

O backend foi refatorado para suportar **multi-tenancy completo**, onde um usuário pode ter acesso a múltiplas organizações e projetos. As principais mudanças são:

### O que mudou?

| Campo | Antes | Agora |
|-------|-------|-------|
| `User.organization_id` | ✅ Existia | ❌ Removido |
| `User.project_id` | ✅ Existia | ❌ Removido |
| `User.role` | ✅ Existia | ❌ Removido |
| `User.active` | ❌ Não existia | ✅ Adicionado |
| Relacionamentos | ❌ Não existia | ✅ `UserOrganization` e `UserProject` |

---

## Mudanças no Login

### Resposta Antiga do Login
```json
{
  "user": {
    "id": "uuid",
    "organization_id": "uuid",
    "project_id": "uuid",
    "name": "João Silva",
    "email": "joao@email.com",
    "role": "admin",
    "permissions": ["admin"]
  },
  "token": "eyJhbGc..."
}
```

### 🆕 Nova Resposta do Login
```json
{
  "user": {
    "id": "uuid",
    "name": "João Silva",
    "email": "joao@email.com",
    "permissions": ["admin"],
    "active": true,
    "created_at": "2024-01-15T10:00:00Z",
    "updated_at": "2024-01-15T10:00:00Z"
  },
  "token": "eyJhbGc...",
  "organizations": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "organization_id": "org-uuid-1",
      "role": "owner",
      "active": true,
      "created_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": "uuid",
      "user_id": "uuid",
      "organization_id": "org-uuid-2",
      "role": "admin",
      "active": true,
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "projects": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "project_id": "proj-uuid-1",
      "role": "admin",
      "active": true,
      "created_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": "uuid",
      "user_id": "uuid",
      "project_id": "proj-uuid-2",
      "role": "member",
      "active": true,
      "created_at": "2024-01-15T10:00:00Z"
    }
  ]
}
```

---

## Implementação Passo a Passo

### Passo 1: Atualizar Interfaces TypeScript

Crie ou atualize o arquivo `src/types/auth.ts`:

```typescript
// src/types/auth.ts

export interface User {
  id: string;
  name: string;
  email: string;
  permissions: string[];
  active: boolean;
  created_at: string;
  updated_at: string;
  deleted_at?: string;
}

export interface UserOrganization {
  id: string;
  user_id: string;
  organization_id: string;
  role: string; // "owner" | "admin" | "member"
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface UserProject {
  id: string;
  user_id: string;
  project_id: string;
  role: string; // "admin" | "manager" | "waiter" | "member"
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Organization {
  id: string;
  name: string;
  email?: string;
  phone?: string;
  active: boolean;
}

export interface Project {
  id: string;
  organization_id: string;
  name: string;
  description?: string;
  active: boolean;
}

export interface LoginResponse {
  user: User;
  token: string;
  organizations: UserOrganization[];
  projects: UserProject[];
}
```

### Passo 2: Atualizar AuthContext

```typescript
// src/contexts/AuthContext.tsx

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { User, UserOrganization, UserProject, Organization, Project } from '@/types/auth';
import api from '@/services/api';

interface AuthContextData {
  user: User | null;
  token: string | null;
  organizations: UserOrganization[];
  projects: UserProject[];
  currentOrganization: string | null;
  currentProject: string | null;
  organizationDetails: Organization | null;
  projectDetails: Project | null;
  loading: boolean;
  setCurrentOrganization: (orgId: string) => Promise<void>;
  setCurrentProject: (projectId: string) => Promise<void>;
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextData>({} as AuthContextData);

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [organizations, setOrganizations] = useState<UserOrganization[]>([]);
  const [projects, setProjects] = useState<UserProject[]>([]);
  const [currentOrganization, setCurrentOrganizationState] = useState<string | null>(null);
  const [currentProject, setCurrentProjectState] = useState<string | null>(null);
  const [organizationDetails, setOrganizationDetails] = useState<Organization | null>(null);
  const [projectDetails, setProjectDetails] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);

  // Carregar dados do localStorage na inicialização
  useEffect(() => {
    const storedUser = localStorage.getItem('@LEP:user');
    const storedToken = localStorage.getItem('@LEP:token');
    const storedOrgs = localStorage.getItem('@LEP:organizations');
    const storedProjs = localStorage.getItem('@LEP:projects');
    const storedCurrentOrg = localStorage.getItem('@LEP:currentOrganization');
    const storedCurrentProj = localStorage.getItem('@LEP:currentProject');

    if (storedUser && storedToken) {
      setUser(JSON.parse(storedUser));
      setToken(storedToken);
      setOrganizations(storedOrgs ? JSON.parse(storedOrgs) : []);
      setProjects(storedProjs ? JSON.parse(storedProjs) : []);
      setCurrentOrganizationState(storedCurrentOrg);
      setCurrentProjectState(storedCurrentProj);

      // Carregar detalhes da org/projeto
      if (storedCurrentOrg) {
        loadOrganizationDetails(storedCurrentOrg);
      }
      if (storedCurrentProj) {
        loadProjectDetails(storedCurrentProj);
      }
    }
    setLoading(false);
  }, []);

  const loadOrganizationDetails = async (orgId: string) => {
    try {
      const response = await api.get(`/organization/${orgId}`);
      setOrganizationDetails(response.data);
    } catch (error) {
      console.error('Erro ao carregar organização:', error);
    }
  };

  const loadProjectDetails = async (projectId: string) => {
    try {
      const response = await api.get(`/project/${projectId}`);
      setProjectDetails(response.data);
    } catch (error) {
      console.error('Erro ao carregar projeto:', error);
    }
  };

  const login = async (email: string, password: string) => {
    const response = await api.post('/login', { email, password });
    const { user, token, organizations, projects } = response.data;

    setUser(user);
    setToken(token);
    setOrganizations(organizations || []);
    setProjects(projects || []);

    // Definir primeira organização e projeto como padrão
    if (organizations?.length > 0) {
      const firstOrgId = organizations[0].organization_id;
      setCurrentOrganizationState(firstOrgId);
      await loadOrganizationDetails(firstOrgId);
    }

    if (projects?.length > 0) {
      const firstProjId = projects[0].project_id;
      setCurrentProjectState(firstProjId);
      await loadProjectDetails(firstProjId);
    }

    // Persistir no localStorage
    localStorage.setItem('@LEP:user', JSON.stringify(user));
    localStorage.setItem('@LEP:token', token);
    localStorage.setItem('@LEP:organizations', JSON.stringify(organizations || []));
    localStorage.setItem('@LEP:projects', JSON.stringify(projects || []));
    localStorage.setItem('@LEP:currentOrganization', organizations[0]?.organization_id || '');
    localStorage.setItem('@LEP:currentProject', projects[0]?.project_id || '');
  };

  const setCurrentOrganization = async (orgId: string) => {
    setCurrentOrganizationState(orgId);
    localStorage.setItem('@LEP:currentOrganization', orgId);
    await loadOrganizationDetails(orgId);

    // Filtrar projetos da organização selecionada
    const orgProjects = await getProjectsForOrganization(orgId);
    if (orgProjects.length > 0) {
      await setCurrentProject(orgProjects[0].project_id);
    }
  };

  const setCurrentProject = async (projectId: string) => {
    setCurrentProjectState(projectId);
    localStorage.setItem('@LEP:currentProject', projectId);
    await loadProjectDetails(projectId);
  };

  const getProjectsForOrganization = async (orgId: string): Promise<UserProject[]> => {
    if (!user) return [];
    try {
      const response = await api.get(`/user/${user.id}/organization/${orgId}/projects`);
      return response.data;
    } catch (error) {
      console.error('Erro ao buscar projetos da organização:', error);
      return projects.filter(p => {
        // Fallback: verificar se o projeto pertence à org via projectDetails
        return true; // Simplificado - idealmente buscar do backend
      });
    }
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    setOrganizations([]);
    setProjects([]);
    setCurrentOrganizationState(null);
    setCurrentProjectState(null);
    setOrganizationDetails(null);
    setProjectDetails(null);

    localStorage.removeItem('@LEP:user');
    localStorage.removeItem('@LEP:token');
    localStorage.removeItem('@LEP:organizations');
    localStorage.removeItem('@LEP:projects');
    localStorage.removeItem('@LEP:currentOrganization');
    localStorage.removeItem('@LEP:currentProject');
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        organizations,
        projects,
        currentOrganization,
        currentProject,
        organizationDetails,
        projectDetails,
        loading,
        setCurrentOrganization,
        setCurrentProject,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => useContext(AuthContext);
```

### Passo 3: Atualizar Axios Interceptor

```typescript
// src/services/api.ts

import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
});

// Request interceptor
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('@LEP:token');
  const currentOrg = localStorage.getItem('@LEP:currentOrganization');
  const currentProj = localStorage.getItem('@LEP:currentProject');

  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }

  if (currentOrg) {
    config.headers['X-Lpe-Organization-Id'] = currentOrg;
  }

  if (currentProj) {
    config.headers['X-Lpe-Project-Id'] = currentProj;
  }

  return config;
});

// Response interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expirado ou inválido
      localStorage.clear();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api;
```

### Passo 4: Criar Componente de Seleção Org/Projeto

```typescript
// src/components/OrgProjectSelector.tsx

import { useAuth } from '@/contexts/AuthContext';
import { Building2, FolderKanban, ChevronDown } from 'lucide-react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

export const OrgProjectSelector = () => {
  const {
    organizations,
    projects,
    currentOrganization,
    currentProject,
    organizationDetails,
    projectDetails,
    setCurrentOrganization,
    setCurrentProject,
  } = useAuth();

  // Filtrar projetos da organização atual
  const currentOrgProjects = projects.filter((p) => {
    // Você pode melhorar isso buscando os detalhes do projeto
    return true; // Simplificado
  });

  return (
    <div className="flex items-center gap-3">
      {/* Seletor de Organização */}
      <div className="flex items-center gap-2">
        <Building2 className="w-4 h-4 text-gray-500" />
        <Select
          value={currentOrganization || undefined}
          onValueChange={setCurrentOrganization}
        >
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Selecione a organização">
              {organizationDetails?.name || 'Carregando...'}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            {organizations.map((org) => (
              <SelectItem key={org.id} value={org.organization_id}>
                <div className="flex flex-col">
                  <span className="font-medium">Organização {org.organization_id.slice(0, 8)}...</span>
                  <span className="text-xs text-gray-500">Role: {org.role}</span>
                </div>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="w-px h-6 bg-gray-300" />

      {/* Seletor de Projeto */}
      <div className="flex items-center gap-2">
        <FolderKanban className="w-4 h-4 text-gray-500" />
        <Select
          value={currentProject || undefined}
          onValueChange={setCurrentProject}
        >
          <SelectTrigger className="w-[200px]">
            <SelectValue placeholder="Selecione o projeto">
              {projectDetails?.name || 'Carregando...'}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            {currentOrgProjects.map((proj) => (
              <SelectItem key={proj.id} value={proj.project_id}>
                <div className="flex flex-col">
                  <span className="font-medium">Projeto {proj.project_id.slice(0, 8)}...</span>
                  <span className="text-xs text-gray-500">Role: {proj.role}</span>
                </div>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    </div>
  );
};
```

### Passo 5: Adicionar no Header/Layout

```typescript
// src/components/Header.tsx

import { OrgProjectSelector } from './OrgProjectSelector';
import { useAuth } from '@/contexts/AuthContext';

export const Header = () => {
  const { user, logout } = useAuth();

  return (
    <header className="border-b bg-white px-6 py-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-6">
          <h1 className="text-xl font-bold">LEP System</h1>

          {/* Seletores de Org/Projeto */}
          {user && <OrgProjectSelector />}
        </div>

        <div className="flex items-center gap-4">
          <span className="text-sm text-gray-600">{user?.name}</span>
          <button
            onClick={logout}
            className="text-sm text-red-600 hover:text-red-700"
          >
            Sair
          </button>
        </div>
      </div>
    </header>
  );
};
```

### Passo 6: Criar Serviços para Novos Endpoints

```typescript
// src/services/userOrganizationService.ts

import api from './api';
import { UserOrganization } from '@/types/auth';

export const userOrganizationService = {
  addUserToOrganization: async (userId: string, data: {
    organization_id: string;
    role: string;
  }): Promise<UserOrganization> => {
    const response = await api.post(`/user/${userId}/organization`, data);
    return response.data;
  },

  removeUserFromOrganization: async (userId: string, orgId: string): Promise<void> => {
    await api.delete(`/user/${userId}/organization/${orgId}`);
  },

  updateUserOrganization: async (id: string, data: Partial<UserOrganization>): Promise<UserOrganization> => {
    const response = await api.put(`/user-organization/${id}`, data);
    return response.data;
  },

  getUserOrganizations: async (userId: string): Promise<UserOrganization[]> => {
    const response = await api.get(`/user/${userId}/organizations`);
    return response.data;
  },

  getOrganizationUsers: async (orgId: string): Promise<UserOrganization[]> => {
    const response = await api.get(`/organization/${orgId}/users`);
    return response.data;
  },
};

// src/services/userProjectService.ts

import api from './api';
import { UserProject } from '@/types/auth';

export const userProjectService = {
  addUserToProject: async (userId: string, data: {
    project_id: string;
    role: string;
  }): Promise<UserProject> => {
    const response = await api.post(`/user/${userId}/project`, data);
    return response.data;
  },

  removeUserFromProject: async (userId: string, projectId: string): Promise<void> => {
    await api.delete(`/user/${userId}/project/${projectId}`);
  },

  updateUserProject: async (id: string, data: Partial<UserProject>): Promise<UserProject> => {
    const response = await api.put(`/user-project/${id}`, data);
    return response.data;
  },

  getUserProjects: async (userId: string): Promise<UserProject[]> => {
    const response = await api.get(`/user/${userId}/projects`);
    return response.data;
  },

  getUserProjectsByOrganization: async (userId: string, orgId: string): Promise<UserProject[]> => {
    const response = await api.get(`/user/${userId}/organization/${orgId}/projects`);
    return response.data;
  },

  getProjectUsers: async (projectId: string): Promise<UserProject[]> => {
    const response = await api.get(`/project/${projectId}/users`);
    return response.data;
  },
};
```

---

## Componentes Prontos

### Componente com ShadcN UI (Recomendado)

```bash
# Instalar dependências necessárias
npx shadcn-ui@latest add select
```

```typescript
// src/components/OrgProjectSelector.tsx (versão completa com loading)

import { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { Building2, FolderKanban, Loader2 } from 'lucide-react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { toast } from 'sonner';

export const OrgProjectSelector = () => {
  const {
    organizations,
    projects,
    currentOrganization,
    currentProject,
    organizationDetails,
    projectDetails,
    setCurrentOrganization,
    setCurrentProject,
  } = useAuth();

  const [loadingOrg, setLoadingOrg] = useState(false);
  const [loadingProj, setLoadingProj] = useState(false);

  const handleOrgChange = async (orgId: string) => {
    setLoadingOrg(true);
    try {
      await setCurrentOrganization(orgId);
      toast.success('Organização alterada com sucesso');
    } catch (error) {
      toast.error('Erro ao trocar organização');
      console.error(error);
    } finally {
      setLoadingOrg(false);
    }
  };

  const handleProjectChange = async (projectId: string) => {
    setLoadingProj(true);
    try {
      await setCurrentProject(projectId);
      toast.success('Projeto alterado com sucesso');
    } catch (error) {
      toast.error('Erro ao trocar projeto');
      console.error(error);
    } finally {
      setLoadingProj(false);
    }
  };

  if (organizations.length === 0 && projects.length === 0) {
    return null;
  }

  return (
    <div className="flex items-center gap-3 bg-gray-50 rounded-lg px-4 py-2">
      {/* Seletor de Organização */}
      <div className="flex items-center gap-2">
        <Building2 className="w-4 h-4 text-gray-500" />
        <Select
          value={currentOrganization || undefined}
          onValueChange={handleOrgChange}
          disabled={loadingOrg}
        >
          <SelectTrigger className="w-[220px] bg-white">
            <SelectValue placeholder="Selecione a organização">
              {loadingOrg ? (
                <span className="flex items-center gap-2">
                  <Loader2 className="w-3 h-3 animate-spin" />
                  Carregando...
                </span>
              ) : (
                organizationDetails?.name || 'Selecione...'
              )}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            {organizations.map((org) => (
              <SelectItem key={org.id} value={org.organization_id}>
                <div className="flex items-center justify-between gap-4">
                  <span className="font-medium">Org {org.organization_id.slice(0, 8)}...</span>
                  <span className="text-xs px-2 py-0.5 bg-blue-100 text-blue-700 rounded-full">
                    {org.role}
                  </span>
                </div>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="w-px h-6 bg-gray-300" />

      {/* Seletor de Projeto */}
      <div className="flex items-center gap-2">
        <FolderKanban className="w-4 h-4 text-gray-500" />
        <Select
          value={currentProject || undefined}
          onValueChange={handleProjectChange}
          disabled={loadingProj}
        >
          <SelectTrigger className="w-[220px] bg-white">
            <SelectValue placeholder="Selecione o projeto">
              {loadingProj ? (
                <span className="flex items-center gap-2">
                  <Loader2 className="w-3 h-3 animate-spin" />
                  Carregando...
                </span>
              ) : (
                projectDetails?.name || 'Selecione...'
              )}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            {projects.map((proj) => (
              <SelectItem key={proj.id} value={proj.project_id}>
                <div className="flex items-center justify-between gap-4">
                  <span className="font-medium">Proj {proj.project_id.slice(0, 8)}...</span>
                  <span className="text-xs px-2 py-0.5 bg-green-100 text-green-700 rounded-full">
                    {proj.role}
                  </span>
                </div>
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    </div>
  );
};
```

---

## Testes e Validação

### 1. Testar Login

```bash
# Terminal - testar endpoint de login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "teste@gmail.com",
    "password": "password"
  }'
```

**Resposta esperada:**
```json
{
  "user": { ... },
  "token": "eyJ...",
  "organizations": [...],
  "projects": [...]
}
```

### 2. Testar Headers nas Requisições

Após login, verifique se os headers estão sendo enviados:

```typescript
// Abrir DevTools > Network
// Fazer qualquer requisição autenticada
// Verificar se os headers estão presentes:
// - Authorization: Bearer eyJ...
// - X-Lpe-Organization-Id: uuid
// - X-Lpe-Project-Id: uuid
```

### 3. Testar Troca de Contexto

1. Fazer login
2. Trocar organização no select
3. Verificar se `localStorage.getItem('@LEP:currentOrganization')` mudou
4. Fazer uma requisição e verificar se o header `X-Lpe-Organization-Id` foi atualizado

### 4. Validar Credenciais de Teste

Após rodar o seeder, use:

```
Email: teste@gmail.com
Senha: password
```

ou

```
Email: pablo@lep.com
Senha: senha123
```

---

## Checklist de Implementação

- [ ] Atualizar interfaces TypeScript (`User`, `UserOrganization`, `UserProject`)
- [ ] Atualizar `AuthContext` com novos estados
- [ ] Atualizar interceptor do Axios para incluir headers
- [ ] Criar componente `OrgProjectSelector`
- [ ] Adicionar seletor no Header/Layout
- [ ] Criar serviços para novos endpoints
- [ ] Testar login e verificar resposta
- [ ] Testar troca de organização
- [ ] Testar troca de projeto
- [ ] Validar headers nas requisições
- [ ] Testar logout e limpeza de dados

---

## Troubleshooting

### Problema: Headers não estão sendo enviados

**Solução:**
- Verificar se `localStorage` tem os valores corretos
- Verificar se o interceptor está configurado corretamente
- Verificar ordem de importação do `api.ts`

### Problema: Token expirado

**Solução:**
- Backend retorna 401
- Interceptor redireciona para `/login`
- Limpar `localStorage` e fazer login novamente

### Problema: Usuário não tem acesso à org/projeto

**Solução:**
- Verificar se relacionamentos `UserOrganization` e `UserProject` foram criados
- Usar endpoints para adicionar usuário às orgs/projetos:
  ```typescript
  await userOrganizationService.addUserToOrganization(userId, {
    organization_id: orgId,
    role: 'member'
  });
  ```

---

## Próximos Passos

1. ✅ Implementar AuthContext com multi-tenancy
2. ✅ Criar componente de seleção org/projeto
3. ✅ Atualizar interceptor do Axios
4. ⏭️ Criar tela de gerenciamento de usuários em orgs
5. ⏭️ Criar tela de gerenciamento de usuários em projetos
6. ⏭️ Implementar permissões baseadas em role

---

## Suporte

Para dúvidas ou problemas:
- 📄 [FRONTEND_MIGRATION_GUIDE.md](./FRONTEND_MIGRATION_GUIDE.md) - Guia completo de migração
- 📄 [CLAUDE.md](./CLAUDE.md) - Documentação do backend
- 📄 [README.md](./README.md) - Documentação geral
