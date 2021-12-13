package main

import (
	"log"

	"github.com/spf13/cobra"
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

func init() {
	cmd.AddCommand(load.Cmd)
	cmd.AddCommand(message.Cmd)
}

func main() {
	setupLogger()

	if err := cmd.Execute(); err != nil {
		zap.S().Errorw("error executing command", "cmd", cmd.Short, "error", err)
	}
}
