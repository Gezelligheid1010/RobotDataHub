package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"robotdatahub/server/internal/dataset"
	"robotdatahub/server/internal/episode"
	"robotdatahub/server/internal/job"
	"robotdatahub/server/internal/platform/database"
	"robotdatahub/server/internal/platform/storage"
)

// env 读取环境变量；未设置时使用默认值。
func env(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}

// requiredEnv 读取必填环境变量。
func requiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("required environment variable %s is missing", key)
	}
	return value
}

func main() {
	// 1. 加载本地开发环境配置。
	// 从 server 目录执行 go run 时，会读取 server/.env。
	// Docker Compose 注入的环境变量优先级更高。
	if err := godotenv.Load(".env"); err != nil {
		log.Println("未找到 .env，使用操作系统环境变量")
	}

	// 2. 读取应用配置。
	databaseURL := requiredEnv("DATABASE_URL")

	minioEndpoint := env("MINIO_ENDPOINT", "localhost:9000")
	minioAccessKey := requiredEnv("MINIO_ACCESS_KEY")
	minioSecretKey := requiredEnv("MINIO_SECRET_KEY")
	minioBucket := env("MINIO_BUCKET", "robotdatahub")

	apiToken := requiredEnv("API_TOKEN")
	corsOrigins := env("CORS_ORIGINS", "http://localhost:5173")

	// 3. 初始化 PostgreSQL。
	initCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	db, err := database.Open(initCtx, databaseURL)
	if err != nil {
		log.Fatalf("PostgreSQL connection failed: %v", err)
	}
	defer db.Close()

	log.Println("PostgreSQL connected successfully")

	// 4. 数据库表结构迁移。
	if err := database.Migrate(initCtx, db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	log.Println("database migration completed")

	// 5. 初始化 MinIO / S3 兼容对象存储。
	objects, err := storage.NewObjectStorage(
		minioEndpoint,
		minioAccessKey,
		minioSecretKey,
		minioBucket,
		false,
	)
	if err != nil {
		log.Fatalf("object storage initialization failed: %v", err)
	}

	// 确保 Bucket 存在。
	if err := objects.EnsureBucket(initCtx); err != nil {
		log.Fatalf("object storage bucket initialization failed: %v", err)
	}

	log.Println("object storage initialized successfully")

	// 6. 初始化业务模块。
	datasets := dataset.NewService(
		dataset.NewRepository(db),
	)

	files := dataset.NewFileService(
		datasets,
		objects,
	)

	episodes := episode.NewService(
		episode.NewRepository(db),
	)

	jobs := job.NewService(
		job.NewRepository(db),
		datasets,
	)

	// 7. 初始化 Gin HTTP 服务。
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(corsOrigins, ","),
		AllowMethods: []string{
			"GET",
			"POST",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Content-Type",
			"X-API-Key",
			"X-Worker-Token",
		},
	}))

	// 注册现有业务路由。
	registerRoutes(
		r,
		dataset.NewHandler(datasets),
		dataset.NewFileHandler(files),
		episode.NewHandler(episodes),
		job.NewHandler(jobs),
		apiToken,
	)

	// 8. 启动任务恢复协程。
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			recoverCtx, recoverCancel := context.WithTimeout(
				context.Background(),
				10*time.Second,
			)

			if err := jobs.Recover(
				recoverCtx,
				5*time.Minute,
			); err != nil {
				log.Printf("job recovery failed: %v", err)
			}

			recoverCancel()
		}
	}()

	// 9. 启动 HTTP 服务。
	log.Println("RobotDataHub server listening on :8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
