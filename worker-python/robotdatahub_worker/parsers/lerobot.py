"""Parse a locally available LeRobot parquet file (supported simple MVP layout)."""
import pyarrow.parquet as pq

def load_action_rows(parquet_path):
    table = pq.read_table(parquet_path, columns=["episode_index", "frame_index", "action"])
    return table.to_pylist()
