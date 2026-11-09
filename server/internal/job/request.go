package job

import "encoding/json"

type CreateRequest struct {
	DatasetID string `json:"dataset_id" binding:"required,uuid"`
}
type CompleteRequest struct {
	Result json.RawMessage `json:"result"`
	Error  string          `json:"error"`
}
