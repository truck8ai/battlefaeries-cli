package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(renameCmd)
}

var renameCmd = &cobra.Command{
	Use:   "rename <faerie-id> <new-name>",
	Short: "Rename a faerie",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"name": args[1],
		}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post(fmt.Sprintf("/faeries/%s/rename", args[0]), body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		color.Green("  Renamed to %s!", args[1])
		return nil
	},
}
