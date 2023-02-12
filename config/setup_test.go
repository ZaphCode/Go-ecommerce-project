package config

import (
	"testing"

	"github.com/ZaphCode/clean-arch/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestMustLoadConfig(t *testing.T) {
	assert.Panics(t, func() {
		MustLoadConfig("./random/uwu")
	})

	cfg := Get()

	assert.Empty(t, cfg.Api.Port, "should be empty")

	assert.NotPanics(t, func() {
		MustLoadConfig("./")
	})

	cfg = Get()

	assert.NotEmpty(t, cfg.Api.Port, "should not be empty")

	utils.PrettyPrint(cfg)
}

func TestMustLoadFirebaseConfig(t *testing.T) {
	assert.Panics(t, func() {
		MustLoadFirebaseConfig("./random/uwu")
	})

	app := GetFirebaseApp()

	assert.Nil(t, app, "should be nil")

	assert.NotPanics(t, func() {
		MustLoadFirebaseConfig("./")
	})

	app = GetFirebaseApp()

	assert.NotNil(t, app, "should not be nil")
}
