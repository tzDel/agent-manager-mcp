package application

import (
	"context"
	"errors"
	"testing"

	"github.com/tzDel/agent-manager-mcp/internal/domain"
)

type mockGitOperations struct {
	createWorktreeFunc func(ctx context.Context, path string, branch string) error
	branchExistsFunc   func(ctx context.Context, branch string) (bool, error)
}

func (mock *mockGitOperations) CreateWorktree(ctx context.Context, path string, branch string) error {
	if mock.createWorktreeFunc != nil {
		return mock.createWorktreeFunc(ctx, path, branch)
	}
	return nil
}

func (mock *mockGitOperations) RemoveWorktree(ctx context.Context, path string) error {
	return nil
}

func (mock *mockGitOperations) BranchExists(ctx context.Context, branch string) (bool, error) {
	if mock.branchExistsFunc != nil {
		return mock.branchExistsFunc(ctx, branch)
	}
	return false, nil
}

type mockAgentRepository struct {
	agents map[string]*domain.Agent
}

func newMockAgentRepository() *mockAgentRepository {
	return &mockAgentRepository{
		agents: make(map[string]*domain.Agent),
	}
}

func (mock *mockAgentRepository) Save(ctx context.Context, agent *domain.Agent) error {
	mock.agents[agent.ID().String()] = agent
	return nil
}

func (mock *mockAgentRepository) FindByID(ctx context.Context, agentID domain.AgentID) (*domain.Agent, error) {
	agent, exists := mock.agents[agentID.String()]
	if !exists {
		return nil, errors.New("not found")
	}
	return agent, nil
}

func (mock *mockAgentRepository) Exists(ctx context.Context, agentID domain.AgentID) (bool, error) {
	_, exists := mock.agents[agentID.String()]
	return exists, nil
}

func TestCreateWorktreeUseCase_Execute_Success(t *testing.T) {
	// arrange
	gitOps := &mockGitOperations{}
	agentRepo := newMockAgentRepository()
	useCase := NewCreateWorktreeUseCase(gitOps, agentRepo, "/repo/root")
	request := CreateWorktreeRequest{AgentID: "test-agent"}
	ctx := context.Background()

	// act
	response, err := useCase.Execute(ctx, request)

	// assert
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	if response.AgentID != "test-agent" {
		t.Errorf("AgentID = %q, want %q", response.AgentID, "test-agent")
	}

	if response.BranchName != "agent-test-agent" {
		t.Errorf("BranchName = %q, want %q", response.BranchName, "agent-test-agent")
	}

	if response.Status != "created" {
		t.Errorf("Status = %q, want %q", response.Status, "created")
	}
}

func TestCreateWorktreeUseCase_Execute_InvalidAgentID(t *testing.T) {
	// arrange
	gitOps := &mockGitOperations{}
	agentRepo := newMockAgentRepository()
	useCase := NewCreateWorktreeUseCase(gitOps, agentRepo, "/repo/root")
	request := CreateWorktreeRequest{AgentID: "Invalid_ID"}
	ctx := context.Background()

	// act
	_, err := useCase.Execute(ctx, request)

	// assert
	if err == nil {
		t.Error("Execute() expected error for invalid agent ID")
	}
}

func TestCreateWorktreeUseCase_Execute_AgentAlreadyExists(t *testing.T) {
	// arrange
	gitOps := &mockGitOperations{}
	agentRepo := newMockAgentRepository()
	useCase := NewCreateWorktreeUseCase(gitOps, agentRepo, "/repo/root")

	agentID, _ := domain.NewAgentID("test-agent")
	agent, _ := domain.NewAgent(agentID, "/path")
	agentRepo.Save(context.Background(), agent)

	request := CreateWorktreeRequest{AgentID: "test-agent"}
	ctx := context.Background()

	// act
	_, err := useCase.Execute(ctx, request)

	// assert
	if err == nil {
		t.Error("Execute() expected error for existing agent")
	}
}

func TestCreateWorktreeUseCase_Execute_BranchAlreadyExists(t *testing.T) {
	// arrange
	gitOps := &mockGitOperations{
		branchExistsFunc: func(ctx context.Context, branch string) (bool, error) {
			return true, nil
		},
	}
	agentRepo := newMockAgentRepository()
	useCase := NewCreateWorktreeUseCase(gitOps, agentRepo, "/repo/root")
	request := CreateWorktreeRequest{AgentID: "test-agent"}
	ctx := context.Background()

	// act
	_, err := useCase.Execute(ctx, request)

	// assert
	if err == nil {
		t.Error("Execute() expected error for existing branch")
	}
}

func TestCreateWorktreeUseCase_Execute_GitOperationFails(t *testing.T) {
	// arrange
	gitOps := &mockGitOperations{
		createWorktreeFunc: func(ctx context.Context, path string, branch string) error {
			return errors.New("git error")
		},
	}
	agentRepo := newMockAgentRepository()
	useCase := NewCreateWorktreeUseCase(gitOps, agentRepo, "/repo/root")
	request := CreateWorktreeRequest{AgentID: "test-agent"}
	ctx := context.Background()

	// act
	_, err := useCase.Execute(ctx, request)

	// assert
	if err == nil {
		t.Error("Execute() expected error when git operation fails")
	}
}
