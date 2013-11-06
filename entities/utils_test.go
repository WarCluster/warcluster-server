package entities

import "testing"

func TestUsernameHash(t *testing.T) {
	hash := generateHash("Gopher")
	if len(hash) != 64 {
		t.Error("Wrong hash length")
	}
}
