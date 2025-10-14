-- Script para criar Master Admins diretamente no banco
-- Senha: senha123 (já hasheada com bcrypt)

-- 1. Verificar e inserir organização
INSERT INTO organizations (id, name, email, phone, address, website, description, active, created_at, updated_at)
VALUES (
    '123e4567-e89b-12d3-a456-426614174000',
    'LEP Restaurante Demo',
    'teste1@gmail.com',
    '+55 11 9999-8888',
    'Rua das Flores, 123 - São Paulo, SP',
    'https://lep-demo.com',
    'Restaurante demo para demonstração do sistema LEP',
    true,
    NOW(),
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    updated_at = NOW();

-- 2. Inserir Master Admins
INSERT INTO users (id, name, email, password, permissions, active, created_at, updated_at)
VALUES
    (
        '123e4567-e89b-12d3-a456-426614174010',
        'Pablo Master Admin',
        'pablo@lep.com',
        '$2a$10$C83MMWJFkG/djLU.UfWEZuQ4Xl2gJPz.ABP//wEsbzBggjfRV4kF.',
        ARRAY['master_admin']::text[],
        true,
        NOW(),
        NOW()
    ),
    (
        '123e4567-e89b-12d3-a456-426614174011',
        'Luan Master Admin',
        'luan@lep.com',
        '$2a$10$C83MMWJFkG/djLU.UfWEZuQ4Xl2gJPz.ABP//wEsbzBggjfRV4kF.',
        ARRAY['master_admin']::text[],
        true,
        NOW(),
        NOW()
    ),
    (
        '123e4567-e89b-12d3-a456-426614174012',
        'Eduardo Master Admin',
        'eduardo@lep.com',
        '$2a$10$C83MMWJFkG/djLU.UfWEZuQ4Xl2gJPz.ABP//wEsbzBggjfRV4kF.',
        ARRAY['master_admin']::text[],
        true,
        NOW(),
        NOW()
    )
ON CONFLICT (email) DO UPDATE SET
    name = EXCLUDED.name,
    password = EXCLUDED.password,
    permissions = EXCLUDED.permissions,
    updated_at = NOW();

-- 3. Criar relacionamentos usuário-organização
INSERT INTO user_organizations (id, user_id, organization_id, role, active, created_at, updated_at)
VALUES
    (
        '223e4567-e89b-12d3-a456-426614174010',
        '123e4567-e89b-12d3-a456-426614174010',
        '123e4567-e89b-12d3-a456-426614174000',
        'owner',
        true,
        NOW(),
        NOW()
    ),
    (
        '223e4567-e89b-12d3-a456-426614174011',
        '123e4567-e89b-12d3-a456-426614174011',
        '123e4567-e89b-12d3-a456-426614174000',
        'admin',
        true,
        NOW(),
        NOW()
    ),
    (
        '223e4567-e89b-12d3-a456-426614174012',
        '123e4567-e89b-12d3-a456-426614174012',
        '123e4567-e89b-12d3-a456-426614174000',
        'admin',
        true,
        NOW(),
        NOW()
    )
ON CONFLICT (id) DO UPDATE SET
    role = EXCLUDED.role,
    active = EXCLUDED.active,
    updated_at = NOW();

-- Mensagem de sucesso
SELECT '✅ Master Admins criados com sucesso!' AS status;
SELECT
    u.name,
    u.email,
    u.permissions,
    uo.role AS org_role
FROM users u
LEFT JOIN user_organizations uo ON u.id = uo.user_id
WHERE u.permissions @> ARRAY['master_admin']::text[];
