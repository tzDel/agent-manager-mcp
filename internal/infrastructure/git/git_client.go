package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Client struct {
	repoRoot string
}

func NewClient(repoRoot string) *Client {
	return &Client{repoRoot: repoRoot}
}

func (client *Client) CreateWorktree(ctx context.Context, worktreePath string, branchName string) error {
	cmd := exec.CommandContext(ctx, "git", "worktree", "add", "-b", branchName, worktreePath)
	cmd.Dir = client.repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree add failed: %w (output: %s)", err, string(output))
	}

	return nil
}

func (client *Client) RemoveWorktree(ctx context.Context, worktreePath string) error {
	cmd := exec.CommandContext(ctx, "git", "worktree", "remove", worktreePath, "--force")
	cmd.Dir = client.repoRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree remove failed: %w (output: %s)", err, string(output))
	}

	return nil
}

func (client *Client) BranchExists(ctx context.Context, branchName string) (bool, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "--list", branchName)
	cmd.Dir = client.repoRoot

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("git branch --list failed: %w", err)
	}

	return strings.TrimSpace(string(output)) != "", nil
}
