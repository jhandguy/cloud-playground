package message

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

	msg := Message{
		ID:      uuid.NewString(),
		Content: "content",
	}

	buf := bytes.NewBufferString("")
	Cmd.SetOut(buf)
	Cmd.SetArgs([]string{"create", "-i", msg.ID, "-c", msg.Content})

	if err := Cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out, err := ioutil.ReadAll(buf)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(out), fmt.Sprintf("%v", msg)))

	buf.Reset()
	Cmd.SetArgs([]string{"get", "-i", msg.ID})

	if err := Cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out, err = ioutil.ReadAll(buf)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(out), fmt.Sprintf("%v", msg)))

	buf.Reset()
	Cmd.SetArgs([]string{"delete", "-i", msg.ID})

	if err := Cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	out, err = ioutil.ReadAll(buf)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(out), fmt.Sprintf("%v", nil)))
}
