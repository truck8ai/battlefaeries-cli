package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/truck8ai/battlefaeries-cli/internal/client"
	"github.com/truck8ai/battlefaeries-cli/internal/format"
)

var historyOutcome string

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().StringVar(&historyOutcome, "outcome", "", "Filter by outcome: win, loss, draw")
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View battle history",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		path := "/history"
		if historyOutcome != "" {
			path += "?outcome=" + historyOutcome
		}

		data, err := c.Get(path)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Battles []client.BattleHistory `json:"battles"`
		}
		json.Unmarshal(data, &resp)

		if len(resp.Battles) == 0 {
			fmt.Println("  No battles found.")
			return nil
		}

		headers := []string{"Log ID", "Attacker", "Defender", "Gold", "Trophies", "Date"}
		var rows [][]string
		for _, b := range resp.Battles {
			rows = append(rows, []string{
				b.BattleLogID[:8],
				b.AttackerName, b.DefenderName,
				fmt.Sprintf("%+d", b.AttackerGoldChange),
				fmt.Sprintf("%+d", b.AttackerTrophyChange),
				b.CreatedAt[:10],
			})
		}
		fmt.Print(format.Table(headers, rows))
		return nil
	},
}
