package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValid(t *testing.T) {
	secret := "mysecret"
	id, _ := uuid.NewUUID()

	jwt, _ := MakeJWT(id, secret, time.Second)

	alsoID, _ := ValidateJWT(jwt, secret)

	if id != alsoID {
		t.Errorf("IDs mismatch: %v and %v", id, alsoID)
	}
}

func TestTimeout(t *testing.T) {
	secret := "mysecret"
	id, _ := uuid.NewUUID()
	jwt, _ := MakeJWT(id, secret, time.Microsecond)
	time.Sleep(time.Millisecond * 5)

	alsoID, _ := ValidateJWT(jwt, secret)

	if id == alsoID {
		t.Errorf("ID timeout: %v and %v", id, alsoID)
	}
}

func TestSecretmismatch(t *testing.T) {
	id, _ := uuid.NewUUID()
	jwt, _ := MakeJWT(id, "secret1", time.Microsecond)
	time.Sleep(time.Millisecond * 5)

	alsoID, _ := ValidateJWT(jwt, "secret2")

	if id == alsoID {
		t.Errorf("Secret mismatch: %v and %v", id, alsoID)
	}
}
