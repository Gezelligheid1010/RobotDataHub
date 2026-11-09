package database

import (
	"context"
	"embed"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

//go:embed schema.sql
var schema embed.FS

func Migrate(ctx context.Context, db *pgxpool.Pool) error {
	data, err := schema.ReadFile("schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, strings.TrimSpace(string(data)))
	return err
}
