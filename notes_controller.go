package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type INoteController interface {
	Get(c *gin.Context)
	GetById(c *gin.Context)
	Insert(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type NoteController struct {
	noteService INoteService
}

func getIdFromParamString(idString string) (uint, error) {
	i64, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		return 0, err
	}
	id := uint(i64)
	return id, nil
}

func (nc *NoteController) Get(c *gin.Context) {
	notes := nc.noteService.Get()
	c.IndentedJSON(http.StatusOK, notes)
}

func (nc *NoteController) GetById(c *gin.Context) {
	idString := c.Param("id")
	id, err := getIdFromParamString(idString)
	if err != nil {
		response := ApiResponse{400, "Invalid ID"}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	note, found := nc.noteService.GetById(id)
	if !found {
		response := ApiResponse{404, "Not found"}
		c.IndentedJSON(http.StatusNotFound, response)
		return
	}
	c.IndentedJSON(http.StatusOK, note)
}

func (nc *NoteController) Insert(c *gin.Context) {
	var note Note
	if err := c.BindJSON(&note); err != nil {
		response := ApiResponse{400, "Invalid request body"}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	if _, inserted := nc.noteService.Insert(note); !inserted {
		response := ApiResponse{500, "Failed to insert data"}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}
	response := ApiResponse{200, "Success"}
	c.IndentedJSON(http.StatusOK, response)
}

func (nc *NoteController) Update(c *gin.Context) {
	idString := c.Param("id")
	id, err := getIdFromParamString(idString)
	if err != nil {
		response := ApiResponse{400, "Invalid ID"}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	var note Note
	if err := c.BindJSON(&note); err != nil {
		response := ApiResponse{400, "Invalid request body"}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	if _, inserted := nc.noteService.Update(id, note); !inserted {
		response := ApiResponse{500, "Failed to update data"}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}
	response := ApiResponse{200, "Success"}
	c.IndentedJSON(http.StatusOK, response)
}

func (nc *NoteController) Delete(c *gin.Context) {
	idString := c.Param("id")
	id, err := getIdFromParamString(idString)
	if err != nil {
		response := ApiResponse{400, "Invalid ID"}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	if deleted := nc.noteService.Delete(id); !deleted {
		response := ApiResponse{404, "Not found"}
		c.IndentedJSON(http.StatusNotFound, response)
		return
	}
	response := ApiResponse{200, "Success"}
	c.IndentedJSON(http.StatusOK, response)
}
