package urls

import (
	"context"
	"database/sql"
)

type Repo interface {
	Save(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Init(ctx context.Context) error
}

type RepoLayer struct {
	db *sql.DB
}

func NewRepoLayer(db *sql.DB) *RepoLayer {
	return &RepoLayer{
		db: db,
	}
}

func (r *RepoLayer) Save(ctx context.Context, key, value string) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO cache(key, value) VALUES (?, ?)`, key, value)
	return err
}

func (r *RepoLayer) Get(ctx context.Context, key string) (string, error) {
	rowResult := r.db.QueryRowContext(ctx, `SELECT value FROM cache WHERE key=?`, key)
	var value string
	err := rowResult.Scan(&value)
	return value, err
}

func (r *RepoLayer) Init(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS cache (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						key TEXT CHECK (LENGTH(key) > 0) NOT NULL,
						value TEXT CHECK (LENGTH(value) > 0) NOT NULL
						)`)
	return err
}
