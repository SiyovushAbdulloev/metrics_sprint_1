package httpserver

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewServer_Default(t *testing.T) {
	s := New()
	assert.NotNil(t, s)
	assert.NotNil(t, s.App)
	assert.Equal(t, "", s.Address)
}

func TestNewServer_WithAddress(t *testing.T) {
	s := New(WithAddress(":8080"))
	assert.NotNil(t, s)
	assert.Equal(t, ":8080", s.Address)
}

func TestWithAddress(t *testing.T) {
	s := &Server{}
	opt := WithAddress(":9000")
	opt(s)
	assert.Equal(t, ":9000", s.Address)
}
