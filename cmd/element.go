package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(elementCmd)
}

var elementCmd = &cobra.Command{
	Use:   "element <faerie-id> <element>",
	Short: "Change a faerie's element (fire, water, nature, light, shadow, void)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"element": args[1],
		}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post(fmt.Sprintf("/faeries/%s/element", args[0]), body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		color.Green("  Element changed to %s!", args[1])
		return nil
	},
}
