package main

import (
	"github.com/gin-gonic/gin"
	"robotdatahub/server/internal/dataset"
	"robotdatahub/server/internal/episode"
	"robotdatahub/server/internal/job"
)

func registerRoutes(r *gin.Engine, d *dataset.Handler, files *dataset.FileHandler, e *episode.Handler, j *job.Handler, token string) {
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	v := r.Group("/api/v1")
	v.Use(func(c *gin.Context) {
		if c.GetHeader("X-API-Key") != token {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	})
	v.GET("/datasets", d.List)
	v.POST("/datasets", d.Create)
	v.GET("/datasets/:id", d.Get)
	v.DELETE("/datasets/:id", d.Delete)
	v.GET("/datasets/:id/episodes", e.List)
	v.POST("/datasets/:id/import-manifest", e.ImportManifest)
	v.POST("/datasets/:id/files", files.Upload)
	v.GET("/datasets/:id/file-url", files.URL)
	v.GET("/jobs", j.List)
	v.POST("/jobs", j.Create)
	worker := r.Group("/internal/v1")
	worker.Use(func(c *gin.Context) {
		if c.GetHeader("X-Worker-Token") != token {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	})
	worker.POST("/jobs/claim", j.Claim)
	worker.POST("/jobs/:id/complete", j.Complete)
}
