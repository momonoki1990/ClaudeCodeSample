package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"todo-api/handler"
	"todo-api/repository"
	"todo-api/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func initDB() *gorm.DB {
	user := getEnv("DB_USER", "root")
	password := getEnv("DB_PASSWORD", "password")
	host := getEnv("DB_HOST", "db:3306")
	dbName := getEnv("DB_NAME", "todos")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4", user, password, host, dbName)

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
		log.Printf("waiting for db (%d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("could not connect to db:", err)
	}

	log.Println("db ready")
	return db
}

func main() {
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
