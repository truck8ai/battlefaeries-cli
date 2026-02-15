package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/truck8ai/battlefaeries-cli/internal/client"
	"github.com/truck8ai/battlefaeries-cli/internal/format"
)

func init() {
	rootCmd.AddCommand(battleCmd)
	battleCmd.AddCommand(battleListCmd)
}

var battleCmd = &cobra.Command{
	Use:   "battle <defender-id>",
	Short: "Initiate a battle",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"defenderId": args[0],
		}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post("/battle", body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			BattleResult struct {
				Winner       string `json:"winner"`
				GoldChange   int    `json:"goldChange"`
				TrophyChange int    `json:"trophyChange"`
				WinStreak    int    `json:"winStreak"`
				BattleLogId  string `json:"battleLogId"`
			} `json:"battleResult"`
		}
		json.Unmarshal(data, &resp)
		r := resp.BattleResult

		switch r.Winner {
		case "attacker":
			color.Green("  Victory!")
		case "draw":
			color.Yellow("  Draw")
		default:
			color.Red("  Defeat")
		}

		fmt.Printf("  Gold: %+d  Trophies: %+d  Streak: %d\n", r.GoldChange, r.TrophyChange, r.WinStreak)
		fmt.Printf("  Replay: bf replay %s\n", r.BattleLogId)
		return nil
	},
}

var battleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available opponents",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/opponents")
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Opponents []client.Opponent `json:"opponents"`
		}
		json.Unmarshal(data, &resp)

		headers := []string{"ID", "Name", "Trophies", "Power", "Faeries", "Agent"}
		var rows [][]string
		for _, o := range resp.Opponents {
			agent := ""
			if o.IsAgentControlled {
				agent = "🤖"
			}
			rows = append(rows, []string{
				o.ID[:8], o.DisplayName,
				fmt.Sprintf("%d", o.Trophies), fmt.Sprintf("%d", o.TotalPower),
				fmt.Sprintf("%d", o.FaerieCount), agent,
			})
		}
		fmt.Print(format.Table(headers, rows))
		return nil
	},
}
