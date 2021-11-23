package main

const (
	UNSPECIFIED_ID uint = 0
)

type INoteService interface {
	Get() []Note
	GetById(id uint) (Note, bool)
	Create(note Note) (uint, error)
	Update(id uint, note Note) (uint, error)
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

func (ns *NoteService) Create(note Note) (uint, error) {
	if note.ID != UNSPECIFIED_ID {
		return UNSPECIFIED_ID, &IllegalIdError{}
	}

	id, ok := ns.noteRepository.Create(note)
	if !ok {
		return UNSPECIFIED_ID, &InternalError{}
	}
	return id, nil
}

func (ns *NoteService) Update(id uint, note Note) (uint, error) {
	if note.ID != UNSPECIFIED_ID {
		return UNSPECIFIED_ID, &IllegalIdError{}
	}

	id, ok := ns.noteRepository.Update(id, note)
	if !ok {
		return UNSPECIFIED_ID, &InternalError{}
	}
	return id, nil
}

func (ns *NoteService) Delete(id uint) bool {
	ok := ns.noteRepository.Delete(id)
	return ok
}
