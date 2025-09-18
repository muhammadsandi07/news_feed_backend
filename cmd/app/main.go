package main

import (
	"log"
	"news-feed/internal/domain/user"
	"news-feed/internal/middleware"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// db connection
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("failed to migrate database", err)
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	app := fiber.New(middleware.NewFiberConfig())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	api := app.Group("/api")
	api.Post("/register", userHandler.Register)
	api.Post("/login", userHandler.Login)

	port := ":3000"
	log.Println("🚀 Server running at http://localhost" + port)
	log.Fatal(app.Listen(port))

}
