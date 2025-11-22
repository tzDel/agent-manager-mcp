package domain

import (
	"errors"
	"time"
)

type SessionStatus string

const (
	StatusCreated SessionStatus = "created"
	StatusMerged  SessionStatus = "merged"
	StatusFailed  SessionStatus = "failed"
	StatusRemoved SessionStatus = "removed"
)

type Session struct {
	id           SessionID
	status       SessionStatus
	worktreePath string
	branchName   string
	createdAt    time.Time
	updatedAt    time.Time
}

func NewSession(sessionID SessionID, worktreePath string) (*Session, error) {
	if worktreePath == "" {
		return nil, errors.New("worktree path cannot be empty")
	}

	now := time.Now()
	return &Session{
		id:           sessionID,
		status:       StatusCreated,
		worktreePath: worktreePath,
		branchName:   sessionID.BranchName(),
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func (session *Session) ID() SessionID {
	return session.id
}

func (session *Session) Status() SessionStatus {
	return session.status
}

func (session *Session) WorktreePath() string {
	return session.worktreePath
}

func (session *Session) BranchName() string {
	return session.branchName
}

func (session *Session) MarkMerged() {
	session.status = StatusMerged
	session.updatedAt = time.Now()
}

func (session *Session) MarkFailed() {
	session.status = StatusFailed
	session.updatedAt = time.Now()
}

func (session *Session) MarkRemoved() {
	session.status = StatusRemoved
	session.updatedAt = time.Now()
}
