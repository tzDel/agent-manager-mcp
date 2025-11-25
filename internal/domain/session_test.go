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

	if session.Status() != StatusOpen {
		t.Errorf("Status() = %q, want %q", session.Status(), StatusOpen)
	}

	if session.WorktreePath() != worktreePath {
		t.Errorf("WorktreePath() = %q, want %q", session.WorktreePath(), worktreePath)
	}

	if session.BranchName() != "orchestragent-test-session" {
		t.Errorf("BranchName() = %q, want %q", session.BranchName(), "orchestragent-test-session")
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

func TestSession_MarkReviewed(t *testing.T) {
	// arrange
	sessionID, _ := NewSessionID("test-session")
	session, _ := NewSession(sessionID, "/path")

	// act
	session.MarkReviewed()

	// assert
	if session.Status() != StatusReviewed {
		t.Errorf("Status after MarkReviewed() = %q, want %q", session.Status(), StatusReviewed)
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
