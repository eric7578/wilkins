package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_isErrMessageContains(t *testing.T) {
	b1 := isErrMessageContains(errors.New("foo bar"), "foo")
	b2 := isErrMessageContains(errors.New("foo bar"), "foo bar")
	b3 := isErrMessageContains(errors.New("foo bar"), "foo nan")
	b4 := isErrMessageContains(errors.New("foo bar"), "nan")
	b5 := isErrMessageContains(errors.New("foo bar"), "")
	b6 := isErrMessageContains(nil, "foo")

	assert.True(t, b1)
	assert.True(t, b2)
	assert.False(t, b3)
	assert.False(t, b4)
	assert.False(t, b5)
	assert.False(t, b6)
}
