package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tzDel/orchestrAIgent/internal/adapters/mcp"
	"github.com/tzDel/orchestrAIgent/internal/application"
	"github.com/tzDel/orchestrAIgent/internal/infrastructure/git"
	"github.com/tzDel/orchestrAIgent/internal/infrastructure/persistence"
)

func main() {
	repositoryPath, databasePath := parseFlags()
	sessionRepository, cleanup := initializeSessionRepository(databasePath)
	defer cleanup()

	server := initializeMCPServer(repositoryPath, sessionRepository)
	startMCPServer(server, repositoryPath)
}

func parseFlags() (string, string) {
	repositoryPath := flag.String("repo", "", "path to git repository (defaults to current directory)")
	databasePath := flag.String("db", ".orchestrAIgent.db", "path to SQLite database (defaults to .orchestrAIgent.db in repository)")
	flag.Parse()

	repoPath := *repositoryPath
	if repoPath == "" {
		repoPath = resolveCurrentWorkingDirectory()
	}

	return repoPath, *databasePath
}

func resolveCurrentWorkingDirectory() string {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to resolve current working directory: %v", err)
	}
	return currentWorkingDirectory
}

func initializeSessionRepository(databasePath string) (*persistence.SQLiteSessionRepository, func()) {
	sessionRepository, err := persistence.NewSQLiteSessionRepository(databasePath)
	if err != nil {
		log.Fatalf("failed to initialize session repository: %v", err)
	}

	cleanup := func() {
		if err := sessionRepository.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}

	return sessionRepository, cleanup
}

func initializeMCPServer(repositoryPath string, sessionRepository *persistence.SQLiteSessionRepository) *mcp.MCPServer {
	gitOperations := git.NewGitClient(repositoryPath)
	baseBranch := "main"

	createWorktreeUseCase := application.NewCreateWorktreeUseCase(gitOperations, sessionRepository, repositoryPath)
	removeSessionUseCase := application.NewRemoveSessionUseCase(gitOperations, sessionRepository, baseBranch)
	getSessionsUseCase := application.NewGetSessionsUseCase(gitOperations, sessionRepository, baseBranch)

	server, err := mcp.NewMCPServer(createWorktreeUseCase, removeSessionUseCase, getSessionsUseCase)
	if err != nil {
		log.Fatalf("failed to initialize MCP server: %v", err)
	}

	return server
}

func startMCPServer(server *mcp.MCPServer, repositoryPath string) {
	fmt.Fprintf(os.Stderr, "Starting MCP server for repository: %s\n", repositoryPath)

	serverContext := context.Background()
	if err := server.Run(serverContext); err != nil {
		log.Fatalf("MCP server terminated with error: %v", err)
	}
}
