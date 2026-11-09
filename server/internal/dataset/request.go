package dataset

type CreateRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=120"`
	Description string `json:"description"`
	Source      string `json:"source"`
}
