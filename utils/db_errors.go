package utils

import (
	"errors"

	"github.com/lib/pq"
)

// PostgreSQL error codes
const (
	PgUniqueViolation     = "23505" // unique_violation
	PgForeignKeyViolation = "23503" // foreign_key_violation
)

// IsDuplicateKeyError verifica se o erro é uma violação de unique constraint
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == PgUniqueViolation
	}
	return false
}
