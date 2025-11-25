package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tzDel/orchestragent-mcp/internal/adapters/mcp"
	"github.com/tzDel/orchestragent-mcp/internal/application"
	"github.com/tzDel/orchestragent-mcp/internal/infrastructure/git"
	"github.com/tzDel/orchestragent-mcp/internal/infrastructure/persistence"
)

const databaseFileName = ".orchestragent-mcp.db"
const defaultDatabaseDirectory = "."

func main() {
	repositoryPath, databaseDirectory := parseFlags()

	databasePath, err := resolveDatabasePath(databaseDirectory)
	if err != nil {
		log.Fatalf("failed to resolve database path: %v", err)
	}

	sessionRepository, cleanup := initializeSessionRepository(databasePath)
	defer cleanup()

	server := initializeMCPServer(repositoryPath, sessionRepository)
	startMCPServer(server, repositoryPath)
}

func parseFlags() (string, string) {
	repositoryPath := flag.String("repo", resolveCurrentWorkingDirectory(), "path to git repository (defaults to current directory)")
	databaseDirectory := flag.String("db", defaultDatabaseDirectory, "directory where SQLite database should be created (defaults to current working directory)")
	flag.Parse()

	return *repositoryPath, *databaseDirectory
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

func resolveDatabasePath(databaseDirectory string) (string, error) {
	baseDirectory := databaseDirectory
	if baseDirectory == "" {
		baseDirectory = defaultDatabaseDirectory
	}

	if !filepath.IsAbs(baseDirectory) {
		absoluteDirectory, err := filepath.Abs(baseDirectory)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path for %s: %w", baseDirectory, err)
		}
		baseDirectory = absoluteDirectory
	}

	if err := os.MkdirAll(baseDirectory, 0o755); err != nil {
		return "", fmt.Errorf("failed to create database directory %s: %w", baseDirectory, err)
	}

	return filepath.Join(baseDirectory, databaseFileName), nil
}
