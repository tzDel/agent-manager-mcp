package application

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tzDel/agent-manager-mcp/internal/domain"
)

type CreateWorktreeRequest struct {
	AgentID string
}

type CreateWorktreeResponse struct {
	AgentID      string
	WorktreePath string
	BranchName   string
	Status       string
}

type CreateWorktreeUseCase struct {
	gitOps      domain.GitOperations
	agentRepo   domain.AgentRepository
	repoRoot    string
	worktreeDir string
}

func NewCreateWorktreeUseCase(
	gitOps domain.GitOperations,
	agentRepo domain.AgentRepository,
	repoRoot string,
) *CreateWorktreeUseCase {
	return &CreateWorktreeUseCase{
		gitOps:      gitOps,
		agentRepo:   agentRepo,
		repoRoot:    repoRoot,
		worktreeDir: filepath.Join(repoRoot, ".worktrees"),
	}
}

func (useCase *CreateWorktreeUseCase) Execute(ctx context.Context, request CreateWorktreeRequest) (*CreateWorktreeResponse, error) {
	agentID, err := domain.NewAgentID(request.AgentID)
	if err != nil {
		return nil, fmt.Errorf("invalid agent ID: %w", err)
	}

	exists, err := useCase.agentRepo.Exists(ctx, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check agent existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("agent already exists: %s", agentID.String())
	}

	branchName := agentID.BranchName()
	branchExists, err := useCase.gitOps.BranchExists(ctx, branchName)
	if err != nil {
		return nil, fmt.Errorf("failed to check branch existence: %w", err)
	}
	if branchExists {
		return nil, fmt.Errorf("branch already exists: %s", branchName)
	}

	worktreePath := filepath.Join(useCase.worktreeDir, agentID.WorktreeDirName())

	if err := useCase.gitOps.CreateWorktree(ctx, worktreePath, branchName); err != nil {
		return nil, fmt.Errorf("failed to create worktree: %w", err)
	}

	agent, err := domain.NewAgent(agentID, worktreePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	if err := useCase.agentRepo.Save(ctx, agent); err != nil {
		return nil, fmt.Errorf("failed to save agent: %w", err)
	}

	return &CreateWorktreeResponse{
		AgentID:      agent.ID().String(),
		WorktreePath: agent.WorktreePath(),
		BranchName:   agent.BranchName(),
		Status:       string(agent.Status()),
	}, nil
}
