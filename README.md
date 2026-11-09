# RobotDataHub

A Go + Vue3 + Python embodied-AI dataset management and validation platform MVP.

## Features

- PostgreSQL-backed Dataset and Episode catalogue.
- Manifest import, episode slice queries and MinIO raw-file upload.
- Vue3 dataset manager, episode table and quality-job dashboard.
- Python/PyArrow local parquet action validator (dimension, NaN/Inf, frame continuity).
- PostgreSQL job claiming with `SKIP LOCKED`, timeout recovery scaffold.

## Quickstart (Windows PowerShell with Docker Desktop)

```powershell
cd deploy
Copy-Item .env.example .env
# Edit API_TOKEN in .env first
docker compose up --build -d
```

Open http://localhost:5173 and paste your `API_TOKEN` from `deploy/.env` in the sidebar, then create a dataset.

API health: http://localhost:8080/healthz

```powershell
# In deploy folder, after replacing <dataset-uuid> with a real ID:
$token = (Select-String -Path .env -Pattern '^API_TOKEN=').Line.Split('=')[1]
Invoke-RestMethod -Uri http://localhost:8080/api/v1/datasets -Headers @{'X-API-Key'=$token}
```

## Import contract

Select a dataset and upload `examples/manifest.example.json` in the import dialog. This only creates episode database metadata. Upload original Parquet and MP4 separately using corresponding relative object paths. To make the worker actually validate a dataset, put its local Parquet at `sample-data/<dataset-uuid>/data.parquet` (mirrored read-only into the worker), containing `episode_index`, `frame_index` and `action` columns, then create a quality job in UI. **Uploading to MinIO alone does not make the file available to the worker in this MVP.**

## Run tests

```bash
cd server && go test ./...
cd ../worker-python && pip install -r requirements.txt pytest && python -m pytest tests
```

## Scope and tradeoffs

This repository is a runnable first engineering baseline, **not** a completed production AI Infra system. Review [Architecture and limitations](docs/ARCHITECTURE.md). In particular, end-to-end MinIO-to-worker processing, native LeRobot import (v2/v3), action/state chart, secure auth and robust job retries remain unfinished.

## Attribution

Original scaffold; no Gin-Vue-Admin source copied. Based on public Gin, Vue, MinIO and PyArrow APIs. Consult dependency licenses before redistribution.


## 目录结构重构

按业务模块组织 Go 与 Vue 代码，见 [docs/REFACTOR.md](docs/REFACTOR.md)。Go 启动入口为 `server/cmd/server/main.go`。
