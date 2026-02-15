package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/truck8ai/battlefaeries-cli/internal/client"
	"github.com/truck8ai/battlefaeries-cli/internal/format"
)

func init() {
	rootCmd.AddCommand(teamCmd)
}

var teamCmd = &cobra.Command{
	Use:   "team [faerie-name]",
	Short: "Show your team",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/team")
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Team []client.Faerie `json:"team"`
		}
		json.Unmarshal(data, &resp)

		if len(args) > 0 {
			name := strings.ToLower(args[0])
			for _, f := range resp.Team {
				if strings.ToLower(f.Name) == name {
					printFaerieDetail(f)
					return nil
				}
			}
			return fmt.Errorf("faerie '%s' not found", args[0])
		}

		headers := []string{"Name", "Lv", "Element", "HP", "STR", "AGI", "MAG", "Pts", "Weapon", "Armor", "Accessory"}
		var rows [][]string
		for _, f := range resp.Team {
			rows = append(rows, []string{
				f.Name, fmt.Sprintf("%d", f.Level), f.Element,
				fmt.Sprintf("%d", f.HP), fmt.Sprintf("%d", f.Strength),
				fmt.Sprintf("%d", f.Agility), fmt.Sprintf("%d", f.Magic),
				fmt.Sprintf("%d", f.UnallocatedPoints),
				strOrDash(f.WeaponName), strOrDash(f.ArmorName), strOrDash(f.AccessoryName),
			})
		}

		fmt.Print(format.Table(headers, rows))
		return nil
	},
}

func printFaerieDetail(f client.Faerie) {
	bold := color.New(color.Bold)
	bold.Printf("  %s", f.Name)
	fmt.Printf("  [%s]  Lv.%d\n\n", f.Element, f.Level)
	fmt.Printf("  HP: %d  STR: %d  AGI: %d  MAG: %d\n", f.HP, f.Strength, f.Agility, f.Magic)
	fmt.Printf("  XP: %d  Unallocated: %d\n", f.TotalXP, f.UnallocatedPoints)

	if len(f.Skills) > 0 {
		fmt.Println("\n  Skills:")
		for _, s := range f.Skills {
			fmt.Printf("    [%d] %s (%s/%s) Power:%d\n", s.SkillSlot, s.Name, s.SkillType, s.Element, s.Power)
		}
	}
}

func strOrDash(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}
