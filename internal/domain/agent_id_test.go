package domain

import "testing"

func TestNewAgentID_Valid(t *testing.T) {
	// arrange
	tests := []struct {
		input    string
		expected string
	}{
		{"test-agent", "test-agent"},
		{"Test-Agent", "test-agent"},
		{"  copilot-123  ", "copilot-123"},
		{"a1", "a1"},
	}

	for _, testCase := range tests {
		// act
		agentID, err := NewAgentID(testCase.input)

		// assert
		if err != nil {
			t.Errorf("NewAgentID(%q) unexpected error: %v", testCase.input, err)
		}
		if agentID.String() != testCase.expected {
			t.Errorf("NewAgentID(%q) = %q, want %q", testCase.input, agentID.String(), testCase.expected)
		}
	}
}

func TestNewAgentID_Invalid(t *testing.T) {
	// arrange
	tests := []string{
		"",
		"a",
		"Test_Agent",
		"test agent",
		"-test",
		"test-",
	}

	for _, input := range tests {
		// act
		_, err := NewAgentID(input)

		// assert
		if err == nil {
			t.Errorf("NewAgentID(%q) expected error, got nil", input)
		}
	}
}

func TestAgentID_BranchName(t *testing.T) {
	// arrange
	agentID, _ := NewAgentID("copilot-123")
	expected := "agent-copilot-123"

	// act
	result := agentID.BranchName()

	// assert
	if result != expected {
		t.Errorf("BranchName() = %q, want %q", result, expected)
	}
}
