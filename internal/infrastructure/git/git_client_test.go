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

	temporaryDirectory := t.TempDir()

	gitInitCommand := exec.Command("git", "init")
	gitInitCommand.Dir = temporaryDirectory
	if err := gitInitCommand.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}

	gitConfigNameCommand := exec.Command("git", "config", "user.name", "Test User")
	gitConfigNameCommand.Dir = temporaryDirectory
	gitConfigNameCommand.Run()

	gitConfigEmailCommand := exec.Command("git", "config", "user.email", "test@example.com")
	gitConfigEmailCommand.Dir = temporaryDirectory
	gitConfigEmailCommand.Run()

	testFilePath := filepath.Join(temporaryDirectory, "README.md")
	if err := os.WriteFile(testFilePath, []byte("# Test Repo"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	gitAddCommand := exec.Command("git", "add", "README.md")
	gitAddCommand.Dir = temporaryDirectory
	if err := gitAddCommand.Run(); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	gitCommitCommand := exec.Command("git", "commit", "-m", "Initial commit")
	gitCommitCommand.Dir = temporaryDirectory
	if err := gitCommitCommand.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(temporaryDirectory)
	}

	return temporaryDirectory, cleanup
}

func TestGitClient_CreateWorktree(t *testing.T) {
	// arrange
	repositoryRoot, cleanup := setupTestRepo(t)
	defer cleanup()

	gitClient := NewGitClient(repositoryRoot)
	ctx := context.Background()
	worktreePath := filepath.Join(repositoryRoot, ".worktrees", "test-agent")
	branchName := "agent-test"

	// act
	err := gitClient.CreateWorktree(ctx, worktreePath, branchName)

	// assert
	if err != nil {
		t.Fatalf("CreateWorktree() error: %v", err)
	}

	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		t.Error("Worktree directory was not created")
	}

	exists, err := gitClient.BranchExists(ctx, branchName)
	if err != nil {
		t.Fatalf("BranchExists() error: %v", err)
	}
	if !exists {
		t.Error("Branch was not created")
	}
}

func TestGitClient_RemoveWorktree(t *testing.T) {
	// arrange
	repositoryRoot, cleanup := setupTestRepo(t)
	defer cleanup()

	gitClient := NewGitClient(repositoryRoot)
	ctx := context.Background()
	worktreePath := filepath.Join(repositoryRoot, ".worktrees", "test-agent")
	branchName := "agent-test"

	gitClient.CreateWorktree(ctx, worktreePath, branchName)

	// act
	err := gitClient.RemoveWorktree(ctx, worktreePath)

	// assert
	if err != nil {
		t.Fatalf("RemoveWorktree() error: %v", err)
	}

	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
		t.Error("Worktree directory still exists after removal")
	}
}

func TestGitClient_BranchExists(t *testing.T) {
	// arrange
	repositoryRoot, cleanup := setupTestRepo(t)
	defer cleanup()

	gitClient := NewGitClient(repositoryRoot)
	ctx := context.Background()

	// act
	exists, err := gitClient.BranchExists(ctx, "nonexistent")

	// assert
	if err != nil {
		t.Fatalf("BranchExists() error: %v", err)
	}
	if exists {
		t.Error("BranchExists() returned true for non-existent branch")
	}

	// arrange
	worktreePath := filepath.Join(repositoryRoot, ".worktrees", "test-agent")
	branchName := "agent-test"
	gitClient.CreateWorktree(ctx, worktreePath, branchName)

	// act
	exists, err = gitClient.BranchExists(ctx, branchName)

	// assert
	if err != nil {
		t.Fatalf("BranchExists() error: %v", err)
	}
	if !exists {
		t.Error("BranchExists() returned false for existing branch")
	}
}
