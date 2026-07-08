package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetBearerToken(t *testing.T) {
	secret := "Nwn6xysO0WDB2g3dfGyyEFCzv69jRNjeYn5KLeoOwNQ1Sr19nalCisR4qNIPeEuVy8FSkdfe26l6VFhwsOX69A"

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer Nwn6xysO0WDB2g3dfGyyEFCzv69jRNjeYn5KLeoOwNQ1Sr19nalCisR4qNIPeEuVy8FSkdfe26l6VFhwsOX69A")

	token, err := GetBearerToken(req.Header)

	if err != nil {
		t.Errorf("Failed: %v", err)
	}

	if token != secret {
		t.Errorf("Failed: token and secret do not match.\n\tsecret: %v\n\ttoken: %v\n", secret, token)
	}
}

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
