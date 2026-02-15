package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/truck8ai/battlefaeries-cli/internal/client"
	"github.com/truck8ai/battlefaeries-cli/internal/format"
)

func init() {
	rootCmd.AddCommand(inventoryCmd)
}

var inventoryCmd = &cobra.Command{
	Use:   "inventory",
	Short: "View your equipment inventory",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/inventory")
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Equipment []client.Equipment `json:"equipment"`
		}
		json.Unmarshal(data, &resp)

		if len(resp.Equipment) == 0 {
			fmt.Println("  No equipment owned.")
			return nil
		}

		headers := []string{"ID", "Name", "Type", "Tier", "ATK", "DEF", "HP", "SPD", "CRIT", "Equipped"}
		var rows [][]string
		for _, e := range resp.Equipment {
			equipped := ""
			if e.EquippedOn != nil {
				equipped = (*e.EquippedOn)[:8]
			}
			rows = append(rows, []string{
				e.ID[:8], e.Name, e.EquipmentType,
				fmt.Sprintf("T%d", e.Tier),
				fmt.Sprintf("%d", e.AttackBonus), fmt.Sprintf("%d", e.DefenseBonus),
				fmt.Sprintf("%d", e.HPBonus), fmt.Sprintf("%d", e.SpeedBonus),
				fmt.Sprintf("%d", e.CritBonus), equipped,
			})
		}
		fmt.Print(format.Table(headers, rows))
		return nil
	},
}
