package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432", host, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Note{})

	noteRepository := &NoteRepository{db}
	noteService := &NoteService{noteRepository}
	noteController := NoteController{noteService}

	router := gin.Default()

	router.GET("/notes", noteController.Get)
	router.GET("/notes/:id", noteController.GetById)
	router.POST("/notes", noteController.Create)
	router.PUT("/notes/:id", noteController.Update)
	router.DELETE("/notes/:id", noteController.Delete)

	router.Run(":8080")
}
