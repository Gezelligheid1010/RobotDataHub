package httptransport

import "github.com/gin-gonic/gin"

func OK(c *gin.Context, value any)                    { c.JSON(200, gin.H{"data": value}) }
func Fail(c *gin.Context, status int, message string) { c.JSON(status, gin.H{"error": message}) }
