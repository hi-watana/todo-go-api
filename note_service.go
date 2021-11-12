package main

type INoteService interface {
	Get() []Note
	GetById(id uint) (Note, bool)
	Insert(note Note) (uint, bool)
	Update(id uint, note Note) (uint, bool)
	Delete(id uint) bool
}

type NoteService struct {
	noteRepository INoteRepository
}

func (ns *NoteService) Get() []Note {
	notes := ns.noteRepository.GetAll()
	return notes
}

func (ns *NoteService) GetById(id uint) (Note, bool) {
	note, found := ns.noteRepository.GetById(id)
	return note, found
}

func (ns *NoteService) Insert(note Note) (uint, bool) {
	id, ok := ns.noteRepository.Insert(note)
	return id, ok
}

func (ns *NoteService) Update(id uint, note Note) (uint, bool) {
	id, ok := ns.noteRepository.Update(id, note)
	return id, ok
}

func (ns *NoteService) Delete(id uint) bool {
	ok := ns.noteRepository.Delete(id)
	return ok
}
