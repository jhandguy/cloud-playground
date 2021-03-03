package message

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	id, content string
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("gateway")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	Cmd.AddCommand(createMessageCmd)
	Cmd.AddCommand(getMessageCmd)
	Cmd.AddCommand(deleteMessageCmd)

	Cmd.PersistentFlags().StringP("token", "t", "", "gateway auth token")
	handleMissingFlag(viper.BindPFlag("api-key", Cmd.PersistentFlags().Lookup("token")))

	Cmd.PersistentFlags().StringP("url", "u", "", "gateway URL")
	handleMissingFlag(viper.BindPFlag("url", Cmd.PersistentFlags().Lookup("url")))

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
	url := fmt.Sprintf("http://%s", viper.GetString("url"))
	token := viper.GetString("api-key")

	return resty.
		New().
		SetHostURL(url).
		SetAuthToken(token)
}

func createMessage(cmd *cobra.Command, _ []string) {
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

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "successfully created message: %v\n", res.Result()); err != nil {
		log.Fatal(err)
	}
}

func getMessage(cmd *cobra.Command, _ []string) {
	res, err := newClient().
		R().
		SetResult(message{}).
		SetPathParam("id", id).
		Get("/message/{id}")
	if err != nil {
		log.Fatalf("failed to get message: %v", err)
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "successfully got message: %v\n", res.Result()); err != nil {
		log.Fatal(err)
	}
}

func deleteMessage(cmd *cobra.Command, _ []string) {
	res, err := newClient().
		R().
		SetPathParam("id", id).
		Delete("/message/{id}")
	if err != nil {
		log.Fatalf("failed to delete message: %v", err)
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "successfully deleted message: %v\n", res.Result()); err != nil {
		log.Fatal(err)
	}
}
