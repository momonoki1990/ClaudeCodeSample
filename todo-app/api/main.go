package main

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
	_ "time/tzdata"
	"todo-api/handler"
	"todo-api/repository"
	"todo-api/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type healthResponse struct {
	Status string `json:"status"`
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func initLogger() {
	var level slog.Level
	switch getEnv("LOG_LEVEL", "INFO") {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})))
}

func initDB() *gorm.DB {
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	host := getEnv("DB_HOST", "db:3306")
	dbName := getEnv("DB_NAME", "todos")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&loc=Asia%%2FTokyo", user, password, host, dbName)

	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			err = sqlDB.Ping()
		}
		if err == nil {
			break
		}
		slog.Warn("waiting for db", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		slog.Error("could not connect to db", "error", err)
		os.Exit(1)
	}

	slog.Info("db ready")
	return db
}

func main() {
	initLogger()

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		slog.Error("failed to load timezone", "error", err)
		os.Exit(1)
	}
	time.Local = loc

	now := time.Now()
	slog.Info("server starting",
		"time", now.Format(time.RFC3339),
		"timezone", now.Location().String(),
		"tz_offset", now.Format("-07:00"),
		"go", runtime.Version(),
		"os", runtime.GOOS+"/"+runtime.GOARCH,
	)

	db := initDB()

	todoRepo := repository.NewTodoRepository(db)
	todoSvc := service.NewTodoService(todoRepo)
	todoHandler := handler.NewTodoHandler(todoSvc)

	categoryRepo := repository.NewCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categorySvc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, healthResponse{Status: "ok"})
	})

	api := e.Group("/api")
	api.GET("/todos", todoHandler.GetAll)
	api.POST("/todos", todoHandler.Create)
	api.PUT("/todos/reorder", todoHandler.Reorder)
	api.PUT("/todos/:id", todoHandler.Update)
	api.DELETE("/todos/done", todoHandler.DeleteDone)
	api.DELETE("/todos/:id", todoHandler.Delete)

	api.GET("/categories", categoryHandler.GetAll)
	api.POST("/categories", categoryHandler.Create)
	api.PUT("/categories/reorder", categoryHandler.Reorder)
	api.PUT("/categories/:id", categoryHandler.Update)
	api.DELETE("/categories/:id", categoryHandler.Delete)

	e.Logger.Fatal(e.Start(":8080"))
}
