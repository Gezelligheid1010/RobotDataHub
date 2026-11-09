package episode

type Episode struct {
	ID           string `json:"id"`
	DatasetID    string `json:"dataset_id"`
	EpisodeIndex int    `json:"episode_index"`
	Length       int    `json:"length"`
	FromIndex    int64  `json:"from_index"`
	ToIndex      int64  `json:"to_index"`
	DataFile     string `json:"data_file"`
	VideoFile    string `json:"video_file"`
}
