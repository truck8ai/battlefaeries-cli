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
	rootCmd.AddCommand(tournamentsCmd)
	tournamentsCmd.AddCommand(tournamentJoinCmd)
	tournamentsCmd.AddCommand(tournamentLeaveCmd)
	tournamentsCmd.AddCommand(tournamentInfoCmd)
}

var tournamentsCmd = &cobra.Command{
	Use:   "tournaments",
	Short: "List available tournaments",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/tournaments")
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Tournaments []client.Tournament `json:"tournaments"`
		}
		json.Unmarshal(data, &resp)

		if len(resp.Tournaments) == 0 {
			fmt.Println("  No tournaments available.")
			return nil
		}

		headers := []string{"ID", "Name", "Status", "Max", "Prize", "Type", "Players", "Joined"}
		var rows [][]string
		for _, t := range resp.Tournaments {
			joined := ""
			if t.IsRegistered {
				joined = "✓"
			}
			rows = append(rows, []string{
				t.ID[:8], t.Name, t.Status,
				fmt.Sprintf("%d", t.MaxPlayers),
				fmt.Sprintf("%d", t.PrizeTrophies),
				t.ParticipantType,
				fmt.Sprintf("%d", t.ParticipantCount),
				joined,
			})
		}
		fmt.Print(format.Table(headers, rows))
		return nil
	},
}

var tournamentJoinCmd = &cobra.Command{
	Use:   "join <tournament-id>",
	Short: "Register for a tournament",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Post(fmt.Sprintf("/tournaments/%s/register", args[0]), body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		color.Green("  Registered for tournament!")
		return nil
	},
}

var tournamentLeaveCmd = &cobra.Command{
	Use:   "leave <tournament-id>",
	Short: "Withdraw from a tournament",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body := map[string]interface{}{}
		if reason != "" {
			body["reasoning"] = reason
		}

		data, err := c.Delete(fmt.Sprintf("/tournaments/%s/register", args[0]), body)
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		color.Green("  Withdrew from tournament!")
		return nil
	},
}

var tournamentInfoCmd = &cobra.Command{
	Use:   "info <tournament-id>",
	Short: "View tournament details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get(fmt.Sprintf("/tournaments/%s", args[0]))
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		fmt.Println(string(data))
		return nil
	},
}
