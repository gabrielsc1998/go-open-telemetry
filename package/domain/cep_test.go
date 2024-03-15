package domain

import "testing"

func TestNewCep(t *testing.T) {
	t.Run("valid cep", func(t *testing.T) {
		_, err := NewCep("12345678")
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("invalid cep", func(t *testing.T) {
		_, err := NewCep("invalid")
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
