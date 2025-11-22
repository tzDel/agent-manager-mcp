package domain

import "testing"

func TestNewSession(t *testing.T) {
	// arrange
	sessionID, _ := NewSessionID("test-session")
	worktreePath := "/path/to/worktree"

	// act
	session, err := NewSession(sessionID, worktreePath)

	// assert
	if err != nil {
		t.Fatalf("NewSession() unexpected error: %v", err)
	}

	if session.ID().String() != "test-session" {
		t.Errorf("ID() = %q, want %q", session.ID().String(), "test-session")
	}

	if session.Status() != StatusCreated {
		t.Errorf("Status() = %q, want %q", session.Status(), StatusCreated)
	}

	if session.WorktreePath() != worktreePath {
		t.Errorf("WorktreePath() = %q, want %q", session.WorktreePath(), worktreePath)
	}

	if session.BranchName() != "session-test-session" {
		t.Errorf("BranchName() = %q, want %q", session.BranchName(), "session-test-session")
	}
}

func TestNewSession_InvalidWorktreePath(t *testing.T) {
	// arrange
	sessionID, _ := NewSessionID("test-session")

	// act
	_, err := NewSession(sessionID, "")

	// assert
	if err == nil {
		t.Error("NewSession() with empty path expected error, got nil")
	}
}

func TestSession_MarkMerged(t *testing.T) {
	// arrange
	sessionID, _ := NewSessionID("test-session")
	session, _ := NewSession(sessionID, "/path")

	// act
	session.MarkMerged()

	// assert
	if session.Status() != StatusMerged {
		t.Errorf("Status after MarkMerged() = %q, want %q", session.Status(), StatusMerged)
	}
}

func TestSession_MarkFailed(t *testing.T) {
	// arrange
	sessionID, _ := NewSessionID("test-session")
	session, _ := NewSession(sessionID, "/path")

	// act
	session.MarkFailed()

	// assert
	if session.Status() != StatusFailed {
		t.Errorf("Status after MarkFailed() = %q, want %q", session.Status(), StatusFailed)
	}
}

func TestSession_MarkRemoved(t *testing.T) {
	// arrange
	sessionID, _ := NewSessionID("test-session")
	session, _ := NewSession(sessionID, "/path")

	// act
	session.MarkRemoved()

	// assert
	if session.Status() != StatusRemoved {
		t.Errorf("Status after MarkRemoved() = %q, want %q", session.Status(), StatusRemoved)
	}
}
