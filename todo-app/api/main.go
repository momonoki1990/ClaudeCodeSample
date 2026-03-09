package main

import (
	"database/sql"
	"log"
	"os"
	"time"
	"todo-api/handler"
	"todo-api/repository"
	"todo-api/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func initDB() *sql.DB {
	dsn := getEnv("DB_USER", "root") + ":" +
		getEnv("DB_PASSWORD", "password") + "@tcp(" +
		getEnv("DB_HOST", "db:3306") + ")/" +
		getEnv("DB_NAME", "todos") + "?parseTime=true"

	var db *sql.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id   INT AUTO_INCREMENT PRIMARY KEY,
		text TEXT NOT NULL,
		done BOOLEAN NOT NULL DEFAULT FALSE
	)`)
	if err != nil {
		log.Fatal("could not create table:", err)
	}
	log.Println("db ready")
	return db
}

func main() {
	db := initDB()

	repo := repository.NewTodoRepository(db)
	svc := service.NewTodoService(repo)
	h := handler.NewTodoHandler(svc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api := e.Group("/api")
	api.GET("/todos", h.GetAll)
	api.POST("/todos", h.Create)
	api.PUT("/todos/:id", h.Update)
	api.DELETE("/todos/:id", h.Delete)

	e.Logger.Fatal(e.Start(":8080"))
}
