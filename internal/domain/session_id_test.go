package domain

import "testing"

func TestNewSessionID_Valid(t *testing.T) {
	// arrange
	tests := []struct {
		input    string
		expected string
	}{
		{"test-session", "test-session"},
		{"Test-Session", "test-session"},
		{"  copilot-123  ", "copilot-123"},
		{"a1", "a1"},
	}

	for _, testCase := range tests {
		// act
		sessionID, err := NewSessionID(testCase.input)

		// assert
		if err != nil {
			t.Errorf("NewSessionID(%q) unexpected error: %v", testCase.input, err)
		}
		if sessionID.String() != testCase.expected {
			t.Errorf("NewSessionID(%q) = %q, want %q", testCase.input, sessionID.String(), testCase.expected)
		}
	}
}

func TestNewSessionID_Invalid(t *testing.T) {
	// arrange
	tests := []string{
		"",
		"a",
		"Test_Session",
		"test session",
		"-test",
		"test-",
	}

	for _, input := range tests {
		// act
		_, err := NewSessionID(input)

		// assert
		if err == nil {
			t.Errorf("NewSessionID(%q) expected error, got nil", input)
		}
	}
}

func TestSessionID_BranchName(t *testing.T) {
	// arrange
	sessionID, _ := NewSessionID("copilot-123")
	expected := "session-copilot-123"

	// act
	result := sessionID.BranchName()

	// assert
	if result != expected {
		t.Errorf("BranchName() = %q, want %q", result, expected)
	}
}
