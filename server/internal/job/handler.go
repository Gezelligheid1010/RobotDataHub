package job

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	httptransport "robotdatahub/server/internal/transport/http"
)

type Handler struct{ service *Service }

func NewHandler(s *Service) *Handler { return &Handler{service: s} }
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httptransport.Fail(c, 400, err.Error())
		return
	}
	v, err := h.service.Create(c.Request.Context(), req.DatasetID)
	if errors.Is(err, pgx.ErrNoRows) {
		httptransport.Fail(c, 404, "dataset not found")
		return
	}
	if err != nil {
		httptransport.Fail(c, 500, err.Error())
		return
	}
	c.JSON(201, gin.H{"data": v})
}
func (h *Handler) List(c *gin.Context) {
	v, err := h.service.List(c.Request.Context())
	if err != nil {
		httptransport.Fail(c, 500, err.Error())
		return
	}
	httptransport.OK(c, v)
}
func (h *Handler) Claim(c *gin.Context) {
	v, err := h.service.Claim(c.Request.Context())
	if errors.Is(err, pgx.ErrNoRows) {
		c.Status(204)
		return
	}
	if err != nil {
		httptransport.Fail(c, 500, err.Error())
		return
	}
	httptransport.OK(c, v)
}
func (h *Handler) Complete(c *gin.Context) {
	id := c.Param("id")
	if !httptransport.ValidID(c, id) {
		return
	}
	var req CompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httptransport.Fail(c, 400, err.Error())
		return
	}
	if err := h.service.Complete(c.Request.Context(), id, req.Result, req.Error); err != nil {
		httptransport.Fail(c, 500, err.Error())
		return
	}
	httptransport.OK(c, gin.H{"saved": true})
}
