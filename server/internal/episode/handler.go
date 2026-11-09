package episode

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"io"
	httptransport "robotdatahub/server/internal/transport/http"
)

type Handler struct{ service *Service }

func NewHandler(s *Service) *Handler { return &Handler{service: s} }
func (h *Handler) List(c *gin.Context) {
	id := c.Param("id")
	if !httptransport.ValidID(c, id) {
		return
	}
	v, err := h.service.List(c.Request.Context(), id)
	if err != nil {
		httptransport.Fail(c, 500, err.Error())
		return
	}
	httptransport.OK(c, v)
}
func (h *Handler) ImportManifest(c *gin.Context) {
	id := c.Param("id")
	if !httptransport.ValidID(c, id) {
		return
	}
	data, err := io.ReadAll(io.LimitReader(c.Request.Body, 4<<20+1))
	if err != nil {
		httptransport.Fail(c, 400, err.Error())
		return
	}
	if len(data) > 4<<20 {
		httptransport.Fail(c, 413, "manifest too large")
		return
	}
	count, total, err := h.service.Import(c.Request.Context(), id, data)
	if errors.Is(err, pgx.ErrNoRows) {
		httptransport.Fail(c, 404, "dataset not found")
		return
	}
	if err != nil {
		httptransport.Fail(c, 400, err.Error())
		return
	}
	httptransport.OK(c, gin.H{"episodes": count, "frames": total})
}
