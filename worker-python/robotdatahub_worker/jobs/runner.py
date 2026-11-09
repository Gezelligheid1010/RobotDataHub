import os
from collections import defaultdict
from robotdatahub_worker.parsers.lerobot import load_action_rows
from robotdatahub_worker.validators.actions import validate_actions

def execute(dataset_id):
    """Validate parquet mirrored under DATA_ROOT/<dataset-id>/data.parquet.

    This deliberately supports only an explicit local file in the MVP; MinIO
    and LeRobot v3 multi-file index resolution are follow-up work.
    """
    path = os.path.join(os.environ.get("DATA_ROOT", "/data"), dataset_id, "data.parquet")
    if not os.path.isfile(path):
        raise FileNotFoundError(f"parquet not mounted: {path}")
    rows = load_action_rows(path)
    grouped = defaultdict(list)
    for row in rows:
        grouped[int(row["episode_index"])].append(row)
    issues = []
    for idx, batch in sorted(grouped.items()):
        batch.sort(key=lambda row: row["frame_index"])
        for issue in validate_actions(batch, int(os.environ.get("ACTION_DIM", "14"))):
            issues.append({"episode_index": idx, **issue})
    return {"dataset_id": dataset_id, "episodes_checked": len(grouped), "frames_checked": len(rows), "issues": issues, "passed": not issues}
