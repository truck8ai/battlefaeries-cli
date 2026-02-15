package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(replayCmd)
}

var replayCmd = &cobra.Command{
	Use:   "replay <battle-log-id>",
	Short: "View a battle replay",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/battle/%s", args[0]))
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			AttackerName string `json:"attackerName"`
			DefenderName string `json:"defenderName"`
			Winner       string `json:"winner"`
			Turns        []struct {
				Round      int    `json:"round"`
				Attacker   string `json:"attacker"`
				Defender   string `json:"defender"`
				Damage     int    `json:"damage"`
				IsCrit     bool   `json:"isCrit"`
				IsSkill    bool   `json:"isSkill"`
				SkillName  string `json:"skillName"`
				Healing    int    `json:"healing"`
			} `json:"turns"`
		}
		json.Unmarshal(data, &resp)

		fmt.Printf("  %s vs %s\n\n", resp.AttackerName, resp.DefenderName)

		currentRound := 0
		for _, t := range resp.Turns {
			if t.Round != currentRound {
				currentRound = t.Round
				fmt.Printf("  --- Round %d ---\n", currentRound)
			}

			action := fmt.Sprintf("%d dmg", t.Damage)
			if t.IsSkill && t.SkillName != "" {
				action = fmt.Sprintf("%s → %d dmg", t.SkillName, t.Damage)
			}
			if t.Healing > 0 {
				action += fmt.Sprintf(" +%d heal", t.Healing)
			}
			if t.IsCrit {
				action += " CRIT!"
			}
			fmt.Printf("  %s → %s: %s\n", t.Attacker, t.Defender, action)
		}

		fmt.Println()
		switch resp.Winner {
		case "attacker":
			color.Green("  Winner: %s", resp.AttackerName)
		case "draw":
			color.Yellow("  Draw")
		default:
			color.Red("  Winner: %s", resp.DefenderName)
		}
		return nil
	},
}
