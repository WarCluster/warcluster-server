package libs

import (
    "testing"
)

func TestUsernameHash(t *testing.T) {
    hash := usernameHash("Gopher")
    if len(hash) != 64 {
        t.Error("Wrong hash length")
    }
}

func TestSimplifyHash(t *testing.T) {
    hash := simplifyHash(usernameHash("Gopher"))
    if len(hash) != 64 {
        t.Error("Wrong simplified hash length")
    }
}
