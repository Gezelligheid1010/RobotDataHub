package dataset

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	httptransport "robotdatahub/server/internal/transport/http"
)

type FileHandler struct{ service *FileService }

func NewFileHandler(s *FileService) *FileHandler { return &FileHandler{service: s} }
func (h *FileHandler) Upload(c *gin.Context) {
	id := c.Param("id")
	if !httptransport.ValidID(c, id) {
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		httptransport.Fail(c, 400, err.Error())
		return
	}
	if file.Size > 1<<30 {
		httptransport.Fail(c, 413, "file exceeds 1 GiB")
		return
	}
	name := c.PostForm("path")
	if name == "" {
		name = file.Filename
	}
	src, err := file.Open()
	if err != nil {
		httptransport.Fail(c, 500, err.Error())
		return
	}
	defer src.Close()
	key, err := h.service.Upload(c.Request.Context(), id, name, src, file.Size)
	if errors.Is(err, pgx.ErrNoRows) {
		httptransport.Fail(c, 404, "dataset not found")
		return
	}
	if err != nil {
		httptransport.Fail(c, 400, err.Error())
		return
	}
	httptransport.OK(c, gin.H{"key": key})
}
func (h *FileHandler) URL(c *gin.Context) {
	id := c.Param("id")
	if !httptransport.ValidID(c, id) {
		return
	}
	url, err := h.service.URL(c.Request.Context(), id, c.Query("path"))
	if errors.Is(err, pgx.ErrNoRows) {
		httptransport.Fail(c, 404, "dataset not found")
		return
	}
	if err != nil {
		httptransport.Fail(c, 400, err.Error())
		return
	}
	httptransport.OK(c, gin.H{"url": url})
}
