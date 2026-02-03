-- Migration: Alterar constraint de email único global para único por organização
-- Data: 2026-02-02
-- Descrição: Permite que o mesmo email exista em organizações diferentes

-- Remove constraint global (se existir)
ALTER TABLE clients DROP CONSTRAINT IF EXISTS uni_clients_email;

-- Remove índice antigo (se existir)
DROP INDEX IF EXISTS uni_clients_email;

-- Cria índice único composto (email + org_id) excluindo registros deletados
CREATE UNIQUE INDEX IF NOT EXISTS idx_client_email_org
ON clients (email, org_id) WHERE deleted_at IS NULL;
