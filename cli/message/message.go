package message

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Cmd message command
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
	id, content string
)

func handleMissingFlag(err error) {
	if err != nil {
		zap.S().Fatalw("missing required flag", "error", err.Error())
	}
}

func init() {
	Cmd.AddCommand(createMessageCmd)
	Cmd.AddCommand(getMessageCmd)
	Cmd.AddCommand(deleteMessageCmd)

	createMessageCmd.Flags().StringVarP(&id, "id", "i", "", "id of the message")
	createMessageCmd.Flags().StringVarP(&content, "content", "c", "", "content of the message")
	handleMissingFlag(createMessageCmd.MarkFlagRequired("content"))

	getMessageCmd.Flags().StringVarP(&id, "id", "i", "", "id of the message")
	handleMissingFlag(getMessageCmd.MarkFlagRequired("id"))

	deleteMessageCmd.Flags().StringVarP(&id, "id", "i", "", "id of the message")
	handleMissingFlag(deleteMessageCmd.MarkFlagRequired("id"))
}

func newClient() *resty.Client {
	url := fmt.Sprintf("http://%s", viper.GetString("gateway-url"))
	token := viper.GetString("gateway-token")
	host := viper.GetString("gateway-host")
	canary := viper.GetString("gateway-canary")

	return resty.
		New().
		SetBaseURL(url).
		SetAuthToken(token).
		SetHeader("Host", host).
		SetHeader("X-Canary", canary)
}

type Message struct {
	ID      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
}

func Create(message Message) (*resty.Response, error) {
	res, err := newClient().
		R().
		SetResult(Message{}).
		SetBody(message).
		Post("/message")
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to create message: %d", res.StatusCode())
	}

	return res, nil
}

func Get(id string) (*resty.Response, error) {
	res, err := newClient().
		R().
		SetResult(Message{}).
		SetPathParam("id", id).
		Get("/message/{id}")
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to get message: %d", res.StatusCode())
	}

	return res, nil
}

func Delete(id string) (*resty.Response, error) {
	res, err := newClient().
		R().
		SetPathParam("id", id).
		Delete("/message/{id}")
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		return nil, fmt.Errorf("failed to delete message: %d", res.StatusCode())
	}

	return res, nil
}

func createMessage(cmd *cobra.Command, _ []string) {
	res, err := Create(Message{
		ID:      id,
		Content: content,
	})
	if err != nil {
		zap.S().Errorw("failed to create message", "error", err.Error())
		return
	}

	zap.S().Infow("successfully created message", "message", res.Result())
}

func getMessage(cmd *cobra.Command, _ []string) {
	res, err := Get(id)
	if err != nil {
		zap.S().Errorw("failed to get message", "error", err.Error())
		return
	}

	zap.S().Infow("successfully got message", "message", res.Result())
}

func deleteMessage(cmd *cobra.Command, _ []string) {
	res, err := Delete(id)
	if err != nil {
		zap.S().Errorw("failed to delete message", "error", err.Error())
		return
	}

	zap.S().Infow("successfully deleted message", "message", res.Result())
}
