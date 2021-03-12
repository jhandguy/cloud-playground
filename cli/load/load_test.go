package load

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommands(t *testing.T) {
	assert.Equal(t, len(Cmd.Commands()), 1)
	assert.Equal(t, Cmd.Commands()[0], testLoadCmd)
}
