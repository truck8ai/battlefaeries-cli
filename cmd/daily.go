package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(dailyCmd)
}

var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Claim your daily reward",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post("/daily", body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Gold   int `json:"gold"`
			Streak int `json:"streak"`
		}
		json.Unmarshal(data, &resp)

		color.Green("  Daily reward claimed!")
		fmt.Printf("  Gold earned: %d  Streak: %d\n", resp.Gold, resp.Streak)
		return nil
	},
}
