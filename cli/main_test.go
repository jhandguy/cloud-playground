package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/jhandguy/cloud-playground/cli/load"
	"github.com/jhandguy/cloud-playground/cli/message"
)

type logger struct {
	*testing.T

	Messages []string
}

func newLogger(t *testing.T) *logger {
	return &logger{T: t}
}

func (t *logger) Logf(format string, args ...interface{}) {
	t.Messages = append(t.Messages, fmt.Sprintf(format, args...))
}

func TestCommands(t *testing.T) {
	assert.Equal(t, len(cmd.Commands()), 2)
	assert.Equal(t, cmd.Commands()[0], load.Cmd)
	assert.Equal(t, cmd.Commands()[1], message.Cmd)
}

func TestEnd2End(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	log := newLogger(t)
	zap.ReplaceGlobals(zaptest.NewLogger(log))

	msg := message.Message{
		ID:      uuid.NewString(),
		Content: "content",
	}

	cmd.SetArgs([]string{"message", "create", "-i", msg.ID, "-c", msg.Content})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	cmd.SetArgs([]string{"message", "get", "-i", msg.ID})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	cmd.SetArgs([]string{"message", "delete", "-i", msg.ID})

	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(log.Messages), 3)

	bytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, strings.Contains(log.Messages[0], "successfully created message"))
	assert.True(t, strings.Contains(log.Messages[0], string(bytes)))

	assert.True(t, strings.Contains(log.Messages[1], "successfully got message"))
	assert.True(t, strings.Contains(log.Messages[1], string(bytes)))

	assert.True(t, strings.Contains(log.Messages[2], "successfully deleted message"))
	assert.True(t, strings.Contains(log.Messages[2], "null"))
}
