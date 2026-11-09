import math

def validate_actions(rows, expected_dim=14):
    """Return structured quality issues for a sequence of episode rows."""
    issues = []
    previous = None
    for row in rows:
        frame = int(row["frame_index"])
        if previous is not None and frame != previous + 1:
            issues.append({"frame_index": frame, "code": "FRAME_GAP", "message": "non-contiguous frame index"})
        previous = frame
        values = row.get("action")
        if values is None or len(values) != expected_dim:
            issues.append({"frame_index": frame, "code": "ACTION_DIM", "message": f"expected {expected_dim} dimensions"})
        elif any(not math.isfinite(float(v)) for v in values):
            issues.append({"frame_index": frame, "code": "NON_FINITE", "message": "action has NaN or infinity"})
    return issues
