package persistence

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/tzDel/orchestrAIgent/internal/domain"
)

func TestSQLiteSessionRepository_Save_PersistsSessionSuccessfully(t *testing.T) {
	// arrange
	repository, cleanup := setupTestRepository(t)
	defer cleanup()

	sessionID, _ := domain.NewSessionID("test-session")
	session, _ := domain.NewSession(sessionID, "/path/to/worktree")
	ctx := context.Background()

	// act
	err := repository.Save(ctx, session)

	// assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	retrieved, err := repository.FindByID(ctx, sessionID)
	if err != nil {
		t.Fatalf("expected to find saved session, got error: %v", err)
	}

	assertSessionEquals(t, session, retrieved)
}

func TestSQLiteSessionRepository_Save_UpdatesExistingSession(t *testing.T) {
	// arrange
	repository, cleanup := setupTestRepository(t)
	defer cleanup()

	sessionID, _ := domain.NewSessionID("test-session")
	session, _ := domain.NewSession(sessionID, "/path/to/worktree")
	ctx := context.Background()

	repository.Save(ctx, session)

	// act
	session.MarkReviewed()
	err := repository.Save(ctx, session)

	// assert
	if err != nil {
		t.Fatalf("expected no error when updating, got: %v", err)
	}

	retrieved, _ := repository.FindByID(ctx, sessionID)
	if retrieved.Status() != domain.StatusReviewed {
		t.Errorf("expected status %s, got %s", domain.StatusReviewed, retrieved.Status())
	}
}

func TestSQLiteSessionRepository_FindByID_ReturnsErrorWhenNotFound(t *testing.T) {
	// arrange
	repository, cleanup := setupTestRepository(t)
	defer cleanup()

	sessionID, _ := domain.NewSessionID("nonexistent")
	ctx := context.Background()

	// act
	_, err := repository.FindByID(ctx, sessionID)

	// assert
	if err == nil {
		t.Error("expected error for nonexistent session, got nil")
	}
}

func TestSQLiteSessionRepository_FindAll_ReturnsAllSessions(t *testing.T) {
	// arrange
	repository, cleanup := setupTestRepository(t)
	defer cleanup()

	sessionID1, _ := domain.NewSessionID("session-01")
	sessionID2, _ := domain.NewSessionID("session-02")
	session1, _ := domain.NewSession(sessionID1, "/path/1")
	session2, _ := domain.NewSession(sessionID2, "/path/2")
	ctx := context.Background()

	repository.Save(ctx, session1)
	repository.Save(ctx, session2)

	// act
	sessions, err := repository.FindAll(ctx)

	// assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(sessions) != 2 {
		t.Errorf("expected 2 sessions, got %d", len(sessions))
	}
}

func TestSQLiteSessionRepository_FindAll_ReturnsEmptySliceWhenNoSessions(t *testing.T) {
	// arrange
	repository, cleanup := setupTestRepository(t)
	defer cleanup()
	ctx := context.Background()

	// act
	sessions, err := repository.FindAll(ctx)

	// assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(sessions) != 0 {
		t.Errorf("expected empty slice, got %d sessions", len(sessions))
	}
}

func TestSQLiteSessionRepository_Exists_ReturnsTrueWhenSessionExists(t *testing.T) {
	// arrange
	repository, cleanup := setupTestRepository(t)
	defer cleanup()

	sessionID, _ := domain.NewSessionID("test-session")
	session, _ := domain.NewSession(sessionID, "/path/to/worktree")
	ctx := context.Background()

	repository.Save(ctx, session)

	// act
	exists, err := repository.Exists(ctx, sessionID)

	// assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !exists {
		t.Error("expected session to exist, but it does not")
	}
}

func TestSQLiteSessionRepository_Exists_ReturnsFalseWhenSessionDoesNotExist(t *testing.T) {
	// arrange
	repository, cleanup := setupTestRepository(t)
	defer cleanup()

	sessionID, _ := domain.NewSessionID("nonexistent")
	ctx := context.Background()

	// act
	exists, err := repository.Exists(ctx, sessionID)

	// assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if exists {
		t.Error("expected session not to exist, but it does")
	}
}

func TestSQLiteSessionRepository_Delete_RemovesSessionSuccessfully(t *testing.T) {
	// arrange
	repository, cleanup := setupTestRepository(t)
	defer cleanup()

	sessionID, _ := domain.NewSessionID("test-session")
	session, _ := domain.NewSession(sessionID, "/path/to/worktree")
	ctx := context.Background()

	repository.Save(ctx, session)

	// act
	err := repository.Delete(ctx, sessionID)

	// assert
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	exists, _ := repository.Exists(ctx, sessionID)
	if exists {
		t.Error("expected session to be deleted, but it still exists")
	}
}

func TestSQLiteSessionRepository_Delete_ReturnsErrorWhenSessionNotFound(t *testing.T) {
	// arrange
	repository, cleanup := setupTestRepository(t)
	defer cleanup()

	sessionID, _ := domain.NewSessionID("nonexistent")
	ctx := context.Background()

	// act
	err := repository.Delete(ctx, sessionID)

	// assert
	if err == nil {
		t.Error("expected error when deleting nonexistent session, got nil")
	}
}

// Helper functions

func setupTestRepository(t *testing.T) (*SQLiteSessionRepository, func()) {
	tempDir, err := os.MkdirTemp("", "sqlite-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	dbPath := filepath.Join(tempDir, "test.db")
	repository, err := NewSQLiteSessionRepository(dbPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("failed to create repository: %v", err)
	}

	cleanup := func() {
		repository.Close()
		os.RemoveAll(tempDir)
	}

	return repository, cleanup
}

func assertSessionEquals(t *testing.T, expected, actual *domain.Session) {
	if expected.ID().String() != actual.ID().String() {
		t.Errorf("expected ID %s, got %s", expected.ID().String(), actual.ID().String())
	}

	if expected.Status() != actual.Status() {
		t.Errorf("expected status %s, got %s", expected.Status(), actual.Status())
	}

	if expected.WorktreePath() != actual.WorktreePath() {
		t.Errorf("expected worktree path %s, got %s", expected.WorktreePath(), actual.WorktreePath())
	}

	if expected.BranchName() != actual.BranchName() {
		t.Errorf("expected branch name %s, got %s", expected.BranchName(), actual.BranchName())
	}
}
