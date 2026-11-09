package episode

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }
func (r *Repository) ListEpisodes(ctx context.Context, id string) ([]Episode, error) {
	rows, err := r.db.Query(ctx, `SELECT id,dataset_id,episode_index,length,from_index,to_index,data_file,video_file FROM episodes WHERE dataset_id=$1 ORDER BY episode_index`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Episode, 0)
	for rows.Next() {
		var e Episode
		if err = rows.Scan(&e.ID, &e.DatasetID, &e.EpisodeIndex, &e.Length, &e.FromIndex, &e.ToIndex, &e.DataFile, &e.VideoFile); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}
func (r *Repository) ImportEpisodes(ctx context.Context, datasetID string, eps []Episode, totalFrames int, fps float64, actionDim int) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var exists string
	if err = tx.QueryRow(ctx, `SELECT id FROM datasets WHERE id=$1 FOR UPDATE`, datasetID).Scan(&exists); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `DELETE FROM episodes WHERE dataset_id=$1`, datasetID); err != nil {
		return err
	}
	for _, e := range eps {
		if _, err = tx.Exec(ctx, `INSERT INTO episodes(id,dataset_id,episode_index,length,from_index,to_index,data_file,video_file) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`, uuid.NewString(), datasetID, e.EpisodeIndex, e.Length, e.FromIndex, e.ToIndex, e.DataFile, e.VideoFile); err != nil {
			return fmt.Errorf("episode %d: %w", e.EpisodeIndex, err)
		}
	}
	if _, err = tx.Exec(ctx, `UPDATE datasets SET status='ready',total_episodes=$2,total_frames=$3,fps=$4,action_dim=$5 WHERE id=$1`, datasetID, len(eps), totalFrames, fps, actionDim); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
