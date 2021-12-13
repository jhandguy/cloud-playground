package message

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
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
	assert.Equal(t, len(Cmd.Commands()), 3)
	assert.Equal(t, Cmd.Commands()[0], createMessageCmd)
	assert.Equal(t, Cmd.Commands()[1], deleteMessageCmd)
	assert.Equal(t, Cmd.Commands()[2], getMessageCmd)
}

func TestEnd2End(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	log := newLogger(t)
	zap.ReplaceGlobals(zaptest.NewLogger(log))

	msg := Message{
		ID:      uuid.NewString(),
		Content: "content",
	}

	Cmd.SetArgs([]string{"create", "-i", msg.ID, "-c", msg.Content})

	if err := Cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	Cmd.SetArgs([]string{"get", "-i", msg.ID})

	if err := Cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	Cmd.SetArgs([]string{"delete", "-i", msg.ID})

	if err := Cmd.Execute(); err != nil {
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
