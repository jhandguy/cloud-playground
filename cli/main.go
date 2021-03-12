package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/jhandguy/devops-playground/cli/load"
	"github.com/jhandguy/devops-playground/cli/message"
)

var cmd = &cobra.Command{
	Use:   "cli",
	Short: "CLI",
	Long:  "CLI that interacts with devops-playground system",
}

func init() {
	cmd.AddCommand(load.Cmd)
	cmd.AddCommand(message.Cmd)
}

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("error executing %s: %v", cmd.Short, err)
	}
}
