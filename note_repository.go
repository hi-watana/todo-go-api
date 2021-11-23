package main

import (
	"errors"

	"gorm.io/gorm"
)

type INoteRepository interface {
	GetAll() []Note
	GetById(id uint64) (Note, bool)
	Create(note Note) (uint64, bool)
	Update(id uint64, note Note) (uint64, bool)
	Delete(id uint64) bool
}

type NoteRepository struct {
	db *gorm.DB
}

func (nr *NoteRepository) GetAll() []Note {
	var notes []Note
	nr.db.Find(&notes)
	return notes
}

func (nr *NoteRepository) GetById(id uint64) (Note, bool) {
	var note Note
	result := nr.db.First(&note, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return note, false
	}
	return note, true
}

func (nr *NoteRepository) Create(note Note) (uint64, bool) {
	result := nr.db.Create(&note)
	if result.Error != nil {
		return 0, false
	}
	return note.ID, true
}

func (nr *NoteRepository) Update(id uint64, note Note) (uint64, bool) {
	result := nr.db.Model(&Note{ID: id}).Updates(note)
	if result.Error != nil {
		return 0, false
	}
	return id, true
}

func (nr *NoteRepository) Delete(id uint64) bool {
	if result := nr.db.Delete(&Note{}, id); result.Error != nil {
		return false
	}
	return true
}
