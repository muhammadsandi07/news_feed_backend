package main

import (
	"fmt"
	"log"
	"os"

	"news-feed/internal/config"
	"news-feed/internal/domain/follow"
	"news-feed/internal/domain/post"
	"news-feed/internal/domain/user"
	"news-feed/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println(" No .env file found, using system env")
	}
	cfg := config.LoadConfig()

	// Init DB
	db := initDB()

	// Run migrations
	if err := db.AutoMigrate(&user.User{}, &post.Post{}, &follow.Follow{}); err != nil {
		log.Fatal(" Failed to migrate database:", err)
	}

	// Init repositories & services
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, cfg.JWTSecret, cfg.AccessExpireMin, cfg.RefreshExpireHr)
	userHandler := user.NewHandler(userService)

	followRepo := follow.NewRepository(db)
	followService := follow.NewService(followRepo)
	followHandler := follow.NewHandler(followService)

	postRepo := post.NewRepository(db)
	postService := post.NewService(postRepo)
	postHandler := post.NewHandler(postService, followService)

	// Setup Fiber
	app := setupFiber()

	// Public routes
	api := app.Group("/api")
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)
	api.Post("/refresh", userHandler.Refresh)

	// Protected routes
	protected := api.Group("", middleware.JWTProtected())
	protected.Post("/posts", postHandler.CreatePost)
	protected.Get("/feed", postHandler.GetFeed)
	protected.Post("/follow/:id", followHandler.Follow)
	protected.Delete("/follow/:id", followHandler.Unfollow)

	// Run server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":3000"
	}
	log.Printf(" Server running at http://localhost:%s\n", port)
	log.Fatal(app.Listen(":" + port))
}

func initDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(" Failed to connect database:", err)
	}
	return db
}

func setupFiber() *fiber.App {
	app := fiber.New(middleware.NewFiberConfig())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("FRONTEND_URL"),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))
	return app
}
