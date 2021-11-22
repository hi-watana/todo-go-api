package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type INoteController interface {
	Get(c *gin.Context)
	GetById(c *gin.Context)
	Create(c *gin.Context)
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

func (nc *NoteController) Create(c *gin.Context) {
	var note Note
	if err := c.BindJSON(&note); err != nil {
		response := ApiResponse{400, "Invalid request body"}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	_, err := nc.noteService.Create(note)
	if errors.Is(err, &IllegalIdError{}) {
		response := ApiResponse{400, "ID must not be specified"}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	if errors.Is(err, &InternalError{}) {
		response := ApiResponse{500, "Unexpected error"}
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

	if _, found := nc.noteService.GetById(id); !found {
		response := ApiResponse{404, "Not found"}
		c.IndentedJSON(http.StatusNotFound, response)
		return
	}
	_, err = nc.noteService.Update(id, note)
	if errors.Is(err, &IllegalIdError{}) {
		response := ApiResponse{400, "Illegal ID in request body"}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}
	if errors.Is(err, &InternalError{}) {
		response := ApiResponse{500, "Unexpected error"}
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
	if _, found := nc.noteService.GetById(id); !found {
		response := ApiResponse{404, "Not found"}
		c.IndentedJSON(http.StatusNotFound, response)
		return
	}
	if deleted := nc.noteService.Delete(id); !deleted {
		response := ApiResponse{500, "Unexpected error"}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}
	response := ApiResponse{200, "Success"}
	c.IndentedJSON(http.StatusOK, response)
}
