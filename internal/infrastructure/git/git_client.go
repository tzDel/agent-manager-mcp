package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type GitClient struct {
	repositoryRoot string
}

func NewGitClient(repositoryRoot string) *GitClient {
	return &GitClient{repositoryRoot: repositoryRoot}
}

// executeGitCommand executes a git command in the repository root directory
// and returns the combined output (stdout and stderr) along with any error
func (gitClient *GitClient) executeGitCommand(ctx context.Context, args ...string) ([]byte, error) {
	gitCommand := exec.CommandContext(ctx, "git", args...)
	gitCommand.Dir = gitClient.repositoryRoot

	commandOutput, err := gitCommand.CombinedOutput()
	if err != nil {
		return commandOutput, fmt.Errorf("git command failed: %w (output: %s)", err, string(commandOutput))
	}

	return commandOutput, nil
}

// executeGitCommandWithOutput executes a git command and returns only stdout
// Used for commands where we need to parse the output (like branch --list)
func (gitClient *GitClient) executeGitCommandWithOutput(ctx context.Context, args ...string) ([]byte, error) {
	gitCommand := exec.CommandContext(ctx, "git", args...)
	gitCommand.Dir = gitClient.repositoryRoot

	commandOutput, err := gitCommand.Output()
	if err != nil {
		return nil, fmt.Errorf("git command failed: %w", err)
	}

	return commandOutput, nil
}

func (gitClient *GitClient) CreateWorktree(ctx context.Context, worktreePath string, branchName string) error {
	_, err := gitClient.executeGitCommand(ctx, "worktree", "add", "-b", branchName, worktreePath)
	if err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	return nil
}

func (gitClient *GitClient) RemoveWorktree(ctx context.Context, worktreePath string) error {
	_, err := gitClient.executeGitCommand(ctx, "worktree", "remove", worktreePath, "--force")
	if err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	return nil
}

func (gitClient *GitClient) BranchExists(ctx context.Context, branchName string) (bool, error) {
	commandOutput, err := gitClient.executeGitCommandWithOutput(ctx, "branch", "--list", branchName)
	if err != nil {
		return false, fmt.Errorf("failed to check branch existence: %w", err)
	}

	return strings.TrimSpace(string(commandOutput)) != "", nil
}
