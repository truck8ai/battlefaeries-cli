package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(skillsCmd)
	skillsCmd.AddCommand(skillsBuyCmd)
	skillsCmd.AddCommand(skillsAssignCmd)
	skillsCmd.AddCommand(skillsUnassignCmd)
}

var skillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "View skills",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/skills")
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		fmt.Println("  Use --json flag for full skill data")
		return nil
	},
}

var skillsBuyCmd = &cobra.Command{
	Use:   "buy <skill-id>",
	Short: "Purchase a skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{"skillId": args[0]}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post("/skills/purchase", body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Skill   map[string]interface{} `json:"skill"`
			NewGold int                    `json:"newGold"`
		}
		json.Unmarshal(data, &resp)
		color.Green("  Purchased %s!", resp.Skill["name"])
		fmt.Printf("  Remaining gold: %d\n", resp.NewGold)
		return nil
	},
}

var skillsAssignCmd = &cobra.Command{
	Use:   "assign <faerie-id> <skill-id> <slot>",
	Short: "Assign a skill to a faerie",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"faerieId":      args[0],
			"playerSkillId": args[1],
			"slot":          args[2],
			"action":        "assign",
		}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post("/skills/assign", body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		color.Green("  Skill assigned!")
		return nil
	},
}

var skillsUnassignCmd = &cobra.Command{
	Use:   "unassign <skill-id>",
	Short: "Unassign a skill",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{
			"playerSkillId": args[0],
			"action":        "unassign",
		}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post("/skills/assign", body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		color.Green("  Skill unassigned!")
		return nil
	},
}
