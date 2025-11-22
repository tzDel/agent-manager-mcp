package mcp

import (
	"context"
	"fmt"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/tzDel/orchestrAIgent/internal/application"
)

type CreateWorktreeArgs struct {
	SessionID string `json:"sessionId" jsonschema:"required" jsonschema_description:"The unique identifier for the session"`
}

type CreateWorktreeOutput struct {
	SessionID    string `json:"sessionId"`
	WorktreePath string `json:"worktreePath"`
	BranchName   string `json:"branchName"`
	Status       string `json:"status"`
}

type RemoveWorktreeArgs struct {
	SessionID string `json:"sessionId" jsonschema:"required" jsonschema_description:"Session identifier"`
	Force     bool   `json:"force" jsonschema_description:"Skip safety checks and force removal"`
}

type RemoveWorktreeOutput struct {
	SessionID          string `json:"sessionId"`
	RemovedAt          string `json:"removedAt,omitempty"`
	HasUnmergedChanges bool   `json:"hasUnmergedChanges"`
	UnmergedCommits    int    `json:"unmergedCommits"`
	UncommittedFiles   int    `json:"uncommittedFiles"`
	Warning            string `json:"warning,omitempty"`
}

type MCPServer struct {
	mcpServer             *mcpsdk.Server
	createWorktreeUseCase *application.CreateWorktreeUseCase
	removeWorktreeUseCase *application.RemoveWorktreeUseCase
}

func NewMCPServer(
	createWorktreeUseCase *application.CreateWorktreeUseCase,
	removeWorktreeUseCase *application.RemoveWorktreeUseCase,
) (*MCPServer, error) {
	impl := &mcpsdk.Implementation{
		Name:    "orchestrAIgent",
		Version: "0.1.0",
	}

	mcpServer := mcpsdk.NewServer(impl, nil)

	server := &MCPServer{
		mcpServer:             mcpServer,
		createWorktreeUseCase: createWorktreeUseCase,
		removeWorktreeUseCase: removeWorktreeUseCase,
	}

	mcpsdk.AddTool(
		mcpServer,
		&mcpsdk.Tool{
			Name:        "create_worktree",
			Description: "Creates an isolated git worktree for a specific session with its own branch",
		},
		server.handleCreateWorktree,
	)

	mcpsdk.AddTool(
		mcpServer,
		&mcpsdk.Tool{
			Name:        "remove_worktree",
			Description: "Removes an session's worktree and branch. Checks for unmerged changes unless force=true.",
		},
		server.handleRemoveWorktree,
	)

	return server, nil
}

func (s *MCPServer) handleCreateWorktree(
	ctx context.Context,
	req *mcpsdk.CallToolRequest,
	args CreateWorktreeArgs,
) (*mcpsdk.CallToolResult, any, error) {
	request := application.CreateWorktreeRequest{
		SessionID: args.SessionID,
	}

	response, err := s.createWorktreeUseCase.Execute(ctx, request)
	if err != nil {
		message := fmt.Sprintf("Failed to create worktree: %v", err)
		return newErrorResult(message), nil, err
	}

	output := CreateWorktreeOutput{
		SessionID:    response.SessionID,
		WorktreePath: response.WorktreePath,
		BranchName:   response.BranchName,
		Status:       response.Status,
	}

	message := fmt.Sprintf("Successfully created worktree for session '%s' at '%s' on branch '%s'", response.SessionID, response.WorktreePath, response.BranchName)
	return newSuccessResult(message), output, nil
}

func (s *MCPServer) handleRemoveWorktree(
	ctx context.Context,
	req *mcpsdk.CallToolRequest,
	args RemoveWorktreeArgs,
) (*mcpsdk.CallToolResult, any, error) {
	request := application.RemoveWorktreeRequest{
		SessionID: args.SessionID,
		Force:     args.Force,
	}

	response, err := s.removeWorktreeUseCase.Execute(ctx, request)
	if err != nil {
		message := fmt.Sprintf("Failed to remove worktree: %v", err)
		return newErrorResult(message), nil, err
	}

	output := RemoveWorktreeOutput{
		SessionID:          response.SessionID,
		HasUnmergedChanges: response.HasUnmergedChanges,
		UnmergedCommits:    response.UnmergedCommits,
		UncommittedFiles:   response.UncommittedFiles,
		Warning:            response.Warning,
	}

	if !response.RemovedAt.IsZero() {
		output.RemovedAt = response.RemovedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	if response.HasUnmergedChanges {
		message := fmt.Sprintf(
			"WARNING: Session '%s' has unmerged changes\n\nUncommitted files: %d\nUnpushed commits: %d\n\n%s",
			response.SessionID,
			response.UncommittedFiles,
			response.UnmergedCommits,
			response.Warning,
		)
		return &mcpsdk.CallToolResult{
			Content: []mcpsdk.Content{newTextContent(message)},
			IsError: false,
		}, output, nil
	}

	message := fmt.Sprintf("Successfully removed worktree for session '%s'", response.SessionID)
	return newSuccessResult(message), output, nil
}

func (s *MCPServer) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcpsdk.StdioTransport{})
}
