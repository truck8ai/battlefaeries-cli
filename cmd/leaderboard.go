package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/truck8ai/battlefaeries-cli/internal/format"
)

var leaderboardType string

func init() {
	leaderboardCmd.Flags().StringVarP(&leaderboardType, "type", "t", "combined", "Leaderboard type: combined, power, trophies")
	rootCmd.AddCommand(leaderboardCmd)
}

var leaderboardCmd = &cobra.Command{
	Use:   "leaderboard",
	Short: "View the leaderboard",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		endpoint := "/leaderboard"
		if leaderboardType != "" && leaderboardType != "combined" {
			endpoint = fmt.Sprintf("/leaderboard?type=%s", leaderboardType)
		}

		data, err := c.Get(endpoint)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Type    string `json:"type"`
			Players []struct {
				Rank              int    `json:"rank"`
				DisplayName       string `json:"display_name"`
				Trophies          int    `json:"trophies"`
				TotalPower        int    `json:"total_power"`
				CombatPower       int    `json:"combat_power"`
				WinStreak         int    `json:"win_streak"`
				IsAgentControlled bool   `json:"is_agent_controlled"`
			} `json:"players"`
			YourRank int `json:"yourRank"`
		}
		json.Unmarshal(data, &resp)

		showCombatPower := resp.Type == "power"

		var headers []string
		if showCombatPower {
			headers = []string{"#", "Name", "Trophies", "CombatPower", "Streak", "Agent"}
		} else {
			headers = []string{"#", "Name", "Trophies", "Power", "Streak", "Agent"}
		}

		var rows [][]string
		for _, p := range resp.Players {
			agent := ""
			if p.IsAgentControlled {
				agent = "🤖"
			}
			power := p.TotalPower
			if showCombatPower {
				power = p.CombatPower
			}
			rows = append(rows, []string{
				fmt.Sprintf("%d", p.Rank), p.DisplayName,
				fmt.Sprintf("%d", p.Trophies), fmt.Sprintf("%d", power),
				fmt.Sprintf("%d", p.WinStreak), agent,
			})
		}
		fmt.Print(format.Table(headers, rows))

		if resp.YourRank > 0 {
			fmt.Printf("\n  Your rank: #%d\n", resp.YourRank)
		}
		return nil
	},
}
