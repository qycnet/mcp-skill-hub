package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/qycnet/mcp-skill-hub/internal/api"
	"github.com/qycnet/mcp-skill-hub/internal/auth"
	"github.com/qycnet/mcp-skill-hub/internal/skill"
	"github.com/qycnet/mcp-skill-hub/internal/storage"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	version   = "0.1.0-dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// 初始化配置
	initConfig()

	// 设置 Gin 模式
	if viper.GetString("mode") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	db := initDatabase()

	// 初始化存储
	objectStorage := initStorage()

	// 初始化服务
	jwtSecret := viper.GetString("jwt_secret")
	if jwtSecret == "" {
		log.Fatal("❌ JWT_SECRET 环境变量必须设置")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("❌ JWT_SECRET 必须至少 32 个字符")
	}
	
	skillService := skill.NewService(db, objectStorage)
	authService := auth.NewService(db, jwtSecret)

	// 创建路由
	router := setupRouter(skillService, authService)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", viper.GetInt("port")),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 优雅关闭
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Println("正在关闭服务器...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("服务器关闭失败：%v", err)
		}
	}()

	// 启动服务器
	log.Printf("🚀 MCP Skill Hub 服务器启动 (v%s)", version)
	log.Printf("📡 监听端口：%d", viper.GetInt("port"))
	log.Printf("📊 运行模式：%s", viper.GetString("mode"))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("服务器启动失败：%v", err)
	}

	log.Println("服务器已关闭")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// 默认配置
	viper.SetDefault("port", 8080)
	viper.SetDefault("mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.name", "mcp_skill_hub")
	viper.SetDefault("database.user", "postgres")
	// 注意：数据库密码必须通过环境变量设置，不提供默认值
	viper.SetDefault("storage.endpoint", "localhost:9000")
	viper.SetDefault("storage.access_key", "minioadmin")
	// 注意：存储密钥必须通过环境变量设置
	viper.SetDefault("storage.bucket", "mcp-skills")
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.db", 0)

	// 支持从环境变量读取
	viper.AutomaticEnv()

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("⚠️  未找到配置文件，使用默认配置：%v", err)
	}

	log.Printf("📄 使用配置文件：%s", viper.ConfigFileUsed())
}

func initDatabase() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		viper.GetString("database.host"),
		viper.GetInt("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.name"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(
		&skill.Skill{},
		&skill.SkillVersion{},
		&auth.User{},
		&auth.APIKey{},
	); err != nil {
		log.Fatalf("数据库迁移失败：%v", err)
	}

	log.Println("✅ 数据库初始化完成")
	return db
}

func initStorage() storage.ObjectStorage {
	config := storage.Config{
		Endpoint:  viper.GetString("storage.endpoint"),
		AccessKey: viper.GetString("storage.access_key"),
		SecretKey: viper.GetString("storage.secret_key"),
		Bucket:    viper.GetString("storage.bucket"),
		UseSSL:    viper.GetBool("storage.use_ssl"),
	}

	storage, err := storage.NewMinIOStorage(config)
	if err != nil {
		log.Fatalf("对象存储初始化失败：%v", err)
	}

	log.Println("✅ 对象存储初始化完成")
	return storage
}

func setupRouter(skillService *skill.Service, authService *auth.Service) *gin.Engine {
	router := gin.Default()

	// CORS 配置（生产环境应限制为具体域名）
	allowedOrigins := viper.GetStringSlice("cors.allowed_origins")
	if len(allowedOrigins) == 0 {
		// 开发环境允许 localhost
		if viper.GetString("mode") == "debug" {
			allowedOrigins = []string{"http://localhost:*", "http://127.0.0.1:*"}
		} else {
			// 生产环境必须配置具体域名
			log.Fatal("❌ 生产环境必须配置 cors.allowed_origins")
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "healthy",
			"version":     version,
			"build_time":  buildTime,
			"git_commit":  gitCommit,
			"timestamp":   time.Now().UTC(),
		})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		// 公开 API
		public := v1.Group("")
		{
			public.GET("/skills", api.ListSkills(skillService))
			public.GET("/skills/:id", api.GetSkill(skillService))
			public.GET("/skills/:id/download", api.DownloadSkill(skillService))
			public.GET("/search", api.SearchSkills(skillService))
			public.GET("/categories", api.ListCategories(skillService))
		}

		// 需要认证的 API
		protected := v1.Group("")
		protected.Use(auth.AuthMiddleware(authService))
		{
			protected.POST("/skills", api.PublishSkill(skillService))
			protected.PUT("/skills/:id", api.UpdateSkill(skillService))
			protected.DELETE("/skills/:id", api.DeleteSkill(skillService))
			protected.POST("/skills/:id/rate", api.RateSkill(skillService))
			
			// 用户管理
			protected.GET("/user/profile", api.GetUserProfile(authService))
			protected.PUT("/user/profile", api.UpdateUserProfile(authService))
			protected.POST("/user/api-keys", api.CreateAPIKey(authService))
			protected.DELETE("/user/api-keys/:keyId", api.RevokeAPIKey(authService))
		}

		// 管理员 API
		admin := v1.Group("/admin")
		admin.Use(auth.AuthMiddleware(authService), auth.AdminMiddleware())
		{
			admin.GET("/skills", api.AdminListSkills(skillService))
			admin.PUT("/skills/:id/approve", api.ApproveSkill(skillService))
			admin.PUT("/skills/:id/reject", api.RejectSkill(skillService))
			admin.GET("/users", api.AdminListUsers(authService))
			admin.GET("/analytics", api.GetAnalytics(skillService))
		}
	}

	// 认证端点
	authGroup := v1.Group("/auth")
	{
		authGroup.POST("/register", api.Register(authService))
		authGroup.POST("/login", api.Login(authService))
		authGroup.POST("/refresh", api.RefreshToken(authService))
		authGroup.POST("/logout", api.Logout(authService))
	}

	return router
}
