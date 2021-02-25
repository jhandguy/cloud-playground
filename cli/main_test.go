package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jhandguy/devops-playground/cli/message"
)

func TestCommands(t *testing.T) {
	assert.Equal(t, len(cmd.Commands()), 1)
	assert.Equal(t, cmd.Commands()[0], message.Cmd)
}
