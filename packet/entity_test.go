package packet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidMessageString(t *testing.T) {
	err1 := IsValidMessageString(`{}`)
	err2 := IsValidMessageString(`[{ "type": "text" }]`)
	err3 := IsValidMessageString(`[{ "type": "invalid type", "content": "text value" }]`)
	err4 := IsValidMessageString(`[{ "type": "text", "content": "text value" }]`)

	assert.Error(t, err1)
	assert.Error(t, err2)
	assert.Error(t, err3)
	assert.Nil(t, err4)
}
