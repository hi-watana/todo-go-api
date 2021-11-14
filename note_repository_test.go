package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type NoteRepositoryTestSuite struct {
	suite.Suite
	noteRepository NoteRepository
	mock           sqlmock.Sqlmock
}

func (ts *NoteRepositoryTestSuite) SetupTest() {
	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	ts.mock = mock
	noteRepository := NoteRepository{}
	noteRepository.db, _ = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	ts.noteRepository = noteRepository
}

func (ts *NoteRepositoryTestSuite) TearDownTest() {
	db, _ := ts.noteRepository.db.DB()
	db.Close()
}

func (ts *NoteRepositoryTestSuite) TestNoteRepository_GetAll() {
	var (
		id       uint = 1
		title         = "title"
		content       = "content"
		id2      uint = 2
		title2        = "title2"
		content2      = "content3"
	)

	rows := sqlmock.NewRows([]string{"id", "title", "content"})
	rows = rows.AddRow(id, title, content)
	rows = rows.AddRow(id2, title2, content2)
	query := ts.noteRepository.db.Session(&gorm.Session{DryRun: true}).Find(&[]Note{}).Statement.SQL.String()
	ts.mock.ExpectQuery(query).WillReturnRows(rows)

	notes := ts.noteRepository.GetAll()

	assert.Equal(ts.T(), 2, len(notes))
	assert.Equal(ts.T(), id, notes[0].ID)
	assert.Equal(ts.T(), title, notes[0].Title)
	assert.Equal(ts.T(), content, notes[0].Content)
	assert.Equal(ts.T(), id2, notes[1].ID)
	assert.Equal(ts.T(), title2, notes[1].Title)
	assert.Equal(ts.T(), content2, notes[1].Content)
}

func (ts *NoteRepositoryTestSuite) TestNoteRepository_GetById() {
	for _, td := range []struct {
		title        string
		inputId      uint
		outputRows   *sqlmock.Rows
		expectedNote Note
		expectedOk   bool
	}{
		{
			title:      "Returns note and true if found",
			inputId:    2,
			outputRows: sqlmock.NewRows([]string{"id", "title", "content"}).AddRow(2, "test_title2", "test_content2"),
			expectedNote: Note{
				ID:      2,
				Title:   "test_title2",
				Content: "test_content2",
			},
			expectedOk: true,
		},
		{
			title:      "Returns empty note and false if not found",
			inputId:    3,
			outputRows: sqlmock.NewRows([]string{"id", "title", "content"}),
			expectedOk: false,
		},
	} {
		ts.Run("GetById: "+td.title, func() {
			var note Note
			query := ts.noteRepository.db.Session(&gorm.Session{DryRun: true}).First(&note, td.inputId).Statement.SQL.String()
			ts.mock.ExpectQuery(query).WillReturnRows(td.outputRows)

			actualNote, actualOk := ts.noteRepository.GetById(td.inputId)

			assert.Equal(ts.T(), td.expectedOk, actualOk)
			assert.Equal(ts.T(), td.expectedNote, actualNote)
		})
	}
}

func (ts *NoteRepositoryTestSuite) TestNoteRepository_Create_success() {
	var (
		id   uint = 1
		note      = Note{
			Title:   "test_title",
			Content: "test_content",
		}
		rows = sqlmock.NewRows([]string{"id"}).AddRow(id)
	)
	ts.mock.ExpectBegin()
	ts.mock.ExpectCommit()
	query := ts.noteRepository.db.Session(&gorm.Session{DryRun: true}).Create(&note).Statement.SQL.String()
	ts.mock.ExpectBegin()
	ts.mock.ExpectQuery(query).WillReturnRows(rows)
	ts.mock.ExpectCommit()

	actualId, actualOk := ts.noteRepository.Create(note)

	assert.Equal(ts.T(), true, actualOk)
	assert.Equal(ts.T(), id, actualId)
}

func (ts *NoteRepositoryTestSuite) TestNoteRepository_Create_failed() {
	var (
		note = Note{
			Title:   "test_title",
			Content: "test_content",
		}
	)
	ts.mock.ExpectBegin()
	ts.mock.ExpectCommit()
	query := ts.noteRepository.db.Session(&gorm.Session{DryRun: true}).Create(&note).Statement.SQL.String()
	ts.mock.ExpectBegin()
	ts.mock.ExpectQuery(query).WillReturnError(gorm.ErrInvalidDB) // Anything error will do.
	ts.mock.ExpectRollback()

	actualId, actualOk := ts.noteRepository.Create(note)

	assert.Equal(ts.T(), false, actualOk)
	assert.Equal(ts.T(), UNSPECIFIED_ID, actualId)
}

func TestNoteRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(NoteRepositoryTestSuite))
}
