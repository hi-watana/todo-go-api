package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mr *MockRepository) GetAll() []Note {
	ret := mr.Called()
	return ret.Get(0).([]Note)
}

func (mr *MockRepository) GetById(id uint) (Note, bool) {
	ret := mr.Called(id)
	return ret.Get(0).(Note), ret.Get(1).(bool)
}

func (mr *MockRepository) Create(note Note) (uint, bool) {
	ret := mr.Called(note)
	return ret.Get(0).(uint), ret.Get(1).(bool)
}

func (mr *MockRepository) Update(id uint, note Note) (uint, bool) {
	ret := mr.Called(id, note)
	return ret.Get(0).(uint), ret.Get(1).(bool)
}

func (mr *MockRepository) Delete(id uint) bool {
	ret := mr.Called(id)
	return ret.Get(0).(bool)
}

func TestNoteService_Get(t *testing.T) {
	mockRepository := &MockRepository{}
	noteService := NoteService{mockRepository}

	notes := []Note{
		{
			ID:      1,
			Title:   "test_title",
			Content: "test_content",
		},
	}

	mockRepository.On("GetAll").Return(notes)

	actualNotes := noteService.Get()
	assert.Equal(t, notes, actualNotes)
}

func TestNoteService_GetById(t *testing.T) {
	for _, td := range []struct {
		title string
		inputId uint
		outputNote Note
		outputOk bool
	} {
		{
			title: "Return note and true if note was found",
			inputId: 1,
			outputNote: Note{
				ID:      1,
				Title:   "test_title",
				Content: "test_content",
			},
			outputOk: true,
		},
		{
			title: "Return empty note and false if note was not found",
			inputId: 1,
			outputNote: Note{},
			outputOk: false,
		},
	} {
		t.Run("GetById: " + td.title, func(t *testing.T) {
			mockRepository := &MockRepository{}
			noteService := NoteService{mockRepository}

			mockRepository.On("GetById", td.inputId).Return(td.outputNote, td.outputOk)

			actualNote, found := noteService.GetById(td.inputId)
			assert.Equal(t, td.outputOk, found)
			assert.Equal(t, td.outputNote, actualNote)
		})
	}
}

func TestNoteService_Create(t *testing.T) {
	for _, td := range []struct {
		title string
		inputNote Note
		okFromRepository bool
		outputId uint
		outputError error
	} {
		{
			title: "Returns ID and true if successfully inserted",
			inputNote: Note{
				Title: "test_title",
				Content: "test_content",
			},
			okFromRepository: true,
			outputId: 1,
			outputError: nil,
		},
		{
			title: "Returns 0 and IllegalIdError if ID is specified in request body",
			inputNote: Note{
				ID: 1,
				Title: "test_title",
				Content: "test_content",
			},
			outputId: UNSPECIFIED_ID,
			outputError: &IllegalIdError{},
		},
		{
			title: "Returns 0 and InternalError if DB operation failed",
			inputNote: Note{
				Title: "test_title",
				Content: "test_content",
			},
			okFromRepository: false,
			outputId: UNSPECIFIED_ID,
			outputError: &InternalError{},
		},
	} {
		t.Run("Create: " + td.title, func(t *testing.T) {
			mockRepository := &MockRepository{}
			noteService := NoteService{mockRepository}

			mockRepository.On("Create", td.inputNote).Return(td.outputId, td.okFromRepository)

			_, err := noteService.Create(td.inputNote)
			assert.IsType(t, td.outputError, err)
		})
	}
}

func TestNoteService_Update(t *testing.T) {
	for _, td := range []struct {
		title string
		inputId uint
		inputNote Note
		okFromRepository bool
		outputId uint
		outputError error
	} {
		{
			title: "Returns ID and nil if successfully updated",
			inputId: 1,
			inputNote: Note{
				Title: "test_title",
				Content: "test_content",
			},
			okFromRepository: true,
			outputId: 1,
			outputError: nil,
		},
		{
			title: "Returns 0 and IllegalIdError if ID in request body and that in path is not the same",
			inputId: 1,
			inputNote: Note{
				ID: 2,
				Title: "test_title",
				Content: "test_content",
			},
			outputId: UNSPECIFIED_ID,
			outputError: &IllegalIdError{},
		},
		{
			title: "Returns 0 and InternalError if update failed",
			inputId: 1,
			inputNote: Note{
				Title: "test_title",
				Content: "test_content",
			},
			okFromRepository: false,
			outputId: UNSPECIFIED_ID,
			outputError: &InternalError{},
		},
	} {
		t.Run("Update: " + td.title, func(t *testing.T) {
			mockRepository := &MockRepository{}
			noteService := NoteService{mockRepository}

			mockRepository.On("Update", td.inputId, td.inputNote).Return(td.outputId, td.okFromRepository)

			_, err := noteService.Update(td.inputId, td.inputNote)
			assert.IsType(t, td.outputError, err)
		})
	}
}

func TestNoteService_Delete(t *testing.T) {
	for _, td := range []struct {
		title string
		inputId uint
		outputOk bool
	} {
		{
			title: "Return true if successfully deleted",
			inputId: 1,
			outputOk: true,
		},
		{
			title: "Return false if deletion failed",
			inputId: 1,
			outputOk: false,
		},
	} {
		t.Run("GetById: " + td.title, func(t *testing.T) {
			mockRepository := &MockRepository{}
			noteService := NoteService{mockRepository}

			mockRepository.On("Delete", td.inputId).Return(td.outputOk)

			ok := noteService.Delete(td.inputId)
			assert.Equal(t, td.outputOk, ok)
		})
	}
}