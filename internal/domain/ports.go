package domain

import "context"

type GitOperations interface {
	CreateWorktree(ctx context.Context, worktreePath string, branchName string) error
	RemoveWorktree(ctx context.Context, worktreePath string) error
	BranchExists(ctx context.Context, branchName string) (bool, error)
}

type AgentRepository interface {
	Save(ctx context.Context, agent *Agent) error
	FindByID(ctx context.Context, agentID AgentID) (*Agent, error)
	Exists(ctx context.Context, agentID AgentID) (bool, error)
}
