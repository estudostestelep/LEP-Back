package migrate

import (
	"fmt"

	"gorm.io/gorm"
)

// migrateFileReferences cria as tabelas de gerenciamento de imagens com deduplicação
func (r *resourceMigrate) migrateFileReferences() {
	migrator := r.db.Migrator()

	// Criar tabela file_references se não existir
	if !migrator.HasTable("file_references") {
		fmt.Println("📋 Criando tabela file_references...")
		sql := `
		CREATE TABLE file_references (
			id UUID PRIMARY KEY,
			organization_id UUID NOT NULL,
			project_id UUID NOT NULL,
			file_hash VARCHAR(64) NOT NULL,
			file_path VARCHAR(512) NOT NULL,
			file_size BIGINT NOT NULL,
			category VARCHAR(50) NOT NULL,
			mime_type VARCHAR(50) NOT NULL,
			reference_count INTEGER DEFAULT 1 NOT NULL,
			created_at TIMESTAMP NOT NULL,
			last_accessed_at TIMESTAMP,
			deleted_at TIMESTAMP,
			UNIQUE(organization_id, project_id, file_hash)
		);`

		if err := r.db.Exec(sql).Error; err != nil {
			fmt.Printf("⚠️  Erro ao criar tabela file_references: %v\n", err)
			return
		}

		// Criar índices
		indexSQL := []string{
			"CREATE INDEX idx_file_references_org_proj ON file_references(organization_id, project_id);",
			"CREATE INDEX idx_file_references_file_hash ON file_references(file_hash);",
			"CREATE INDEX idx_file_references_deleted ON file_references(deleted_at);",
		}

		for _, sql := range indexSQL {
			if err := r.db.Exec(sql).Error; err != nil {
				fmt.Printf("⚠️  Aviso ao criar índice: %v\n", err)
			}
		}
		fmt.Println("✅ Tabela file_references criada com sucesso")
	}

	// Criar tabela entity_file_references se não existir
	if !migrator.HasTable("entity_file_references") {
		fmt.Println("📋 Criando tabela entity_file_references...")
		sql := `
		CREATE TABLE entity_file_references (
			id UUID PRIMARY KEY,
			file_id UUID NOT NULL,
			entity_type VARCHAR(50) NOT NULL,
			entity_id UUID NOT NULL,
			entity_field VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			FOREIGN KEY (file_id) REFERENCES file_references(id),
			UNIQUE(entity_type, entity_id, entity_field)
		);`

		if err := r.db.Exec(sql).Error; err != nil {
			fmt.Printf("⚠️  Erro ao criar tabela entity_file_references: %v\n", err)
			return
		}

		// Criar índices
		indexSQL := []string{
			"CREATE INDEX idx_entity_file_references_file_id ON entity_file_references(file_id);",
			"CREATE INDEX idx_entity_file_references_entity ON entity_file_references(entity_type, entity_id);",
			"CREATE INDEX idx_entity_file_references_deleted ON entity_file_references(deleted_at);",
		}

		for _, sql := range indexSQL {
			if err := r.db.Exec(sql).Error; err != nil {
				fmt.Printf("⚠️  Aviso ao criar índice: %v\n", err)
			}
		}
		fmt.Println("✅ Tabela entity_file_references criada com sucesso")
	}

	// Adicionar coluna file_hash nas entidades se não existir (retrocompatibilidade futura)
	// Esta será preenchida gradualmente conforme imagens são reutilizadas
	addFileHashColumns(r.db)
}

// addFileHashColumns adiciona coluna file_hash em tabelas que armazenam imagens
func addFileHashColumns(db *gorm.DB) {
	migrator := db.Migrator()

	// Tabelas que podem ter imagens
	tablesToUpdate := []string{
		"products",
		"categories",
		"menus",
		"subcategories",
	}

	for _, tableName := range tablesToUpdate {
		if migrator.HasTable(tableName) && !migrator.HasColumn(tableName, "file_hash") {
			sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN file_hash VARCHAR(64);", tableName)
			if err := db.Exec(sql).Error; err != nil {
				fmt.Printf("⚠️  Aviso ao adicionar file_hash em %s: %v\n", tableName, err)
			}
		}
	}
}
