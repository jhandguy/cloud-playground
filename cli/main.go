package main

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/jhandguy/devops-playground/cli/load"
	"github.com/jhandguy/devops-playground/cli/message"
)

var cmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI",
	Long:  "CLI that interacts with devops-playground system",
}

func setupLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	zap.ReplaceGlobals(logger)
}

func handleUnbindableFlag(err error) {
	if err != nil {
		zap.S().Fatalw("could not bind flag", "error", err)
	}
}

func init() {
	setupLogger()

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	cmd.AddCommand(load.Cmd)
	cmd.AddCommand(message.Cmd)

	cmd.PersistentFlags().StringP("token", "t", "", "gateway auth token")
	handleUnbindableFlag(viper.BindPFlag("gateway-token", cmd.PersistentFlags().Lookup("token")))

	cmd.PersistentFlags().StringP("url", "u", "", "gateway URL")
	handleUnbindableFlag(viper.BindPFlag("gateway-url", cmd.PersistentFlags().Lookup("url")))

	cmd.PersistentFlags().StringP("host", "o", "", "gateway host")
	handleUnbindableFlag(viper.BindPFlag("gateway-host", cmd.PersistentFlags().Lookup("host")))

	cmd.PersistentFlags().StringP("canary", "a", "", "canary header")
	handleUnbindableFlag(viper.BindPFlag("gateway-canary", cmd.PersistentFlags().Lookup("canary")))
}

func main() {
	if err := cmd.Execute(); err != nil {
		zap.S().Errorw("error executing command", "cmd", cmd.Short, "error", err)
	}
}
