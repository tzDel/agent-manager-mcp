package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDatabasePath_WhenDatabaseDirectoryNotProvided_UsesDefaultDirectory(t *testing.T) {
	// arrange
	defaultDirectory := t.TempDir()
	originalDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalDirectory)
	})
	if chdirErr := os.Chdir(defaultDirectory); chdirErr != nil {
		t.Fatalf("failed to change working directory: %v", chdirErr)
	}

	// act
	databasePath, err := resolveDatabasePath("")

	// assert
	if err != nil {
		t.Fatalf("resolveDatabasePath returned error: %v", err)
	}

	expectedDatabasePath := filepath.Join(defaultDirectory, databaseFileName)
	if databasePath != expectedDatabasePath {
		t.Fatalf("database path = %s, expected %s", databasePath, expectedDatabasePath)
	}
}

func TestResolveDatabasePath_WhenCustomDirectoryProvided_CreatesDatabaseInThatDirectory(t *testing.T) {
	// arrange
	workingDirectory := t.TempDir()
	originalDirectory, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalDirectory)
	})
	if chdirErr := os.Chdir(workingDirectory); chdirErr != nil {
		t.Fatalf("failed to change working directory: %v", chdirErr)
	}

	databaseDirectory := filepath.Join("data", "sqlite")

	// act
	databasePath, err := resolveDatabasePath(databaseDirectory)

	// assert
	if err != nil {
		t.Fatalf("resolveDatabasePath returned error: %v", err)
	}

	expectedDirectory := filepath.Join(workingDirectory, databaseDirectory)
	if _, statErr := os.Stat(expectedDirectory); statErr != nil {
		t.Fatalf("expected database directory %s to be created: %v", expectedDirectory, statErr)
	}

	expectedDatabasePath := filepath.Join(expectedDirectory, databaseFileName)
	if databasePath != expectedDatabasePath {
		t.Fatalf("database path = %s, expected %s", databasePath, expectedDatabasePath)
	}
}
