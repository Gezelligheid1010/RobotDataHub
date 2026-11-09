package dataset

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }
func (r *Repository) CreateDataset(ctx context.Context, name, description, source string) (Dataset, error) {
	d := Dataset{ID: uuid.NewString(), Name: name, Description: description, Source: source, Status: "created"}
	err := r.db.QueryRow(ctx, `INSERT INTO datasets(id,name,description,source) VALUES($1,$2,$3,$4) RETURNING created_at`, d.ID, d.Name, d.Description, d.Source).Scan(&d.CreatedAt)
	return d, err
}

func (r *Repository) ListDatasets(ctx context.Context) ([]Dataset, error) {
	rows, err := r.db.Query(ctx, `SELECT id,name,description,source,status,total_episodes,total_frames,fps,action_dim,created_at FROM datasets ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Dataset, 0)
	for rows.Next() {
		var d Dataset
		if err = rows.Scan(&d.ID, &d.Name, &d.Description, &d.Source, &d.Status, &d.TotalEpisodes, &d.TotalFrames, &d.FPS, &d.ActionDim, &d.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (r *Repository) GetDataset(ctx context.Context, id string) (Dataset, error) {
	var d Dataset
	err := r.db.QueryRow(ctx, `SELECT id,name,description,source,status,total_episodes,total_frames,fps,action_dim,created_at FROM datasets WHERE id=$1`, id).Scan(&d.ID, &d.Name, &d.Description, &d.Source, &d.Status, &d.TotalEpisodes, &d.TotalFrames, &d.FPS, &d.ActionDim, &d.CreatedAt)
	return d, err
}

func (r *Repository) DeleteDataset(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM datasets WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("dataset not found")
	}
	return nil
}
