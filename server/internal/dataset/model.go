package dataset

import "time"

type Dataset struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Source        string    `json:"source"`
	Status        string    `json:"status"`
	TotalEpisodes int       `json:"total_episodes"`
	TotalFrames   int       `json:"total_frames"`
	FPS           float64   `json:"fps"`
	ActionDim     int       `json:"action_dim"`
	CreatedAt     time.Time `json:"created_at"`
}
