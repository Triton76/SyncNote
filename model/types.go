package model

import (
	"context"
	"time"
)

type SaveCode string

const (
	SaveCodeOK              SaveCode = "OK"
	SaveCodeVersionConflict SaveCode = "VERSION_CONFLICT"
	SaveCodeNotFound        SaveCode = "NOT_FOUND"
	SaveCodeInvalidParam    SaveCode = "INVALID_PARAM"
)

type SaveResult struct {
	Success       bool
	Code          SaveCode
	Message       string
	LatestVersion int64
	LatestContent string
	Note          *Note
}

type NoteStore interface {
	CreateNote(ctx context.Context, note *Note) (*Note, error)
	GetNoteByID(ctx context.Context, noteID string) (*Note, error)
	GetNotesByUserID(ctx context.Context, userID string) ([]*Note, error)
	SaveNote(ctx context.Context, noteID string, userID string, content string, expectedVersion int64) (*SaveResult, error)
	DeleteNote(ctx context.Context, noteID string, userID string) error
}

type Note struct {
	NoteID       string    `db:"id"`
	UserID       string    `db:"user_id"`
	Title        string    `db:"title"`
	Content      string    `db:"content"`
	Version      int64     `db:"version"`
	CreatedAt    time.Time `db:"created_at"`
	LastModified time.Time `db:"updated_at"`
}
