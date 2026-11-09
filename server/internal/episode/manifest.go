package episode

import (
	"encoding/json"
	"errors"
	"fmt"
)

type manifest struct {
	FPS       float64 `json:"fps"`
	ActionDim int     `json:"action_dim"`
	Episodes  []struct {
		EpisodeIndex int    `json:"episode_index"`
		Length       int    `json:"length"`
		FromIndex    int64  `json:"from_index"`
		ToIndex      int64  `json:"to_index"`
		DataFile     string `json:"data_file"`
		VideoFile    string `json:"video_file"`
	} `json:"episodes"`
}

func ParseManifest(input []byte) ([]Episode, int, float64, int, error) {
	var v manifest
	if err := json.Unmarshal(input, &v); err != nil {
		return nil, 0, 0, 0, err
	}
	if len(v.Episodes) == 0 {
		return nil, 0, 0, 0, errors.New("episodes cannot be empty")
	}
	seen := map[int]bool{}
	eps := make([]Episode, 0, len(v.Episodes))
	total := 0
	for _, e := range v.Episodes {
		if e.EpisodeIndex < 0 || seen[e.EpisodeIndex] || e.Length <= 0 || e.FromIndex < 0 || e.ToIndex-e.FromIndex != int64(e.Length) {
			return nil, 0, 0, 0, fmt.Errorf("invalid episode %d", e.EpisodeIndex)
		}
		seen[e.EpisodeIndex] = true
		total += e.Length
		eps = append(eps, Episode{EpisodeIndex: e.EpisodeIndex, Length: e.Length, FromIndex: e.FromIndex, ToIndex: e.ToIndex, DataFile: e.DataFile, VideoFile: e.VideoFile})
	}
	if v.FPS <= 0 || v.ActionDim <= 0 {
		return nil, 0, 0, 0, errors.New("fps and action_dim must be positive")
	}
	return eps, total, v.FPS, v.ActionDim, nil
}
