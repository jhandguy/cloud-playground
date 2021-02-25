package message

import (
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "message",
	Short: "Message sub commands",
	Long:  "Message sub commands that interacts with gateway service to manage objects with S3 and manage items with Dynamo",
}

var createMessageCmd = &cobra.Command{
	Use:   "create",
	Short: "Create message",
	Long:  "Create message with id and content",
	Run:   createMessage,
}

var getMessageCmd = &cobra.Command{
	Use:   "get",
	Short: "Get message",
	Long:  "Get message by id",
	Run:   getMessage,
}

var deleteMessageCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete message",
	Long:  "Delete message by id",
	Run:   deleteMessage,
}

var (
	token, url  string
	id, content string
)

func init() {
	Cmd.AddCommand(createMessageCmd)
	Cmd.AddCommand(getMessageCmd)
	Cmd.AddCommand(deleteMessageCmd)

	Cmd.PersistentFlags().StringVarP(&token, "token", "t", "", "gateway auth token")
	Cmd.PersistentFlags().StringVarP(&url, "url", "u", "", "gateway URL")
	handleMissingFlag(Cmd.MarkPersistentFlagRequired("token"))
	handleMissingFlag(Cmd.MarkPersistentFlagRequired("url"))

	createMessageCmd.Flags().StringVarP(&id, "id", "i", "", "id of the message")
	createMessageCmd.Flags().StringVarP(&content, "content", "c", "", "content of the message")
	handleMissingFlag(createMessageCmd.MarkFlagRequired("content"))

	getMessageCmd.Flags().StringVarP(&id, "id", "i", "", "id of the message")
	handleMissingFlag(getMessageCmd.MarkFlagRequired("id"))

	deleteMessageCmd.Flags().StringVarP(&id, "id", "i", "", "id of the message")
	handleMissingFlag(deleteMessageCmd.MarkFlagRequired("id"))
}

func handleMissingFlag(err error) {
	if err != nil {
		log.Fatalf("missing required flag: %v", err)
	}
}

type message struct {
	ID      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
}

func newClient() *resty.Client {
	return resty.
		New().
		SetHostURL(url).
		SetAuthToken(token).
		SetDebug(true)
}

func createMessage(*cobra.Command, []string) {
	res, err := newClient().
		R().
		SetResult(message{}).
		SetBody(message{
			ID:      id,
			Content: content,
		}).
		Post("/message")
	if err != nil {
		log.Fatalf("failed to create message: %v", err)
	}

	log.Printf("successfully created message: %v", res.Result().(*message))
}

func getMessage(*cobra.Command, []string) {
	res, err := newClient().
		R().
		SetResult(message{}).
		SetPathParam("id", id).
		Get("/message/{id}")
	if err != nil {
		log.Fatalf("failed to get message: %v", err)
	}

	log.Printf("successfully got message: %v", res.Result().(*message))
}

func deleteMessage(*cobra.Command, []string) {
	_, err := newClient().
		R().
		SetPathParam("id", id).
		Delete("/message/{id}")
	if err != nil {
		log.Fatalf("failed to delete message: %v", err)
	}

	log.Print("successfully deleted message")
}
