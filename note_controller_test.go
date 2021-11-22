package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (ms *MockService) Get() []Note {
	ret := ms.Called()
	return ret.Get(0).([]Note)
}

func (ms *MockService) GetById(id uint) (Note, bool) {
	ret := ms.Called(id)
	return ret.Get(0).(Note), ret.Get(1).(bool)
}

func (ms *MockService) Create(note Note) (uint, error) {
	ret := ms.Called(note)
	return ret.Get(0).(uint), ret.Error(1)
}

func (ms *MockService) Update(id uint, note Note) (uint, error) {
	ret := ms.Called(id, note)
	return ret.Get(0).(uint), ret.Error(1)
}

func (ms *MockService) Delete(id uint) bool {
	ret := ms.Called(id)
	return ret.Get(0).(bool)
}

func TestNoteController_Get(t *testing.T) {
	mockService := &MockService{}
	noteController := NoteController{mockService}
	response := httptest.NewRecorder()
	ginContext, _ := gin.CreateTestContext(response)

	mockService.On("Get").Return([]Note{{1, "test_title", "test_content"}})

	req, _ := http.NewRequest("GET", "/notes", nil)
	ginContext.Request = req

	noteController.Get(ginContext)

	assert.Equal(t, http.StatusOK, response.Code)
	expected, _ := json.MarshalIndent(&[]Note{{1, "test_title", "test_content"}}, "", "    ")
	assert.Equal(t, expected, response.Body.Bytes())
}

func TestNoteController_GetById(t *testing.T) {
	for _, td := range []struct {
		title                  string
		inputId                uint
		inputPathParameter     string
		outputNote             Note
		outputOk               bool
		expectedStatus         int
		expectedResponseObject interface{}
	}{
		{
			title:              "Returns note if found",
			inputId:            1,
			inputPathParameter: "1",
			outputNote: Note{
				ID:      1,
				Title:   "test_title",
				Content: "test_content",
			},
			outputOk:       true,
			expectedStatus: http.StatusOK,
			expectedResponseObject: &Note{
				ID:      1,
				Title:   "test_title",
				Content: "test_content",
			},
		},
		{
			title:              "Returns \"Invalid ID\" message if invalid ID was specified",
			inputPathParameter: "xxx",
			expectedStatus:     http.StatusBadRequest,
			expectedResponseObject: &ApiResponse{
				Status:  400,
				Message: "Invalid ID",
			},
		},
		{
			title:              "Returns \"Not found\" message if not found",
			inputId:            2,
			inputPathParameter: "2",
			outputNote:         Note{},
			outputOk:           false,
			expectedStatus:     http.StatusNotFound,
			expectedResponseObject: &ApiResponse{
				Status:  404,
				Message: "Not found",
			},
		},
	} {
		t.Run("GetById: "+td.title, func(t *testing.T) {
			var (
				mockService    = &MockService{}
				noteController = NoteController{mockService}
				response       = httptest.NewRecorder()
				ginContext, _  = gin.CreateTestContext(response)
				req, _         = http.NewRequest("GET", "/notes/"+td.inputPathParameter, nil)
			)

			mockService.On("GetById", td.inputId).Return(td.outputNote, td.outputOk)

			ginContext.Request = req
			ginContext.Params = append(ginContext.Params, gin.Param{Key: "id", Value: td.inputPathParameter})

			noteController.GetById(ginContext)

			assert.Equal(t, td.expectedStatus, response.Code)
			expected, _ := json.MarshalIndent(td.expectedResponseObject, "", "    ")
			assert.Equal(t, expected, response.Body.Bytes())
		})
	}
}

func noteToBytes(note Note) []byte {
	ret, _ := json.Marshal(note)
	return ret
}

func TestNoteController_Create(t *testing.T) {
	for _, td := range []struct {
		title                  string
		requestBody            []byte
		inputNote              Note
		outputError            error
		expectedStatus         int
		expectedResponseObject interface{}
	}{
		{
			title: "Returns success message",
			requestBody: noteToBytes(Note{
				Title:   "test_title",
				Content: "test_content",
			}),
			inputNote: Note{
				Title:   "test_title",
				Content: "test_content",
			},
			outputError:    nil,
			expectedStatus: http.StatusOK,
			expectedResponseObject: &ApiResponse{
				Status:  200,
				Message: "Success",
			},
		},
		{
			title:          "Returns \"Invalid request body\" message",
			requestBody:    []byte("not json"),
			expectedStatus: http.StatusBadRequest,
			expectedResponseObject: &ApiResponse{
				Status:  400,
				Message: "Invalid request body",
			},
		},
		{
			title: "Returns \"ID must not be specified\" message",
			requestBody: noteToBytes(Note{
				ID:      1,
				Title:   "test_title",
				Content: "test_content",
			}),
			inputNote: Note{
				ID:      1,
				Title:   "test_title",
				Content: "test_content",
			},
			outputError:    &IllegalIdError{},
			expectedStatus: http.StatusBadRequest,
			expectedResponseObject: &ApiResponse{
				Status:  400,
				Message: "ID must not be specified",
			},
		},
		{
			title: "Returns \"Unexpected error\" message",
			requestBody: noteToBytes(Note{
				Title:   "test_title",
				Content: "test_content",
			}),
			inputNote: Note{
				Title:   "test_title",
				Content: "test_content",
			},
			outputError:    &InternalError{},
			expectedStatus: http.StatusInternalServerError,
			expectedResponseObject: &ApiResponse{
				Status:  500,
				Message: "Unexpected error",
			},
		},
	} {
		t.Run("Create: "+td.title, func(t *testing.T) {
			mockService := &MockService{}
			noteController := NoteController{mockService}
			response := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(response)

			var id uint = 1
			mockService.On("Create", td.inputNote).Return(id, td.outputError)

			req, _ := http.NewRequest("POST", "/notes/", bytes.NewReader(td.requestBody))
			ginContext.Request = req

			noteController.Create(ginContext)

			assert.Equal(t, td.expectedStatus, response.Code)
			expected, _ := json.MarshalIndent(td.expectedResponseObject, "", "    ")
			assert.Equal(t, expected, response.Body.Bytes())
		})
	}
}

func TestNoteController_Update(t *testing.T) {
	for _, td := range []struct {
		title                  string
		inputId                uint
		inputPathParameter     string
		requestBody            []byte
		inputNote              Note
		found                  bool
		outputError            error
		expectedStatus         int
		expectedResponseObject interface{}
	}{
		{
			title:              "Returns success message",
			inputId:            1,
			inputPathParameter: "1",
			requestBody: noteToBytes(Note{
				Title:   "test_title",
				Content: "test_content",
			}),
			inputNote: Note{
				Title:   "test_title",
				Content: "test_content",
			},
			found:          true,
			outputError:    nil,
			expectedStatus: http.StatusOK,
			expectedResponseObject: &ApiResponse{
				Status:  200,
				Message: "Success",
			},
		},
		{
			title:              "Returns \"Not found\" message",
			inputId:            1,
			inputPathParameter: "1",
			requestBody: noteToBytes(Note{
				Title:   "test_title",
				Content: "test_content",
			}),
			inputNote: Note{
				Title:   "test_title",
				Content: "test_content",
			},
			found:          false,
			expectedStatus: http.StatusNotFound,
			expectedResponseObject: &ApiResponse{
				Status:  404,
				Message: "Not found",
			},
		},
		{
			title:              "Returns \"Invalid ID\" message",
			inputPathParameter: "xxx",
			expectedStatus:     http.StatusBadRequest,
			expectedResponseObject: &ApiResponse{
				Status:  400,
				Message: "Invalid ID",
			},
		},
		{
			title:              "Returns \"Invalid request\" message",
			inputId:            1,
			inputPathParameter: "1",
			requestBody:        []byte("not json"),
			expectedStatus:     http.StatusBadRequest,
			expectedResponseObject: &ApiResponse{
				Status:  400,
				Message: "Invalid request body",
			},
		},
		{
			title:              "Returns \"Illegal ID in request body\" message",
			inputId:            1,
			inputPathParameter: "1",
			requestBody: noteToBytes(Note{
				ID:      2,
				Title:   "test_title",
				Content: "test_content",
			}),
			inputNote: Note{
				ID:      2,
				Title:   "test_title",
				Content: "test_content",
			},
			found:          true,
			outputError:    &IllegalIdError{},
			expectedStatus: http.StatusBadRequest,
			expectedResponseObject: &ApiResponse{
				Status:  400,
				Message: "Illegal ID in request body",
			},
		},
		{
			title:              "Returns \"Unexpected error\" message",
			inputId:            1,
			inputPathParameter: "1",
			requestBody: noteToBytes(Note{
				Title:   "test_title",
				Content: "test_content",
			}),
			inputNote: Note{
				Title:   "test_title",
				Content: "test_content",
			},
			found:          true,
			outputError:    &InternalError{},
			expectedStatus: http.StatusInternalServerError,
			expectedResponseObject: &ApiResponse{
				Status:  500,
				Message: "Unexpected error",
			},
		},
	} {
		t.Run("Update: "+td.title, func(t *testing.T) {
			mockService := &MockService{}
			noteController := NoteController{mockService}
			response := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(response)

			mockService.On("GetById", td.inputId).Return(Note{}, td.found)
			mockService.On("Update", td.inputId, td.inputNote).Return(td.inputId, td.outputError)

			req, _ := http.NewRequest("PUT", "/notes/"+td.inputPathParameter, bytes.NewReader(td.requestBody))
			ginContext.Request = req
			ginContext.Params = append(ginContext.Params, gin.Param{Key: "id", Value: td.inputPathParameter})

			noteController.Update(ginContext)

			assert.Equal(t, td.expectedStatus, response.Code)
			expected, _ := json.MarshalIndent(td.expectedResponseObject, "", "    ")
			assert.Equal(t, expected, response.Body.Bytes())
		})
	}
}

func TestNoteController_Delete(t *testing.T) {
	for _, td := range []struct {
		title                  string
		inputId                uint
		inputPathParameter     string
		found                  bool
		outputOk               bool
		expectedStatus         int
		expectedResponseObject interface{}
	}{
		{
			title:              "Returns success message",
			inputId:            1,
			inputPathParameter: "1",
			found:              true,
			outputOk:           true,
			expectedStatus:     http.StatusOK,
			expectedResponseObject: &ApiResponse{
				Status:  200,
				Message: "Success",
			},
		},
		{
			title:              "Returns \"Invalid ID\" message",
			inputPathParameter: "xxx",
			expectedStatus:     http.StatusBadRequest,
			expectedResponseObject: &ApiResponse{
				Status:  400,
				Message: "Invalid ID",
			},
		},
		{
			title:              "Returns \"Not found\" message",
			inputId:            2,
			inputPathParameter: "2",
			found:              false,
			expectedStatus:     http.StatusNotFound,
			expectedResponseObject: &ApiResponse{
				Status:  404,
				Message: "Not found",
			},
		},
		{
			title:              "Returns \"Unexpected error\" message",
			inputId:            1,
			inputPathParameter: "1",
			found:              true,
			outputOk:           false,
			expectedStatus:     http.StatusInternalServerError,
			expectedResponseObject: &ApiResponse{
				Status:  500,
				Message: "Unexpected error",
			},
		},
	} {
		t.Run("Delete: "+td.title, func(t *testing.T) {
			mockService := &MockService{}
			noteController := NoteController{mockService}
			response := httptest.NewRecorder()
			ginContext, _ := gin.CreateTestContext(response)

			mockService.On("GetById", td.inputId).Return(Note{}, td.found)
			mockService.On("Delete", td.inputId).Return(td.outputOk)

			req, _ := http.NewRequest("DELETE", "/notes/"+td.inputPathParameter, nil)
			ginContext.Request = req
			ginContext.Params = append(ginContext.Params, gin.Param{Key: "id", Value: td.inputPathParameter})

			noteController.Delete(ginContext)

			assert.Equal(t, td.expectedStatus, response.Code)
			expected, _ := json.MarshalIndent(td.expectedResponseObject, "", "    ")
			assert.Equal(t, expected, response.Body.Bytes())
		})
	}
}
