package httptransport

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ValidID(c *gin.Context, id string) bool {
	if _, err := uuid.Parse(id); err != nil {
		Fail(c, 400, "invalid id")
		return false
	}
	return true
}
