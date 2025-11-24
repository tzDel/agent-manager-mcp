package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/tzDel/orchestrAIgent/internal/domain"
	_ "modernc.org/sqlite"
)

type SQLiteSessionRepository struct {
	database *sql.DB
}

const createTableSQL = `
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    status TEXT NOT NULL CHECK(status IN ('open', 'reviewed', 'merged')),
    worktree_path TEXT NOT NULL,
    branch_name TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);
`

func NewSQLiteSessionRepository(databasePath string) (*SQLiteSessionRepository, error) {
	database, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	repository := &SQLiteSessionRepository{database: database}

	if err := repository.initializeSchema(); err != nil {
		database.Close()
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	return repository, nil
}

func (repository *SQLiteSessionRepository) initializeSchema() error {
	_, err := repository.database.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create sessions table: %w", err)
	}
	return nil
}

func (repository *SQLiteSessionRepository) Save(ctx context.Context, session *domain.Session) error {
	query := `
		INSERT INTO sessions (id, status, worktree_path, branch_name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			status = excluded.status,
			worktree_path = excluded.worktree_path,
			branch_name = excluded.branch_name,
			updated_at = excluded.updated_at
	`

	createdAt := time.Now().Unix()
	updatedAt := time.Now().Unix()

	_, err := repository.database.ExecContext(
		ctx,
		query,
		session.ID().String(),
		string(session.Status()),
		session.WorktreePath(),
		session.BranchName(),
		createdAt,
		updatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save session %s: %w", session.ID().String(), err)
	}

	return nil
}

func (repository *SQLiteSessionRepository) FindByID(ctx context.Context, sessionID domain.SessionID) (*domain.Session, error) {
	query := `
		SELECT id, status, worktree_path, branch_name, created_at, updated_at
		FROM sessions
		WHERE id = ?
	`

	var id, status, worktreePath, branchName string
	var createdAt, updatedAt int64

	err := repository.database.QueryRowContext(ctx, query, sessionID.String()).Scan(
		&id,
		&status,
		&worktreePath,
		&branchName,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found: %s", sessionID.String())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query session %s: %w", sessionID.String(), err)
	}

	return repository.reconstructSession(id, status, worktreePath)
}

func (repository *SQLiteSessionRepository) FindAll(ctx context.Context) ([]*domain.Session, error) {
	query := `
		SELECT id, status, worktree_path, branch_name, created_at, updated_at
		FROM sessions
		ORDER BY created_at ASC
	`

	rows, err := repository.database.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all sessions: %w", err)
	}
	defer rows.Close()

	sessions := make([]*domain.Session, 0)

	for rows.Next() {
		session, err := repository.scanRowIntoSession(rows)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating session rows: %w", err)
	}

	return sessions, nil
}

func (repository *SQLiteSessionRepository) Exists(ctx context.Context, sessionID domain.SessionID) (bool, error) {
	query := `SELECT COUNT(*) FROM sessions WHERE id = ?`

	var count int
	err := repository.database.QueryRowContext(ctx, query, sessionID.String()).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if session exists: %w", err)
	}

	return count > 0, nil
}

func (repository *SQLiteSessionRepository) Delete(ctx context.Context, sessionID domain.SessionID) error {
	query := `DELETE FROM sessions WHERE id = ?`

	result, err := repository.database.ExecContext(ctx, query, sessionID.String())
	if err != nil {
		return fmt.Errorf("failed to delete session %s: %w", sessionID.String(), err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID.String())
	}

	return nil
}

func (repository *SQLiteSessionRepository) Close() error {
	if repository.database != nil {
		return repository.database.Close()
	}
	return nil
}

func (repository *SQLiteSessionRepository) scanRowIntoSession(rows *sql.Rows) (*domain.Session, error) {
	var id, status, worktreePath, branchName string
	var createdAt, updatedAt int64

	err := rows.Scan(&id, &status, &worktreePath, &branchName, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan session row: %w", err)
	}

	session, err := repository.reconstructSession(id, status, worktreePath)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (repository *SQLiteSessionRepository) reconstructSession(id, status, worktreePath string) (*domain.Session, error) {
	sessionID, err := domain.NewSessionID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct session ID: %w", err)
	}

	session, err := domain.NewSession(sessionID, worktreePath)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstruct session: %w", err)
	}

	switch domain.SessionStatus(status) {
	case domain.StatusReviewed:
		session.MarkReviewed()
	case domain.StatusMerged:
		session.MarkMerged()
	}

	return session, nil
}
