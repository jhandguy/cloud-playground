package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommands(t *testing.T) {
	assert.Equal(t, len(Cmd.Commands()), 3)
	assert.Equal(t, Cmd.Commands()[0], createMessageCmd)
	assert.Equal(t, Cmd.Commands()[1], deleteMessageCmd)
	assert.Equal(t, Cmd.Commands()[2], getMessageCmd)
}
