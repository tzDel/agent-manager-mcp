package application

import (
	"context"
	"fmt"
	"time"

	"github.com/tzDel/orchestrAIgent/internal/domain"
)

type RemoveWorktreeRequest struct {
	SessionID string
	Force     bool
}

type RemoveWorktreeResponse struct {
	SessionID          string    `json:"sessionId"`
	RemovedAt          time.Time `json:"removedAt,omitempty"`
	HasUnmergedChanges bool      `json:"hasUnmergedChanges"`
	UnmergedCommits    int       `json:"unmergedCommits"`
	UncommittedFiles   int       `json:"uncommittedFiles"`
	Warning            string    `json:"warning,omitempty"`
}

type RemoveWorktreeUseCase struct {
	gitOperations     domain.GitOperations
	sessionRepository domain.SessionRepository
	baseBranch        string
}

func NewRemoveWorktreeUseCase(
	gitOperations domain.GitOperations,
	sessionRepository domain.SessionRepository,
	baseBranch string,
) *RemoveWorktreeUseCase {
	return &RemoveWorktreeUseCase{
		gitOperations:     gitOperations,
		sessionRepository: sessionRepository,
		baseBranch:        baseBranch,
	}
}

func (removeWorktreeUseCase *RemoveWorktreeUseCase) Execute(
	ctx context.Context,
	request RemoveWorktreeRequest,
) (*RemoveWorktreeResponse, error) {
	sessionID, err := removeWorktreeUseCase.validateSessionID(request.SessionID)
	if err != nil {
		return nil, err
	}
	session, err := removeWorktreeUseCase.fetchSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if err := removeWorktreeUseCase.verifySessionNotRemoved(session); err != nil {
		return nil, err
	}

	response := &RemoveWorktreeResponse{
		SessionID: request.SessionID,
	}

	if !request.Force {
		err := removeWorktreeUseCase.checkForUnmergedWork(ctx, session, response)
		if err != nil {
			return nil, err
		}
		if response.HasUnmergedChanges {
			return response, nil
		}
	}

	if err := removeWorktreeUseCase.removeWorktree(ctx, session, request.Force); err != nil {
		return nil, err
	}
	removeWorktreeUseCase.deleteBranchIfPossible(ctx, session)
	if err := removeWorktreeUseCase.markSessionRemoved(ctx, session); err != nil {
		return nil, err
	}

	response.RemovedAt = time.Now()
	response.HasUnmergedChanges = false
	return response, nil
}

func (removeWorktreeUseCase *RemoveWorktreeUseCase) validateSessionID(sessionIDString string) (domain.SessionID, error) {
	sessionID, err := domain.NewSessionID(sessionIDString)
	if err != nil {
		return domain.SessionID{}, fmt.Errorf("invalid session ID: %w", err)
	}
	return sessionID, nil
}

func (removeWorktreeUseCase *RemoveWorktreeUseCase) fetchSession(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error) {
	session, err := removeWorktreeUseCase.sessionRepository.FindByID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}
	return session, nil
}

func (removeWorktreeUseCase *RemoveWorktreeUseCase) verifySessionNotRemoved(session *domain.Session) error {
	if session.Status() == domain.StatusRemoved {
		return fmt.Errorf("session already removed")
	}
	return nil
}

func (removeWorktreeUseCase *RemoveWorktreeUseCase) checkForUnmergedWork(
	ctx context.Context,
	session *domain.Session,
	response *RemoveWorktreeResponse,
) error {
	hasUncommitted, fileCount, err := removeWorktreeUseCase.gitOperations.HasUncommittedChanges(ctx, session.WorktreePath())
	if err != nil {
		return fmt.Errorf("failed to check uncommitted changes: %w", err)
	}

	unpushedCount, err := removeWorktreeUseCase.gitOperations.HasUnpushedCommits(ctx, removeWorktreeUseCase.baseBranch, session.BranchName())
	if err != nil {
		return fmt.Errorf("failed to check unpushed commits: %w", err)
	}

	response.UncommittedFiles = fileCount
	response.UnmergedCommits = unpushedCount

	if hasUncommitted || unpushedCount > 0 {
		removeWorktreeUseCase.setUnmergedChanges(response, unpushedCount, fileCount)
	}

	return nil
}

func (removeWorktreeUseCase *RemoveWorktreeUseCase) setUnmergedChanges(response *RemoveWorktreeResponse, unpushedCount int, fileCount int) {
	response.HasUnmergedChanges = true
	response.Warning = fmt.Sprintf(
		"Session has %d unpushed commits and %d uncommitted files. Call with force=true to remove anyway.",
		unpushedCount,
		fileCount,
	)
}

func (removeWorktreeUseCase *RemoveWorktreeUseCase) removeWorktree(ctx context.Context, session *domain.Session, force bool) error {
	if err := removeWorktreeUseCase.gitOperations.RemoveWorktree(ctx, session.WorktreePath(), force); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}
	return nil
}

func (removeWorktreeUseCase *RemoveWorktreeUseCase) deleteBranchIfPossible(ctx context.Context, session *domain.Session) {
	removeWorktreeUseCase.gitOperations.DeleteBranch(ctx, session.BranchName(), true)
}

func (removeWorktreeUseCase *RemoveWorktreeUseCase) markSessionRemoved(ctx context.Context, session *domain.Session) error {
	session.MarkRemoved()
	if err := removeWorktreeUseCase.sessionRepository.Save(ctx, session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}
	return nil
}
