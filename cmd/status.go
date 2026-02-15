package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/truck8ai/battlefaeries-cli/internal/client"
)

func init() {
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(whoamiCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show your player status",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/status")
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Player client.PlayerStatus `json:"player"`
		}
		json.Unmarshal(data, &resp)
		p := resp.Player

		bold := color.New(color.Bold)
		bold.Printf("  %s\n\n", p.DisplayName)
		fmt.Printf("  Gold:      %d\n", p.Gold)
		fmt.Printf("  Stamina:   %d\n", p.Stamina)
		fmt.Printf("  Trophies:  %d\n", p.Trophies)
		fmt.Printf("  Power:     %d (combat: %d)\n", p.TotalPower, p.CombatPower)
		fmt.Printf("  Faeries:   %d\n", p.FaerieCount)
		fmt.Printf("  Win Streak: %d (best: %d)\n", p.WinStreak, p.BestWinStreak)

		if p.DailyRewardAvailable {
			color.Yellow("\n  Daily reward available! Run: bf daily")
		}

		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show your player name",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/status")
		if err != nil {
			return err
		}

		var resp struct {
			Player client.PlayerStatus `json:"player"`
		}
		json.Unmarshal(data, &resp)
		fmt.Println(resp.Player.DisplayName)
		return nil
	},
}
