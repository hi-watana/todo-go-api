package main

type Note struct {
	ID uint `gorm:"primaryKey"`
	Title string
	Content string
}

func CopyNote(note *Note) Note {
	newNote := Note{note.ID, note.Title, note.Content}
	return newNote
}