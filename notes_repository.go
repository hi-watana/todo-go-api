package main

import (
	"errors"

	"gorm.io/gorm"
)

type INoteRepository interface {
	GetAll() []Note
	GetById(id uint) (Note, bool)
	Insert(note Note) (uint, bool)
	Update(id uint, note Note) (uint, bool)
	Delete(id uint) bool
}

type NoteRepository struct {
	db *gorm.DB
}

func (nr *NoteRepository) GetAll() []Note {
	var notes []Note
	nr.db.Find(&notes)
	return notes
}

func (nr *NoteRepository) GetById(id uint) (Note, bool) {
	var note Note
	result := nr.db.First(&note, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return note, false
	}
	return note, true
}

func (nr *NoteRepository) Insert(note Note) (uint, bool) {
	result := nr.db.Create(&note)
	if result.Error != nil {
		return 0, false
	}
	return note.ID, true
}

func (nr *NoteRepository) Update(id uint, note Note) (uint, bool) {
	newNote := CopyNote(&note)
	newNote.ID = id
	result := nr.db.Save(&newNote)
	if result.Error != nil {
		return 0, false
	}
	return id, true
}

func (nr *NoteRepository) Delete(id uint) bool {
	result := nr.db.Delete(&Note{}, id)
	return result.RowsAffected > 0
}
