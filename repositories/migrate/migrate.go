package migrate

import (
	"fmt"

	"gorm.io/gorm"
)

type resourceMigrate struct {
	db *gorm.DB
}

type IMigrate interface {
	MigrateRun(modelsToMigrate ...interface{})
}

func (r *resourceMigrate) MigrateRun(modelsToMigrate ...interface{}) {
	// Migração customizada para novos campos de Product
	r.migrateProductFields()

	// Migração de tabelas de referência de imagens (file references)
	r.migrateFileReferences()

	// Migração automática para novas tabelas
	for _, model := range modelsToMigrate {
		migrator := r.db.Migrator()
		if !migrator.HasTable(model) {
			if err := r.db.AutoMigrate(model); err != nil {
				panic(fmt.Sprintf("erro ao migrar tabela %T: %v", model, err))
			}
		}
	}

	// AutoMigrate para adicionar novos campos em tabelas existentes
	if err := r.db.AutoMigrate(modelsToMigrate...); err != nil {
		panic(fmt.Sprintf("erro na migrate AutoMigrate: %v", err))
	}
}

// migrateProductFields adiciona os novos campos da tabela products de forma segura
func (r *resourceMigrate) migrateProductFields() {
	migrator := r.db.Migrator()

	// Verificar se a tabela products existe
	if !migrator.HasTable("products") {
		return // Tabela será criada pelo AutoMigrate principal
	}

	// Adicionar campo 'type' se não existir
	if !migrator.HasColumn("products", "type") {
		// 1. Adicionar coluna como nullable
		if err := r.db.Exec("ALTER TABLE products ADD COLUMN type TEXT").Error; err != nil {
			fmt.Printf("⚠️  Aviso ao adicionar coluna 'type': %v\n", err)
		}

		// 2. Preencher valores existentes com 'prato' como padrão
		if err := r.db.Exec("UPDATE products SET type = 'prato' WHERE type IS NULL").Error; err != nil {
			fmt.Printf("⚠️  Aviso ao atualizar coluna 'type': %v\n", err)
		}

		// 3. Tornar coluna NOT NULL
		if err := r.db.Exec("ALTER TABLE products ALTER COLUMN type SET NOT NULL").Error; err != nil {
			fmt.Printf("⚠️  Aviso ao tornar coluna 'type' NOT NULL: %v\n", err)
		}
	}

	// Adicionar outros novos campos (todos opcionais, sem problemas)
	newColumns := map[string]string{
		`"order"`:            "INTEGER DEFAULT 0",          // Aspas por ser palavra reservada
		"active":             "BOOLEAN DEFAULT true",
		"pdv_code":           "TEXT",
		"category_id":        "UUID",
		"subcategory_id":     "UUID",
		"price_promo":        "NUMERIC",
		"volume":             "INTEGER",
		"alcohol_content":    "NUMERIC",
		"vintage":            "TEXT",
		"country":            "TEXT",
		"region":             "TEXT",
		"winery":             "TEXT",
		"wine_type":          "TEXT",
		"grapes":             "TEXT[]",
		"price_bottle":       "NUMERIC",
		"price_half_bottle":  "NUMERIC",
		"price_glass":        "NUMERIC",
	}

	for columnName, columnType := range newColumns {
		// Remove aspas para verificação
		checkColumnName := columnName
		if columnName[0] == '"' {
			checkColumnName = columnName[1 : len(columnName)-1]
		}

		if !migrator.HasColumn("products", checkColumnName) {
			sql := fmt.Sprintf("ALTER TABLE products ADD COLUMN %s %s", columnName, columnType)
			if err := r.db.Exec(sql).Error; err != nil {
				fmt.Printf("⚠️  Aviso ao adicionar coluna '%s': %v\n", columnName, err)
			}
		}
	}

	// Migrar price_normal de forma segura (produtos antigos podem não ter)
	if !migrator.HasColumn("products", "price_normal") {
		// Adicionar como nullable primeiro
		if err := r.db.Exec("ALTER TABLE products ADD COLUMN price_normal NUMERIC").Error; err != nil {
			fmt.Printf("⚠️  Aviso ao adicionar coluna 'price_normal': %v\n", err)
		}
		// Preencher com valor padrão de 0 onde for NULL
		if err := r.db.Exec("UPDATE products SET price_normal = 0 WHERE price_normal IS NULL").Error; err != nil {
			fmt.Printf("⚠️  Aviso ao atualizar coluna 'price_normal': %v\n", err)
		}
		// Tornar NOT NULL
		if err := r.db.Exec("ALTER TABLE products ALTER COLUMN price_normal SET NOT NULL").Error; err != nil {
			fmt.Printf("⚠️  Aviso ao tornar coluna 'price_normal' NOT NULL: %v\n", err)
		}
	}

	fmt.Println("✅ Migração de campos do Product concluída")
}

func NewConnMigrate(db *gorm.DB) IMigrate {
	return &resourceMigrate{db: db}
}
