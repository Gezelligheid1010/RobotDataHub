CREATE TABLE IF NOT EXISTS datasets (
 id UUID PRIMARY KEY, name TEXT NOT NULL UNIQUE, description TEXT NOT NULL DEFAULT '', source TEXT NOT NULL DEFAULT '',
 status TEXT NOT NULL DEFAULT 'created', total_episodes INT NOT NULL DEFAULT 0,
 total_frames INT NOT NULL DEFAULT 0, fps DOUBLE PRECISION NOT NULL DEFAULT 0,
 action_dim INT NOT NULL DEFAULT 0, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS episodes (
 id UUID PRIMARY KEY, dataset_id UUID NOT NULL REFERENCES datasets(id) ON DELETE CASCADE,
 episode_index INT NOT NULL, length INT NOT NULL CHECK(length >= 0),
 from_index BIGINT NOT NULL, to_index BIGINT NOT NULL,
 data_file TEXT NOT NULL DEFAULT '', video_file TEXT NOT NULL DEFAULT '',
 UNIQUE(dataset_id, episode_index)
);
CREATE TABLE IF NOT EXISTS jobs (
 id UUID PRIMARY KEY, dataset_id UUID NOT NULL REFERENCES datasets(id) ON DELETE CASCADE,
 state TEXT NOT NULL DEFAULT 'pending', attempts INT NOT NULL DEFAULT 0,
 max_attempts INT NOT NULL DEFAULT 3, error TEXT NOT NULL DEFAULT '',
 result JSONB, created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), locked_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_jobs_claim ON jobs(state,created_at);
