package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=db user=postgres password=postgres dbname=todo port=5432"
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
