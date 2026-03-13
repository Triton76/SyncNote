package model

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrNotOwner     = errors.New("note does not belong to user")
)

type MockNoteStore struct {
	mu          sync.RWMutex
	notes       map[string]*Note
	userNoteIDs map[string][]string
}

func NewMockNoteStore() *MockNoteStore {
	return &MockNoteStore{
		notes:       make(map[string]*Note),
		userNoteIDs: make(map[string][]string),
	}
}

func (m *MockNoteStore) CreateNote(ctx context.Context, note *Note) (*Note, error) {
	_ = ctx
	now := time.Now().UTC()
	created := &Note{
		NoteID:       uuid.NewString(),
		UserID:       note.UserID,
		Title:        note.Title,
		Content:      note.Content,
		Version:      1,
		CreatedAt:    now,
		LastModified: now,
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.notes[created.NoteID] = cloneNote(created)
	m.userNoteIDs[created.UserID] = append(m.userNoteIDs[created.UserID], created.NoteID)

	return cloneNote(created), nil
}

func (m *MockNoteStore) GetNoteByID(ctx context.Context, noteID string) (*Note, error) {
	_ = ctx
	m.mu.RLock()
	defer m.mu.RUnlock()
	note, ok := m.notes[noteID]
	if !ok {
		return nil, ErrNoteNotFound
	}
	return cloneNote(note), nil
}

func (m *MockNoteStore) GetNotesByUserID(ctx context.Context, userID string) ([]*Note, error) {
	_ = ctx
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := m.userNoteIDs[userID]
	result := make([]*Note, 0, len(ids))
	for _, id := range ids {
		note, ok := m.notes[id]
		if !ok {
			continue
		}
		result = append(result, cloneNote(note))
	}
	return result, nil
}

func (m *MockNoteStore) SaveNote(ctx context.Context, noteID string, userID string, content string, expectedVersion int64) (*SaveResult, error) {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()

	note, ok := m.notes[noteID]
	if !ok {
		return &SaveResult{
			Success: false,
			Code:    SaveCodeNotFound,
			Message: "note not found",
		}, nil
	}

	if note.UserID != userID {
		return &SaveResult{
			Success: false,
			Code:    SaveCodeNotFound,
			Message: "note not found",
		}, nil
	}

	if note.Version != expectedVersion {
		return &SaveResult{
			Success:       false,
			Code:          SaveCodeVersionConflict,
			Message:       "version conflict",
			LatestVersion: note.Version,
			LatestContent: note.Content,
		}, nil
	}

	note.Content = content
	note.Version++
	note.LastModified = time.Now().UTC()

	return &SaveResult{
		Success: true,
		Code:    SaveCodeOK,
		Message: "save success",
		Note:    cloneNote(note),
	}, nil
}

func (m *MockNoteStore) DeleteNote(ctx context.Context, noteID string, userID string) error {
	_ = ctx
	m.mu.Lock()
	defer m.mu.Unlock()

	note, ok := m.notes[noteID]
	if !ok {
		return ErrNoteNotFound
	}
	if note.UserID != userID {
		return ErrNotOwner
	}

	delete(m.notes, noteID)
	ids := m.userNoteIDs[userID]
	kept := ids[:0]
	for _, id := range ids {
		if id != noteID {
			kept = append(kept, id)
		}
	}
	m.userNoteIDs[userID] = kept
	return nil
}

func cloneNote(n *Note) *Note {
	if n == nil {
		return nil
	}
	clone := *n
	return &clone
}
