package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/truck8ai/battlefaeries-cli/internal/client"
	"github.com/truck8ai/battlefaeries-cli/internal/config"
)

func init() {
	rootCmd.AddCommand(loginCmd)
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with your API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if apiURL != "" {
			cfg.APIURL = apiURL
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your API key: ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)

		if !strings.HasPrefix(key, "bf_live_") {
			return fmt.Errorf("invalid API key format (must start with bf_live_)")
		}

		// Validate by hitting the status endpoint
		c := client.NewWithKey(cfg.APIURL, key)
		_, err = c.Get("/status")
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		cfg.APIKey = key
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		color.Green("✓ Logged in successfully!")
		return nil
	},
}
