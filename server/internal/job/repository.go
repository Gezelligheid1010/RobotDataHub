package job

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Repository struct{ db *pgxpool.Pool }

func NewRepository(db *pgxpool.Pool) *Repository { return &Repository{db: db} }
func (r *Repository) CreateJob(ctx context.Context, datasetID string) (Job, error) {
	j := Job{ID: uuid.NewString(), DatasetID: datasetID, State: "pending", MaxAttempts: 3}
	err := r.db.QueryRow(ctx, `INSERT INTO jobs(id,dataset_id) VALUES($1,$2) RETURNING created_at,updated_at`, j.ID, j.DatasetID).Scan(&j.CreatedAt, &j.UpdatedAt)
	return j, err
}

func (r *Repository) ListJobs(ctx context.Context) ([]Job, error) {
	rows, err := r.db.Query(ctx, `SELECT id,dataset_id,state,attempts,max_attempts,error,result,created_at,updated_at FROM jobs ORDER BY created_at DESC LIMIT 100`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]Job, 0)
	for rows.Next() {
		var j Job
		if err = rows.Scan(&j.ID, &j.DatasetID, &j.State, &j.Attempts, &j.MaxAttempts, &j.Error, &j.Result, &j.CreatedAt, &j.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, j)
	}
	return out, rows.Err()
}

func (r *Repository) ClaimJob(ctx context.Context) (Job, error) {
	var j Job
	// SKIP LOCKED prevents multiple worker instances from claiming the same job.
	err := r.db.QueryRow(ctx, `WITH candidate AS (SELECT id FROM jobs WHERE state='pending' ORDER BY created_at FOR UPDATE SKIP LOCKED LIMIT 1) UPDATE jobs SET state='running',attempts=attempts+1,locked_at=now(),updated_at=now() WHERE id=(SELECT id FROM candidate) RETURNING id,dataset_id,state,attempts,max_attempts,error,created_at,updated_at`).Scan(&j.ID, &j.DatasetID, &j.State, &j.Attempts, &j.MaxAttempts, &j.Error, &j.CreatedAt, &j.UpdatedAt)
	return j, err
}

func (r *Repository) CompleteJob(ctx context.Context, id string, result json.RawMessage, jobError string) error {
	state := "succeeded"
	if jobError != "" {
		state = "failed"
	}
	tag, err := r.db.Exec(ctx, `UPDATE jobs SET state=$2,result=$3,error=$4,locked_at=NULL,updated_at=now() WHERE id=$1 AND state='running'`, id, state, result, jobError)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.New("job is not running or not found")
	}
	return nil
}

func (r *Repository) RecoverJobs(ctx context.Context, olderThan time.Duration) error {
	_, err := r.db.Exec(ctx, `UPDATE jobs SET state=CASE WHEN attempts<max_attempts THEN 'pending' ELSE 'failed' END, error='worker heartbeat timeout',locked_at=NULL,updated_at=now() WHERE state='running' AND locked_at < $1`, time.Now().Add(-olderThan))
	return err
}
