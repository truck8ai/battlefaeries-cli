package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(displayNameCmd)
}

var displayNameCmd = &cobra.Command{
	Use:   "display-name <new-name>",
	Short: "Change your player display name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"displayName": args[0],
		}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post("/display-name", body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		color.Green("  Display name changed to %s!", args[0])
		return nil
	},
}
