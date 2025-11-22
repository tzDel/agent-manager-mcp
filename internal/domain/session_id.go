package domain

import (
	"errors"
	"regexp"
	"strings"
)

type SessionID struct {
	value string
}

var sessionIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)

func NewSessionID(rawID string) (SessionID, error) {
	normalized := strings.ToLower(strings.TrimSpace(rawID))

	if len(normalized) < 2 || len(normalized) > 50 {
		return SessionID{}, errors.New("session ID must be 2-50 characters")
	}

	if !sessionIDPattern.MatchString(normalized) {
		return SessionID{}, errors.New("session ID must contain only lowercase letters, numbers, and hyphens")
	}

	return SessionID{value: normalized}, nil
}

func (sessionID SessionID) String() string {
	return sessionID.value
}

func (sessionID SessionID) BranchName() string {
	return "session-" + sessionID.value
}

func (sessionID SessionID) WorktreeDirName() string {
	return "session-" + sessionID.value
}
