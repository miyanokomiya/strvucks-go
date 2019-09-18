package strava

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	result := Config()

	assert.Equal(t, result.Scopes, []string{"read,activity:read_all"}, "correct scopes")
}

func TestAuthCodeOption(t *testing.T) {
	result := AuthCodeOption()

	assert.Equal(t, len(result), 2, "correct option")
}

func TestClient(t *testing.T) {
	assert.Panics(t, func() {
		Client(nil)
	}, "invalid args")
}
