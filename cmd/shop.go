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
	rootCmd.AddCommand(shopCmd)
	shopCmd.AddCommand(shopBuyCmd)
}

var shopCmd = &cobra.Command{
	Use:   "shop",
	Short: "Browse the shop",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/shop")
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Items []client.ShopItem `json:"items"`
			Gold  int               `json:"gold"`
		}
		json.Unmarshal(data, &resp)

		fmt.Printf("  Your gold: %d\n\n", resp.Gold)

		headers := []string{"ID", "Name", "Type", "Tier", "Price", "ATK", "DEF", "HP", "SPD", "CRIT"}
		var rows [][]string
		for _, item := range resp.Items {
			rows = append(rows, []string{
				item.ID, item.Name, item.EquipmentType,
				fmt.Sprintf("T%d", item.Tier), fmt.Sprintf("%d", item.Price),
				fmt.Sprintf("%d", item.AttackBonus), fmt.Sprintf("%d", item.DefenseBonus),
				fmt.Sprintf("%d", item.HPBonus), fmt.Sprintf("%d", item.SpeedBonus),
				fmt.Sprintf("%d", item.CritBonus),
			})
		}
		fmt.Print(format.Table(headers, rows))
		return nil
	},
}

var shopBuyCmd = &cobra.Command{
	Use:   "buy <item-id>",
	Short: "Purchase an item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"itemId": args[0],
		}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post("/shop/purchase", body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Item    client.ShopItem `json:"item"`
			NewGold int             `json:"newGold"`
		}
		json.Unmarshal(data, &resp)

		color.Green("  Purchased %s!", resp.Item.Name)
		fmt.Printf("  Remaining gold: %d\n", resp.NewGold)
		return nil
	},
}
