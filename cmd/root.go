package cmd

import (
	"github.com/spf13/cobra"
	"github.com/truck8ai/battlefaeries-cli/internal/client"
)

var (
	jsonOutput bool
	reason     string
	apiURL     string
	logFlag    bool
)

var rootCmd = &cobra.Command{
	Use:   "bf",
	Short: "Battle Faeries CLI — play the game from your terminal",
	Long:  "Battle Faeries Agent CLI. Connect any AI agent or play manually via the command line.",
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output raw JSON")
	rootCmd.PersistentFlags().StringVar(&reason, "reason", "", "Attach reasoning to the action")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "Override API URL")
	rootCmd.PersistentFlags().BoolVar(&logFlag, "log", false, "Log all API requests/responses to ~/.battlefaeries/logs/activity.jsonl")
}

// newClient creates an API client with global flags applied.
func newClient() (*client.Client, error) {
	c, err := client.New()
	if err != nil {
		return nil, err
	}
	if logFlag {
		c.SetLogEnabled(true)
	}
	return c, nil
}

func Execute() error {
	return rootCmd.Execute()
}
