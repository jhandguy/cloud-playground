package message

import (
	"log"
	"os"
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

	err := Cmd.PersistentFlags().Set("token", retrieveEnv("GATEWAY_API_KEY"))
	assert.Nil(t, err)

	err = Cmd.PersistentFlags().Set("url", retrieveEnv("GATEWAY_URL"))
	assert.Nil(t, err)

	id := uuid.NewString()

	err = createMessageCmd.Flags().Set("id", id)
	assert.Nil(t, err)

	err = createMessageCmd.Flags().Set("content", "content")
	assert.Nil(t, err)

	err = createMessageCmd.Execute()
	assert.Nil(t, err)

	err = getMessageCmd.Flags().Set("id", id)
	assert.Nil(t, err)

	err = getMessageCmd.Execute()
	assert.Nil(t, err)

	err = deleteMessageCmd.Flags().Set("id", id)
	assert.Nil(t, err)

	err = deleteMessageCmd.Execute()
	assert.Nil(t, err)
}

func retrieveEnv(key string) string {
	env, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("could not lookup env %s", key)
	}
	return env
}
