package tests

import (
	"testing"

	"goproject/app/handlers"
)

func TestGenerateToken(t *testing.T) {
	expectedLength := 32

	token := handlers.GenerateTokenWrapper()

	if len(token) != expectedLength {
		t.Errorf("Generated token length is incorrect. Expected: %d, Got: %d", expectedLength, len(token))
	} else {
		t.Logf("Test passed: GenerateToken")
	}
}
