package utils

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateUUID(t *testing.T) {
	ids := "\n\n"

	for i := 0; i < 5; i++ {
		assert.NotPanics(t, func() {
			id := uuid.New()
			ids += " - " + id.String() + "\n"
		}, "some uuid fail")
	}

	t.Log(ids)
}

func TestGenerateHashPassword(t *testing.T) {
	password := "password1234"

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	assert.NoError(t, err, "should be nil")

	t.Logf("\n\n - %s\n\n", string(hash))
}
