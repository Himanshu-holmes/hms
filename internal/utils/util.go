package util

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetUserIDFromString(subject string) (pgtype.UUID, error) {
	var id pgtype.UUID
	err := id.Scan(subject)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("failed to parse UUID from string: %w", err)
	}
	return id, nil
}