package dataset

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
	value, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		httptransport.Fail(c, 409, err.Error())
		return
	}
	c.JSON(201, gin.H{"data": value})
}
func (h *Handler) List(c *gin.Context) {
	v, err := h.service.List(c.Request.Context())
	if err != nil {
		httptransport.Fail(c, 500, err.Error())
		return
	}
	httptransport.OK(c, v)
}
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	if !httptransport.ValidID(c, id) {
		return
	}
	v, err := h.service.Get(c.Request.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		httptransport.Fail(c, 404, "not found")
		return
	}
	if err != nil {
		httptransport.Fail(c, 500, err.Error())
		return
	}
	httptransport.OK(c, v)
}
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if !httptransport.ValidID(c, id) {
		return
	}
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		httptransport.Fail(c, 404, err.Error())
		return
	}
	httptransport.OK(c, gin.H{"deleted": true})
}
