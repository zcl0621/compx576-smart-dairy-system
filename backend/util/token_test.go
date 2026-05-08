package util

import "testing"

func TestGenerateAgentToken(t *testing.T) {
	first, err := GenerateAgentToken()
	if err != nil {
		t.Fatalf("GenerateAgentToken first: %v", err)
	}

	second, err := GenerateAgentToken()
	if err != nil {
		t.Fatalf("GenerateAgentToken second: %v", err)
	}

	if len(first) < 32 {
		t.Fatalf("token length = %d, want at least 32", len(first))
	}
	if first == second {
		t.Fatalf("tokens should be unique")
	}
}
