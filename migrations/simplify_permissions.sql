-- ============================================================
-- MIGRAÇÃO: Simplificação do Sistema de Permissionamento
-- ============================================================
-- Este script migra o sistema de permissões para o novo formato:
-- - Permissões no formato module:action (ex: orders:read)
-- - Tabelas renomeadas: packages -> plans
-- - Pivot simplificada: role_permission_levels -> role_permissions
-- - Remoção do campo Permissions dos modelos de usuário
-- ============================================================

-- ============================================================
-- FASE 1: Criar novas tabelas
-- ============================================================

-- 1.1 Criar tabela plans (nova versão de packages)
CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(100),
    description VARCHAR(500),
    price_monthly DECIMAL(10,2) DEFAULT 0,
    price_yearly DECIMAL(10,2) DEFAULT 0,
    is_public BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_plans_active ON plans(active);
CREATE INDEX IF NOT EXISTS idx_plans_deleted_at ON plans(deleted_at);

-- 1.2 Criar tabela plan_modules (pivot Plan-Module)
CREATE TABLE IF NOT EXISTS plan_modules (
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (plan_id, module_id)
);

-- 1.3 Criar tabela plan_limits
CREATE TABLE IF NOT EXISTS plan_limits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    limit_type VARCHAR(50) NOT NULL,
    limit_value INT DEFAULT -1, -- -1 = ilimitado
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_plan_limits_plan_id ON plan_limits(plan_id);

-- 1.4 Criar tabela organization_plans
CREATE TABLE IF NOT EXISTS organization_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL UNIQUE REFERENCES organizations(id) ON DELETE CASCADE,
    plan_id UUID NOT NULL REFERENCES plans(id),
    billing_cycle VARCHAR(20) DEFAULT 'monthly',
    custom_price DECIMAL(10,2),
    starts_at TIMESTAMP,
    expires_at TIMESTAMP,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_org_plans_plan_id ON organization_plans(plan_id);
CREATE INDEX IF NOT EXISTS idx_org_plans_active ON organization_plans(active);

-- 1.5 Criar tabela role_permissions (pivot simplificada)
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

-- 1.6 Criar tabela client_roles (pivot Client-Role)
CREATE TABLE IF NOT EXISTS client_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    project_id UUID,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_client_roles_client_id ON client_roles(client_id);
CREATE INDEX IF NOT EXISTS idx_client_roles_role_id ON client_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_client_roles_org_id ON client_roles(organization_id);

-- ============================================================
-- FASE 2: Adicionar colunas às tabelas existentes
-- ============================================================

-- 2.1 Adicionar colunas module e action à tabela permissions
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS module VARCHAR(50);
ALTER TABLE permissions ADD COLUMN IF NOT EXISTS action VARCHAR(20);

-- 2.2 Renomear code_name para code na tabela permissions (se existir)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name = 'permissions' AND column_name = 'code_name') THEN
        ALTER TABLE permissions RENAME COLUMN code_name TO code;
    END IF;
END $$;

-- 2.3 Renomear code_name para code na tabela modules (se existir)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name = 'modules' AND column_name = 'code_name') THEN
        ALTER TABLE modules RENAME COLUMN code_name TO code;
    END IF;
    IF EXISTS (SELECT 1 FROM information_schema.columns
               WHERE table_name = 'modules' AND column_name = 'display_name') THEN
        ALTER TABLE modules RENAME COLUMN display_name TO name;
    END IF;
END $$;

-- ============================================================
-- FASE 3: Migrar dados de packages para plans
-- ============================================================

-- 3.1 Migrar packages para plans
INSERT INTO plans (id, code, name, description, price_monthly, price_yearly, is_public, display_order, active, created_at, updated_at)
SELECT id, code_name, display_name, description, price_monthly, price_yearly, is_public, display_order, active, created_at, updated_at
FROM packages
WHERE deleted_at IS NULL
ON CONFLICT (id) DO NOTHING;

-- 3.2 Migrar package_modules para plan_modules
INSERT INTO plan_modules (plan_id, module_id, created_at)
SELECT package_id, module_id, created_at
FROM package_modules
WHERE deleted_at IS NULL
ON CONFLICT (plan_id, module_id) DO NOTHING;

-- 3.3 Migrar package_limits para plan_limits
INSERT INTO plan_limits (id, plan_id, limit_type, limit_value, created_at, updated_at)
SELECT id, package_id, limit_type, limit_value, created_at, updated_at
FROM package_limits
WHERE deleted_at IS NULL
ON CONFLICT (id) DO NOTHING;

-- 3.4 Migrar organization_packages para organization_plans
INSERT INTO organization_plans (id, organization_id, plan_id, billing_cycle, custom_price, starts_at, expires_at, active, created_at, updated_at)
SELECT id, organization_id, package_id, billing_cycle, custom_price, started_at, expires_at, active, created_at, updated_at
FROM organization_packages
WHERE deleted_at IS NULL
ON CONFLICT (organization_id) DO NOTHING;

-- ============================================================
-- FASE 4: Migrar role_permission_levels para role_permissions
-- ============================================================

-- Migrar apenas níveis > 0 (com acesso)
INSERT INTO role_permissions (role_id, permission_id, created_at)
SELECT role_id, permission_id, created_at
FROM role_permission_levels
WHERE deleted_at IS NULL AND level > 0
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ============================================================
-- FASE 5: Atualizar permissões para novo formato module:action
-- ============================================================

-- 5.1 Mapeamento de permissões antigas para novas
-- Orders
UPDATE permissions SET code = 'orders:read', module = 'orders', action = 'read' WHERE code = 'client_orders_view' OR code = 'view_orders';
UPDATE permissions SET code = 'orders:create', module = 'orders', action = 'create' WHERE code = 'client_orders_create';
UPDATE permissions SET code = 'orders:update', module = 'orders', action = 'update' WHERE code = 'client_orders_edit' OR code = 'manage_orders';
UPDATE permissions SET code = 'orders:delete', module = 'orders', action = 'delete' WHERE code = 'client_orders_delete';

-- Menu
UPDATE permissions SET code = 'menu:read', module = 'menu', action = 'read' WHERE code = 'client_menu_view' OR code = 'view_menus';
UPDATE permissions SET code = 'menu:create', module = 'menu', action = 'create' WHERE code = 'client_menu_create';
UPDATE permissions SET code = 'menu:update', module = 'menu', action = 'update' WHERE code = 'client_menu_edit' OR code = 'manage_menus';
UPDATE permissions SET code = 'menu:delete', module = 'menu', action = 'delete' WHERE code = 'client_menu_delete';

-- Products
UPDATE permissions SET code = 'products:read', module = 'products', action = 'read' WHERE code = 'client_products_view' OR code = 'view_products';
UPDATE permissions SET code = 'products:create', module = 'products', action = 'create' WHERE code = 'client_products_create';
UPDATE permissions SET code = 'products:update', module = 'products', action = 'update' WHERE code = 'client_products_edit' OR code = 'manage_products';
UPDATE permissions SET code = 'products:delete', module = 'products', action = 'delete' WHERE code = 'client_products_delete';

-- Tables
UPDATE permissions SET code = 'tables:read', module = 'tables', action = 'read' WHERE code = 'client_tables_view' OR code = 'view_tables';
UPDATE permissions SET code = 'tables:create', module = 'tables', action = 'create' WHERE code = 'client_tables_create';
UPDATE permissions SET code = 'tables:update', module = 'tables', action = 'update' WHERE code = 'client_tables_edit' OR code = 'manage_tables';
UPDATE permissions SET code = 'tables:delete', module = 'tables', action = 'delete' WHERE code = 'client_tables_delete';

-- Reservations
UPDATE permissions SET code = 'reservations:read', module = 'reservations', action = 'read' WHERE code = 'client_reservations_view' OR code = 'view_reservations';
UPDATE permissions SET code = 'reservations:create', module = 'reservations', action = 'create' WHERE code = 'client_reservations_create';
UPDATE permissions SET code = 'reservations:update', module = 'reservations', action = 'update' WHERE code = 'client_reservations_edit' OR code = 'manage_reservations';
UPDATE permissions SET code = 'reservations:delete', module = 'reservations', action = 'delete' WHERE code = 'client_reservations_delete';

-- Waitlist
UPDATE permissions SET code = 'waitlist:read', module = 'waitlist', action = 'read' WHERE code = 'client_waitlist_view' OR code = 'view_waitlists';
UPDATE permissions SET code = 'waitlist:create', module = 'waitlist', action = 'create' WHERE code = 'client_waitlist_create';
UPDATE permissions SET code = 'waitlist:update', module = 'waitlist', action = 'update' WHERE code = 'client_waitlist_edit' OR code = 'manage_waitlists';
UPDATE permissions SET code = 'waitlist:delete', module = 'waitlist', action = 'delete' WHERE code = 'client_waitlist_delete';

-- Customers
UPDATE permissions SET code = 'customers:read', module = 'customers', action = 'read' WHERE code = 'client_customers_view' OR code = 'view_customers';
UPDATE permissions SET code = 'customers:create', module = 'customers', action = 'create' WHERE code = 'client_customers_create';
UPDATE permissions SET code = 'customers:update', module = 'customers', action = 'update' WHERE code = 'client_customers_edit' OR code = 'manage_customers';
UPDATE permissions SET code = 'customers:delete', module = 'customers', action = 'delete' WHERE code = 'client_customers_delete';

-- Users
UPDATE permissions SET code = 'users:read', module = 'users', action = 'read' WHERE code = 'client_users_view' OR code = 'view_users';
UPDATE permissions SET code = 'users:create', module = 'users', action = 'create' WHERE code = 'client_users_create' OR code = 'manage_users';
UPDATE permissions SET code = 'users:update', module = 'users', action = 'update' WHERE code = 'client_users_edit';
UPDATE permissions SET code = 'users:delete', module = 'users', action = 'delete' WHERE code = 'client_users_delete';

-- Reports
UPDATE permissions SET code = 'reports:read', module = 'reports', action = 'read' WHERE code = 'client_reports_view' OR code = 'view_reports';
UPDATE permissions SET code = 'reports:export', module = 'reports', action = 'export' WHERE code = 'client_reports_export' OR code = 'export_data';

-- Settings
UPDATE permissions SET code = 'settings:read', module = 'settings', action = 'read' WHERE code = 'client_settings_view' OR code = 'view_settings';
UPDATE permissions SET code = 'settings:update', module = 'settings', action = 'update' WHERE code = 'client_settings_edit' OR code = 'manage_settings';

-- Notifications
UPDATE permissions SET code = 'notifications:read', module = 'notifications', action = 'read' WHERE code = 'client_notifications_view' OR code = 'view_notifications';
UPDATE permissions SET code = 'notifications:create', module = 'notifications', action = 'create' WHERE code = 'client_notifications_create';
UPDATE permissions SET code = 'notifications:update', module = 'notifications', action = 'update' WHERE code = 'client_notifications_edit' OR code = 'manage_notifications';
UPDATE permissions SET code = 'notifications:delete', module = 'notifications', action = 'delete' WHERE code = 'client_notifications_delete';
UPDATE permissions SET code = 'notifications:send', module = 'notifications', action = 'send' WHERE code = 'client_notifications_send';

-- Admin - Organizations
UPDATE permissions SET code = 'organizations:read', module = 'organizations', action = 'read' WHERE code = 'admin_organizations_view' OR code = 'view_organizations';
UPDATE permissions SET code = 'organizations:create', module = 'organizations', action = 'create' WHERE code = 'admin_organizations_create' OR code = 'manage_organizations';
UPDATE permissions SET code = 'organizations:update', module = 'organizations', action = 'update' WHERE code = 'admin_organizations_edit';
UPDATE permissions SET code = 'organizations:delete', module = 'organizations', action = 'delete' WHERE code = 'admin_organizations_delete';

-- Admin - Plans/Packages
UPDATE permissions SET code = 'plans:read', module = 'plans', action = 'read' WHERE code = 'admin_packages_view';
UPDATE permissions SET code = 'plans:create', module = 'plans', action = 'create' WHERE code = 'admin_packages_create';
UPDATE permissions SET code = 'plans:update', module = 'plans', action = 'update' WHERE code = 'admin_packages_edit';
UPDATE permissions SET code = 'plans:delete', module = 'plans', action = 'delete' WHERE code = 'admin_packages_delete';

-- 5.2 Atualizar module e action para permissões que já têm o formato correto
UPDATE permissions
SET module = SPLIT_PART(code, ':', 1),
    action = SPLIT_PART(code, ':', 2)
WHERE code LIKE '%:%' AND (module IS NULL OR action IS NULL);

-- ============================================================
-- FASE 6: Criar role master_admin se não existir
-- ============================================================

INSERT INTO roles (id, name, display_name, description, hierarchy_level, scope, is_system, active, created_at, updated_at)
VALUES (
    gen_random_uuid(),
    'master_admin',
    'Master Admin',
    'Acesso total ao sistema',
    10,
    'admin',
    true,
    true,
    NOW(),
    NOW()
)
ON CONFLICT (name) DO UPDATE SET hierarchy_level = 10;

-- ============================================================
-- FASE 7: Migrar user_roles para client_roles (para clientes)
-- ============================================================

-- Migrar apenas registros que referenciam clients (não admins)
INSERT INTO client_roles (id, client_id, role_id, organization_id, project_id, active, created_at, updated_at)
SELECT
    ur.id,
    ur.user_id,
    ur.role_id,
    COALESCE(ur.organization_id, c.org_id),
    ur.project_id,
    ur.active,
    ur.created_at,
    ur.updated_at
FROM user_roles ur
INNER JOIN clients c ON c.id = ur.user_id
WHERE ur.deleted_at IS NULL
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- FASE 8: Limpar dados antigos (OPCIONAL - executar manualmente)
-- ============================================================

-- ATENÇÃO: Executar apenas após validar que a migração funcionou!
-- Descomentar as linhas abaixo quando estiver pronto para limpar

-- -- Remover coluna permissions das tabelas de usuário
-- ALTER TABLE admins DROP COLUMN IF EXISTS permissions;
-- ALTER TABLE clients DROP COLUMN IF EXISTS permissions;
-- ALTER TABLE roles DROP COLUMN IF EXISTS permissions;

-- -- Remover tabelas antigas
-- DROP TABLE IF EXISTS role_permission_levels CASCADE;
-- DROP TABLE IF EXISTS user_roles CASCADE;
-- DROP TABLE IF EXISTS organization_packages CASCADE;
-- DROP TABLE IF EXISTS package_limits CASCADE;
-- DROP TABLE IF EXISTS package_modules CASCADE;
-- DROP TABLE IF EXISTS package_bundles CASCADE;
-- DROP TABLE IF EXISTS packages CASCADE;

-- ============================================================
-- FASE 9: Inserir permissões padrão no novo formato
-- ============================================================

-- Inserir permissões que podem não existir
INSERT INTO permissions (id, code, module, action, display_name, description, active, created_at, updated_at)
VALUES
    -- Orders
    (gen_random_uuid(), 'orders:read', 'orders', 'read', 'Visualizar Pedidos', 'Permite visualizar pedidos', true, NOW(), NOW()),
    (gen_random_uuid(), 'orders:create', 'orders', 'create', 'Criar Pedidos', 'Permite criar novos pedidos', true, NOW(), NOW()),
    (gen_random_uuid(), 'orders:update', 'orders', 'update', 'Editar Pedidos', 'Permite editar pedidos existentes', true, NOW(), NOW()),
    (gen_random_uuid(), 'orders:delete', 'orders', 'delete', 'Excluir Pedidos', 'Permite excluir pedidos', true, NOW(), NOW()),
    -- Menu
    (gen_random_uuid(), 'menu:read', 'menu', 'read', 'Visualizar Cardápio', 'Permite visualizar o cardápio', true, NOW(), NOW()),
    (gen_random_uuid(), 'menu:create', 'menu', 'create', 'Criar Cardápio', 'Permite criar itens no cardápio', true, NOW(), NOW()),
    (gen_random_uuid(), 'menu:update', 'menu', 'update', 'Editar Cardápio', 'Permite editar itens do cardápio', true, NOW(), NOW()),
    (gen_random_uuid(), 'menu:delete', 'menu', 'delete', 'Excluir Cardápio', 'Permite excluir itens do cardápio', true, NOW(), NOW()),
    -- Products
    (gen_random_uuid(), 'products:read', 'products', 'read', 'Visualizar Produtos', 'Permite visualizar produtos', true, NOW(), NOW()),
    (gen_random_uuid(), 'products:create', 'products', 'create', 'Criar Produtos', 'Permite criar produtos', true, NOW(), NOW()),
    (gen_random_uuid(), 'products:update', 'products', 'update', 'Editar Produtos', 'Permite editar produtos', true, NOW(), NOW()),
    (gen_random_uuid(), 'products:delete', 'products', 'delete', 'Excluir Produtos', 'Permite excluir produtos', true, NOW(), NOW()),
    -- Tables
    (gen_random_uuid(), 'tables:read', 'tables', 'read', 'Visualizar Mesas', 'Permite visualizar mesas', true, NOW(), NOW()),
    (gen_random_uuid(), 'tables:create', 'tables', 'create', 'Criar Mesas', 'Permite criar mesas', true, NOW(), NOW()),
    (gen_random_uuid(), 'tables:update', 'tables', 'update', 'Editar Mesas', 'Permite editar mesas', true, NOW(), NOW()),
    (gen_random_uuid(), 'tables:delete', 'tables', 'delete', 'Excluir Mesas', 'Permite excluir mesas', true, NOW(), NOW()),
    -- Reservations
    (gen_random_uuid(), 'reservations:read', 'reservations', 'read', 'Visualizar Reservas', 'Permite visualizar reservas', true, NOW(), NOW()),
    (gen_random_uuid(), 'reservations:create', 'reservations', 'create', 'Criar Reservas', 'Permite criar reservas', true, NOW(), NOW()),
    (gen_random_uuid(), 'reservations:update', 'reservations', 'update', 'Editar Reservas', 'Permite editar reservas', true, NOW(), NOW()),
    (gen_random_uuid(), 'reservations:delete', 'reservations', 'delete', 'Excluir Reservas', 'Permite excluir reservas', true, NOW(), NOW()),
    -- Customers
    (gen_random_uuid(), 'customers:read', 'customers', 'read', 'Visualizar Clientes', 'Permite visualizar clientes', true, NOW(), NOW()),
    (gen_random_uuid(), 'customers:create', 'customers', 'create', 'Criar Clientes', 'Permite criar clientes', true, NOW(), NOW()),
    (gen_random_uuid(), 'customers:update', 'customers', 'update', 'Editar Clientes', 'Permite editar clientes', true, NOW(), NOW()),
    (gen_random_uuid(), 'customers:delete', 'customers', 'delete', 'Excluir Clientes', 'Permite excluir clientes', true, NOW(), NOW()),
    -- Users
    (gen_random_uuid(), 'users:read', 'users', 'read', 'Visualizar Usuários', 'Permite visualizar usuários', true, NOW(), NOW()),
    (gen_random_uuid(), 'users:create', 'users', 'create', 'Criar Usuários', 'Permite criar usuários', true, NOW(), NOW()),
    (gen_random_uuid(), 'users:update', 'users', 'update', 'Editar Usuários', 'Permite editar usuários', true, NOW(), NOW()),
    (gen_random_uuid(), 'users:delete', 'users', 'delete', 'Excluir Usuários', 'Permite excluir usuários', true, NOW(), NOW()),
    -- Reports
    (gen_random_uuid(), 'reports:read', 'reports', 'read', 'Visualizar Relatórios', 'Permite visualizar relatórios', true, NOW(), NOW()),
    (gen_random_uuid(), 'reports:export', 'reports', 'export', 'Exportar Relatórios', 'Permite exportar relatórios', true, NOW(), NOW()),
    -- Settings
    (gen_random_uuid(), 'settings:read', 'settings', 'read', 'Visualizar Configurações', 'Permite visualizar configurações', true, NOW(), NOW()),
    (gen_random_uuid(), 'settings:update', 'settings', 'update', 'Editar Configurações', 'Permite editar configurações', true, NOW(), NOW()),
    -- Notifications
    (gen_random_uuid(), 'notifications:read', 'notifications', 'read', 'Visualizar Notificações', 'Permite visualizar notificações', true, NOW(), NOW()),
    (gen_random_uuid(), 'notifications:create', 'notifications', 'create', 'Criar Notificações', 'Permite criar templates', true, NOW(), NOW()),
    (gen_random_uuid(), 'notifications:update', 'notifications', 'update', 'Editar Notificações', 'Permite editar templates', true, NOW(), NOW()),
    (gen_random_uuid(), 'notifications:delete', 'notifications', 'delete', 'Excluir Notificações', 'Permite excluir templates', true, NOW(), NOW()),
    (gen_random_uuid(), 'notifications:send', 'notifications', 'send', 'Enviar Notificações', 'Permite enviar notificações', true, NOW(), NOW()),
    -- Organizations (Admin)
    (gen_random_uuid(), 'organizations:read', 'organizations', 'read', 'Visualizar Organizações', 'Permite visualizar organizações', true, NOW(), NOW()),
    (gen_random_uuid(), 'organizations:create', 'organizations', 'create', 'Criar Organizações', 'Permite criar organizações', true, NOW(), NOW()),
    (gen_random_uuid(), 'organizations:update', 'organizations', 'update', 'Editar Organizações', 'Permite editar organizações', true, NOW(), NOW()),
    (gen_random_uuid(), 'organizations:delete', 'organizations', 'delete', 'Excluir Organizações', 'Permite excluir organizações', true, NOW(), NOW()),
    -- Plans (Admin)
    (gen_random_uuid(), 'plans:read', 'plans', 'read', 'Visualizar Planos', 'Permite visualizar planos', true, NOW(), NOW()),
    (gen_random_uuid(), 'plans:create', 'plans', 'create', 'Criar Planos', 'Permite criar planos', true, NOW(), NOW()),
    (gen_random_uuid(), 'plans:update', 'plans', 'update', 'Editar Planos', 'Permite editar planos', true, NOW(), NOW()),
    (gen_random_uuid(), 'plans:delete', 'plans', 'delete', 'Excluir Planos', 'Permite excluir planos', true, NOW(), NOW()),
    -- Waitlist
    (gen_random_uuid(), 'waitlist:read', 'waitlist', 'read', 'Visualizar Fila de Espera', 'Permite visualizar fila de espera', true, NOW(), NOW()),
    (gen_random_uuid(), 'waitlist:create', 'waitlist', 'create', 'Criar Fila de Espera', 'Permite adicionar à fila', true, NOW(), NOW()),
    (gen_random_uuid(), 'waitlist:update', 'waitlist', 'update', 'Editar Fila de Espera', 'Permite editar fila', true, NOW(), NOW()),
    (gen_random_uuid(), 'waitlist:delete', 'waitlist', 'delete', 'Excluir Fila de Espera', 'Permite remover da fila', true, NOW(), NOW())
ON CONFLICT (code) DO NOTHING;

-- ============================================================
-- FIM DA MIGRAÇÃO
-- ============================================================

-- Verificar resultados
SELECT 'Planos migrados:' as info, COUNT(*) as count FROM plans;
SELECT 'Permissões atualizadas:' as info, COUNT(*) as count FROM permissions WHERE module IS NOT NULL;
SELECT 'Role permissions:' as info, COUNT(*) as count FROM role_permissions;
SELECT 'Client roles:' as info, COUNT(*) as count FROM client_roles;
SELECT 'Organization plans:' as info, COUNT(*) as count FROM organization_plans;
