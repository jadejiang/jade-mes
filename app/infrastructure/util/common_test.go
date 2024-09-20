package util

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUUIDV4(t *testing.T) {
	uuidStr := NewUUIDV4()
	_, err := uuid.Parse(uuidStr)
	if err != nil {
		t.Errorf("Expected valid UUID, got error: %v", err)
	}
	if uuidStr == "" {
		t.Errorf("Expected non-empty UUID string, got empty string")
	}
}