package persistence

import (
	"context"
	"testing"

	"github.com/tzDel/orchestrAIgent/internal/domain"
)

func TestInMemoryRepository_SaveAndFind(t *testing.T) {
	// arrange
	repository := NewInMemorySessionRepository()
	ctx := context.Background()
	sessionID, _ := domain.NewSessionID("test-session")
	session, _ := domain.NewSession(sessionID, "/path/to/worktree")

	// act
	err := repository.Save(ctx, session)
	if err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	found, err := repository.FindByID(ctx, sessionID)

	// assert
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}

	if found.ID().String() != session.ID().String() {
		t.Errorf("FindByID() returned wrong session")
	}
}

func TestInMemoryRepository_FindByID_NotFound(t *testing.T) {
	// arrange
	repository := NewInMemorySessionRepository()
	ctx := context.Background()
	sessionID, _ := domain.NewSessionID("nonexistent")

	// act
	_, err := repository.FindByID(ctx, sessionID)

	// assert
	if err == nil {
		t.Error("FindByID() expected error for non-existent session")
	}
}

func TestInMemoryRepository_Exists(t *testing.T) {
	// arrange
	repository := NewInMemorySessionRepository()
	ctx := context.Background()
	sessionID, _ := domain.NewSessionID("test-session")

	// act
	exists, err := repository.Exists(ctx, sessionID)

	// assert
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if exists {
		t.Error("Exists() returned true for non-existent session")
	}

	// arrange
	session, _ := domain.NewSession(sessionID, "/path")
	repository.Save(ctx, session)

	// act
	exists, err = repository.Exists(ctx, sessionID)

	// assert
	if err != nil {
		t.Fatalf("Exists() error: %v", err)
	}
	if !exists {
		t.Error("Exists() returned false for existing session")
	}
}
