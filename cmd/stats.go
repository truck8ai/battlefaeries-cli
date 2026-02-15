package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	statsHP  int
	statsSTR int
	statsAGI int
	statsMAG int
)

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().IntVar(&statsHP, "hp", 0, "HP points to allocate")
	statsCmd.Flags().IntVar(&statsSTR, "str", 0, "Strength points to allocate")
	statsCmd.Flags().IntVar(&statsAGI, "agi", 0, "Agility points to allocate")
	statsCmd.Flags().IntVar(&statsMAG, "mag", 0, "Magic points to allocate")
}

var statsCmd = &cobra.Command{
	Use:   "stats <faerie-id>",
	Short: "Allocate stat points",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if statsHP+statsSTR+statsAGI+statsMAG == 0 {
			return fmt.Errorf("specify at least one stat: --hp, --str, --agi, --mag")
		}

		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"hp":       statsHP,
			"strength": statsSTR,
			"agility":  statsAGI,
			"magic":    statsMAG,
		}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post(fmt.Sprintf("/faeries/%s/stats", args[0]), body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Allocated struct {
				HP       int `json:"hp"`
				Strength int `json:"strength"`
				Agility  int `json:"agility"`
				Magic    int `json:"magic"`
			} `json:"allocated"`
		}
		json.Unmarshal(data, &resp)
		a := resp.Allocated
		color.Green("  Stats allocated: HP+%d STR+%d AGI+%d MAG+%d", a.HP, a.Strength, a.Agility, a.Magic)
		return nil
	},
}
