package application

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tzDel/orchestrAIgent/internal/domain"
)

type CreateWorktreeRequest struct {
	SessionID string
}

type CreateWorktreeResponse struct {
	SessionID    string
	WorktreePath string
	BranchName   string
	Status       string
}

type CreateWorktreeUseCase struct {
	gitOperations     domain.GitOperations
	sessionRepository domain.SessionRepository
	repositoryRoot    string
	worktreeDirectory string
}

func NewCreateWorktreeUseCase(
	gitOperations domain.GitOperations,
	sessionRepository domain.SessionRepository,
	repositoryRoot string,
) *CreateWorktreeUseCase {
	return &CreateWorktreeUseCase{
		gitOperations:     gitOperations,
		sessionRepository: sessionRepository,
		repositoryRoot:    repositoryRoot,
		worktreeDirectory: filepath.Join(repositoryRoot, ".worktrees"),
	}
}

func (createWorktreeUseCase *CreateWorktreeUseCase) Execute(ctx context.Context, request CreateWorktreeRequest) (*CreateWorktreeResponse, error) {
	sessionID, err := createWorktreeUseCase.validateSessionID(request.SessionID)
	if err != nil {
		return nil, err
	}

	if err := createWorktreeUseCase.ensureSessionDoesNotExist(ctx, sessionID); err != nil {
		return nil, err
	}

	if err := createWorktreeUseCase.ensureBranchDoesNotExist(ctx, sessionID.BranchName()); err != nil {
		return nil, err
	}

	worktreePath := createWorktreeUseCase.buildWorktreePath(sessionID)

	if err := createWorktreeUseCase.createWorktreeAndBranch(ctx, worktreePath, sessionID.BranchName()); err != nil {
		return nil, err
	}

	session, err := createWorktreeUseCase.createAndSaveSession(ctx, sessionID, worktreePath)
	if err != nil {
		return nil, err
	}

	return createWorktreeUseCase.buildResponse(session), nil
}

func (createWorktreeUseCase *CreateWorktreeUseCase) validateSessionID(sessionIDString string) (domain.SessionID, error) {
	sessionID, err := domain.NewSessionID(sessionIDString)
	if err != nil {
		return domain.SessionID{}, fmt.Errorf("invalid session ID: %w", err)
	}
	return sessionID, nil
}

func (createWorktreeUseCase *CreateWorktreeUseCase) ensureSessionDoesNotExist(ctx context.Context, sessionID domain.SessionID) error {
	exists, err := createWorktreeUseCase.sessionRepository.Exists(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to check session existence: %w", err)
	}
	if exists {
		return fmt.Errorf("session already exists: %s", sessionID.String())
	}
	return nil
}

func (createWorktreeUseCase *CreateWorktreeUseCase) ensureBranchDoesNotExist(ctx context.Context, branchName string) error {
	branchExists, err := createWorktreeUseCase.gitOperations.BranchExists(ctx, branchName)
	if err != nil {
		return fmt.Errorf("failed to check branch existence: %w", err)
	}
	if branchExists {
		return fmt.Errorf("branch already exists: %s", branchName)
	}
	return nil
}

func (createWorktreeUseCase *CreateWorktreeUseCase) buildWorktreePath(sessionID domain.SessionID) string {
	return filepath.Join(createWorktreeUseCase.worktreeDirectory, sessionID.WorktreeDirName())
}

func (createWorktreeUseCase *CreateWorktreeUseCase) createWorktreeAndBranch(ctx context.Context, worktreePath string, branchName string) error {
	if err := createWorktreeUseCase.gitOperations.CreateWorktree(ctx, worktreePath, branchName); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}
	return nil
}

func (createWorktreeUseCase *CreateWorktreeUseCase) createAndSaveSession(ctx context.Context, sessionID domain.SessionID, worktreePath string) (*domain.Session, error) {
	session, err := domain.NewSession(sessionID, worktreePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	if err := createWorktreeUseCase.sessionRepository.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, nil
}

func (createWorktreeUseCase *CreateWorktreeUseCase) buildResponse(session *domain.Session) *CreateWorktreeResponse {
	return &CreateWorktreeResponse{
		SessionID:    session.ID().String(),
		WorktreePath: session.WorktreePath(),
		BranchName:   session.BranchName(),
		Status:       string(session.Status()),
	}
}
