package message

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	id, content   string
	client        *resty.Client
	debug, canary string
)

func handleMissingFlag(err error) {
	if err != nil {
		log.Fatalf("missing required flag: %v", err)
	}
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	Cmd.AddCommand(createMessageCmd)
	Cmd.AddCommand(getMessageCmd)
	Cmd.AddCommand(deleteMessageCmd)

	Cmd.PersistentFlags().StringP("token", "t", "", "gateway auth token")
	handleMissingFlag(viper.BindPFlag("gateway-token", Cmd.PersistentFlags().Lookup("token")))

	Cmd.PersistentFlags().StringP("url", "u", "", "gateway URL")
	handleMissingFlag(viper.BindPFlag("gateway-url", Cmd.PersistentFlags().Lookup("url")))

	Cmd.PersistentFlags().StringP("host", "o", "", "gateway host")
	handleMissingFlag(viper.BindPFlag("gateway-host", Cmd.PersistentFlags().Lookup("host")))

	Cmd.PersistentFlags().StringVarP(&debug, "debug", "d", "", "debug header")
	Cmd.PersistentFlags().StringVarP(&canary, "canary", "a", "", "canary header")

	createMessageCmd.Flags().StringVarP(&id, "id", "i", "", "id of the message")
	createMessageCmd.Flags().StringVarP(&content, "content", "c", "", "content of the message")
	handleMissingFlag(createMessageCmd.MarkFlagRequired("content"))

	getMessageCmd.Flags().StringVarP(&id, "id", "i", "", "id of the message")
	handleMissingFlag(getMessageCmd.MarkFlagRequired("id"))

	deleteMessageCmd.Flags().StringVarP(&id, "id", "i", "", "id of the message")
	handleMissingFlag(deleteMessageCmd.MarkFlagRequired("id"))

	url := fmt.Sprintf("http://%s", viper.GetString("gateway-url"))
	token := viper.GetString("gateway-token")
	host := viper.GetString("gateway-host")

	client = resty.
		New().
		SetBaseURL(url).
		SetAuthToken(token).
		SetHeader("Host", host).
		SetHeader("x-debug", debug).
		SetHeader("x-canary", canary)
}

type Message struct {
	ID      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
}

func Create(message Message) (*resty.Response, error) {
	res, err := client.
		R().
		SetResult(Message{}).
		SetBody(message).
		Post("/message")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Get(id string) (*resty.Response, error) {
	res, err := client.
		R().
		SetResult(Message{}).
		SetPathParam("id", id).
		Get("/message/{id}")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Delete(id string) (*resty.Response, error) {
	res, err := client.
		R().
		SetPathParam("id", id).
		Delete("/message/{id}")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func createMessage(cmd *cobra.Command, _ []string) {
	res, err := Create(Message{
		ID:      id,
		Content: content,
	})
	if err != nil {
		log.Fatalf("failed to create message: %v", err)
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "successfully created message: %v\n", res.Result()); err != nil {
		log.Fatal(err)
	}
}

func getMessage(cmd *cobra.Command, _ []string) {
	res, err := Get(id)
	if err != nil {
		log.Fatalf("failed to get message: %v", err)
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "successfully got message: %v\n", res.Result()); err != nil {
		log.Fatal(err)
	}
}

func deleteMessage(cmd *cobra.Command, _ []string) {
	res, err := Delete(id)
	if err != nil {
		log.Fatalf("failed to delete message: %v", err)
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "successfully deleted message: %v\n", res.Result()); err != nil {
		log.Fatal(err)
	}
}
