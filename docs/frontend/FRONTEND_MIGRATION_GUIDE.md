# Guia de Migração - Multi-Tenant User System

## Resumo das Mudanças

O sistema foi refatorado para suportar **múltiplas organizações e projetos por usuário**, permitindo que um único usuário tenha acesso a várias organizações e projetos sem precisar fazer logout.

## Principais Alterações no Backend

### 1. Modelo de Dados User

**ANTES:**
```typescript
interface User {
  id: string;
  organization_id: string;  // ❌ Removido
  project_id: string;       // ❌ Removido
  name: string;
  email: string;
  password: string;
  role: string;
  permissions: string[];
}
```

**AGORA:**
```typescript
interface User {
  id: string;
  name: string;
  email: string;
  password: string;
  permissions: string[];
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

### 2. Novos Modelos de Relacionamento

```typescript
interface UserOrganization {
  id: string;
  user_id: string;
  organization_id: string;
  role: string;  // "owner", "admin", "member"
  active: boolean;
  created_at: string;
  updated_at: string;
}

interface UserProject {
  id: string;
  user_id: string;
  project_id: string;
  role: string;  // "admin", "manager", "waiter", "member"
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

### 3. Resposta de Login Atualizada

**ANTES:**
```json
{
  "user": { ... },
  "token": "jwt-token"
}
```

**AGORA:**
```json
{
  "user": { ... },
  "token": "jwt-token",
  "organizations": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "organization_id": "uuid",
      "role": "owner",
      "active": true
    }
  ],
  "projects": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "project_id": "uuid",
      "role": "admin",
      "active": true
    }
  ]
}
```

### 4. Novos Endpoints

#### Gerenciar Usuário-Organização
```
POST   /user/:userId/organization              - Adicionar usuário a uma organização
DELETE /user/:userId/organization/:orgId       - Remover usuário de uma organização
PUT    /user-organization/:id                  - Atualizar relacionamento
GET    /user/:userId/organizations             - Listar organizações do usuário
GET    /organization/:orgId/users              - Listar usuários da organização
```

#### Gerenciar Usuário-Projeto
```
POST   /user/:userId/project                        - Adicionar usuário a um projeto
DELETE /user/:userId/project/:projectId             - Remover usuário de um projeto
PUT    /user-project/:id                            - Atualizar relacionamento
GET    /user/:userId/projects                       - Listar projetos do usuário
GET    /user/:userId/organization/:orgId/projects   - Listar projetos do usuário em uma org
GET    /project/:projectId/users                    - Listar usuários do projeto
```

## Mudanças Necessárias no Frontend

### 1. Atualizar Context de Autenticação

```typescript
// auth-context.tsx

interface AuthContextData {
  user: User | null;
  token: string | null;
  organizations: UserOrganization[];
  projects: UserProject[];
  currentOrganization: string | null;  // NOVO
  currentProject: string | null;       // NOVO
  setCurrentOrganization: (orgId: string) => void;  // NOVO
  setCurrentProject: (projectId: string) => void;   // NOVO
  login: (email: string, password: string) => Promise<void>;
  logout: () => void;
}

// Implementação sugerida
const AuthProvider = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [organizations, setOrganizations] = useState<UserOrganization[]>([]);
  const [projects, setProjects] = useState<UserProject[]>([]);
  const [currentOrganization, setCurrentOrganization] = useState<string | null>(null);
  const [currentProject, setCurrentProject] = useState<string | null>(null);

  const login = async (email: string, password: string) => {
    const response = await authService.login(email, password);

    setUser(response.user);
    setToken(response.token);
    setOrganizations(response.organizations || []);
    setProjects(response.projects || []);

    // Definir primeira organização e projeto como padrão
    if (response.organizations?.length > 0) {
      setCurrentOrganization(response.organizations[0].organization_id);
    }
    if (response.projects?.length > 0) {
      setCurrentProject(response.projects[0].project_id);
    }

    // Persistir no localStorage
    localStorage.setItem('@LEP:user', JSON.stringify(response.user));
    localStorage.setItem('@LEP:token', response.token);
    localStorage.setItem('@LEP:organizations', JSON.stringify(response.organizations));
    localStorage.setItem('@LEP:projects', JSON.stringify(response.projects));
    localStorage.setItem('@LEP:currentOrganization', response.organizations[0]?.organization_id);
    localStorage.setItem('@LEP:currentProject', response.projects[0]?.project_id);
  };

  return (
    <AuthContext.Provider value={{
      user,
      token,
      organizations,
      projects,
      currentOrganization,
      currentProject,
      setCurrentOrganization,
      setCurrentProject,
      login,
      logout,
    }}>
      {children}
    </AuthContext.Provider>
  );
};
```

### 2. Atualizar Axios Interceptor

```typescript
// api.ts

axios.interceptors.request.use((config) => {
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
```

### 3. Criar Componente de Seleção Organização/Projeto

```typescript
// OrganizationProjectSelector.tsx

import { useAuth } from '@/contexts/auth-context';
import { Select } from '@/components/ui/select';

export const OrganizationProjectSelector = () => {
  const {
    organizations,
    projects,
    currentOrganization,
    currentProject,
    setCurrentOrganization,
    setCurrentProject
  } = useAuth();

  const handleOrgChange = (orgId: string) => {
    setCurrentOrganization(orgId);

    // Filtrar projetos da organização selecionada
    const orgProjects = projects.filter(p => {
      // Buscar projeto no backend para verificar organization_id
      // ou manter um cache local
      return true; // Simplificado
    });

    if (orgProjects.length > 0) {
      setCurrentProject(orgProjects[0].project_id);
    }
  };

  const handleProjectChange = (projectId: string) => {
    setCurrentProject(projectId);
  };

  return (
    <div className="flex gap-4">
      <Select
        value={currentOrganization}
        onValueChange={handleOrgChange}
      >
        {organizations.map(org => (
          <option key={org.id} value={org.organization_id}>
            {/* Buscar nome da organização via API ou cache */}
            Organização {org.organization_id}
          </option>
        ))}
      </Select>

      <Select
        value={currentProject}
        onValueChange={handleProjectChange}
      >
        {projects.map(proj => (
          <option key={proj.id} value={proj.project_id}>
            {/* Buscar nome do projeto via API ou cache */}
            Projeto {proj.project_id}
          </option>
        ))}
      </Select>
    </div>
  );
};
```

### 4. Criar Serviços para Novos Endpoints

```typescript
// services/userOrganizationService.ts

import api from './api';

export const userOrganizationService = {
  addUserToOrganization: async (userId: string, data: {
    organization_id: string;
    role: string;
  }) => {
    const response = await api.post(`/user/${userId}/organization`, data);
    return response.data;
  },

  removeUserFromOrganization: async (userId: string, orgId: string) => {
    await api.delete(`/user/${userId}/organization/${orgId}`);
  },

  getUserOrganizations: async (userId: string) => {
    const response = await api.get(`/user/${userId}/organizations`);
    return response.data;
  },

  getOrganizationUsers: async (orgId: string) => {
    const response = await api.get(`/organization/${orgId}/users`);
    return response.data;
  },
};

// services/userProjectService.ts

export const userProjectService = {
  addUserToProject: async (userId: string, data: {
    project_id: string;
    role: string;
  }) => {
    const response = await api.post(`/user/${userId}/project`, data);
    return response.data;
  },

  removeUserFromProject: async (userId: string, projectId: string) => {
    await api.delete(`/user/${userId}/project/${projectId}`);
  },

  getUserProjects: async (userId: string) => {
    const response = await api.get(`/user/${userId}/projects`);
    return response.data;
  },

  getUserProjectsByOrganization: async (userId: string, orgId: string) => {
    const response = await api.get(`/user/${userId}/organization/${orgId}/projects`);
    return response.data;
  },

  getProjectUsers: async (projectId: string) => {
    const response = await api.get(`/project/${projectId}/users`);
    return response.data;
  },
};
```

### 5. Adicionar Seletor no Layout

```typescript
// layouts/DashboardLayout.tsx

import { OrganizationProjectSelector } from '@/components/OrganizationProjectSelector';

export const DashboardLayout = ({ children }) => {
  return (
    <div>
      <header>
        <nav>
          {/* Logo, Menu, etc */}
        </nav>

        {/* Adicionar seletor aqui */}
        <OrganizationProjectSelector />
      </header>

      <main>{children}</main>
    </div>
  );
};
```

## Fluxo de Uso

1. **Login**: Usuário faz login e recebe lista de organizações e projetos
2. **Seleção Automática**: Sistema seleciona automaticamente primeira org/projeto
3. **Troca de Contexto**: Usuário pode trocar org/projeto usando os selects
4. **Atualização de Headers**: Headers são atualizados automaticamente no interceptor
5. **Requisições**: Todas as requisições usam a org/projeto selecionados

## Benefícios

- ✅ Usuário pode acessar múltiplas organizações/projetos
- ✅ Não precisa fazer logout para trocar de contexto
- ✅ Roles diferentes por org/projeto
- ✅ Sistema mais flexível e escalável
- ✅ Melhor controle de permissões

## Breaking Changes

⚠️ **ATENÇÃO**: Estas mudanças quebram compatibilidade com a versão anterior:

1. Campo `organization_id` e `project_id` foram removidos do modelo `User`
2. Resposta de login agora inclui arrays `organizations` e `projects`
3. Token JWT agora inclui `user_id` nos claims
4. Middleware valida se usuário tem acesso à org/projeto nos headers

## Migração de Dados

Se você tem dados existentes, execute:

1. Backup do banco de dados
2. Rode as migrations (automático ao iniciar o backend)
3. Execute um script de migração para criar relacionamentos UserOrganization e UserProject para usuários existentes

## Suporte

Para dúvidas ou problemas, consulte:
- [CLAUDE.md](./CLAUDE.md) - Documentação do backend
- [README.md](./README.md) - Documentação geral
