package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"campaign-platform/internal/handler"
	"campaign-platform/internal/model"
	"campaign-platform/internal/repository"
	"campaign-platform/internal/service"
)

func main() {
	dsn := getEnv("DB_DSN", "host=localhost user=app password=app dbname=campaign_platform port=5432 sslmode=disable")
	port := getEnv("PORT", "8080")
	env := getEnv("APP_ENV", "development")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 开发环境自动迁移
	if env == "development" {
		db.AutoMigrate(
			&model.CampaignTemplate{},
			&model.Campaign{},
			&model.Component{},
			&model.CampaignVersion{},
		)
	}

	// 依赖注入
	campaignRepo := repository.NewCampaignRepo(db)
	componentRepo := repository.NewComponentRepo(db)

	campaignSvc := service.NewCampaignService(campaignRepo)
	componentSvc := service.NewComponentService(componentRepo)

	campaignH := handler.NewCampaignHandler(campaignSvc)
	componentH := handler.NewComponentHandler(componentSvc)
	renderH := handler.NewRenderHandler(campaignSvc, componentSvc)

	// 路由
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/templates", campaignH.ListTemplates)
		v1.POST("/templates", campaignH.CreateTemplate)
		v1.GET("/campaigns", campaignH.ListCampaigns)
		v1.GET("/campaigns/:id", campaignH.GetCampaign)
		v1.POST("/campaigns", campaignH.CreateCampaign)
		v1.PUT("/campaigns/:id", campaignH.UpdateCampaign)
		v1.PATCH("/campaigns/:id/status", campaignH.UpdateStatus)
		v1.GET("/campaigns/:id/versions", campaignH.ListVersions)
		v1.GET("/components", componentH.ListComponents)
		v1.GET("/render/:slug", renderH.RenderConfig)
	}

	log.Printf("starting :%s env=%s", port, env)
	r.Run(":" + port)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
