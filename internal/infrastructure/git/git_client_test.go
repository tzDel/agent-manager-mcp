package git

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	cmd.Run()

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	cmd.Run()

	testFile := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test Repo"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd = exec.Command("git", "add", "README.md")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestClient_CreateWorktree(t *testing.T) {
	// arrange
	repoRoot, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient(repoRoot)
	ctx := context.Background()
	worktreePath := filepath.Join(repoRoot, ".worktrees", "test-agent")
	branchName := "agent-test"

	// act
	err := client.CreateWorktree(ctx, worktreePath, branchName)

	// assert
	if err != nil {
		t.Fatalf("CreateWorktree() error: %v", err)
	}

	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Error("Worktree directory was not created")
	}

	exists, err := client.BranchExists(ctx, branchName)
	if err != nil {
		t.Fatalf("BranchExists() error: %v", err)
	}
	if !exists {
		t.Error("Branch was not created")
	}
}

func TestClient_RemoveWorktree(t *testing.T) {
	// arrange
	repoRoot, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient(repoRoot)
	ctx := context.Background()
	worktreePath := filepath.Join(repoRoot, ".worktrees", "test-agent")
	branchName := "agent-test"

	client.CreateWorktree(ctx, worktreePath, branchName)

	// act
	err := client.RemoveWorktree(ctx, worktreePath)

	// assert
	if err != nil {
		t.Fatalf("RemoveWorktree() error: %v", err)
	}

	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		t.Error("Worktree directory still exists after removal")
	}
}

func TestClient_BranchExists(t *testing.T) {
	// arrange
	repoRoot, cleanup := setupTestRepo(t)
	defer cleanup()

	client := NewClient(repoRoot)
	ctx := context.Background()

	// act
	exists, err := client.BranchExists(ctx, "nonexistent")

	// assert
	if err != nil {
		t.Fatalf("BranchExists() error: %v", err)
	}
	if exists {
		t.Error("BranchExists() returned true for non-existent branch")
	}

	// arrange
	worktreePath := filepath.Join(repoRoot, ".worktrees", "test-agent")
	branchName := "agent-test"
	client.CreateWorktree(ctx, worktreePath, branchName)

	// act
	exists, err = client.BranchExists(ctx, branchName)

	// assert
	if err != nil {
		t.Fatalf("BranchExists() error: %v", err)
	}
	if !exists {
		t.Error("BranchExists() returned false for existing branch")
	}
}
