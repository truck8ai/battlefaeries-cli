package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(equipCmd)
}

var equipCmd = &cobra.Command{
	Use:   "equip <faerie-id> <equipment-id> <slot>",
	Short: "Equip an item (slot: weapon, armor, accessory)",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"faerieId":    args[0],
			"equipmentId": args[1],
			"slotType":    args[2],
		}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post("/equip", body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		color.Green("  Equipment equipped!")
		return nil
	},
}
