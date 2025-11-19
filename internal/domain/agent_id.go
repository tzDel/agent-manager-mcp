package domain

import (
	"errors"
	"regexp"
	"strings"
)

type AgentID struct {
	value string
}

var agentIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)

func NewAgentID(rawID string) (AgentID, error) {
	normalized := strings.ToLower(strings.TrimSpace(rawID))

	if len(normalized) < 2 || len(normalized) > 50 {
		return AgentID{}, errors.New("agent ID must be 2-50 characters")
	}

	if !agentIDPattern.MatchString(normalized) {
		return AgentID{}, errors.New("agent ID must contain only lowercase letters, numbers, and hyphens")
	}

	return AgentID{value: normalized}, nil
}

func (agentID AgentID) String() string {
	return agentID.value
}

func (agentID AgentID) BranchName() string {
	return "agent-" + agentID.value
}

func (agentID AgentID) WorktreeDirName() string {
	return "agent-" + agentID.value
}
