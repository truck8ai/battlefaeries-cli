package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(apiRefCmd)
}

var apiRefCmd = &cobra.Command{
	Use:   "api-ref",
	Short: "Show API reference — all endpoints, auth, and CLI flags",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/api-ref")
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			ApiRef json.RawMessage `json:"apiRef"`
		}
		json.Unmarshal(data, &resp)

		var ref struct {
			Authentication struct {
				Method  string `json:"method"`
				Example string `json:"example"`
				Notes   string `json:"notes"`
			} `json:"authentication"`
			Endpoints []struct {
				Category string `json:"category"`
				Routes   []struct {
					Method      string `json:"method"`
					Path        string `json:"path"`
					Description string `json:"description"`
					Body        string `json:"body"`
				} `json:"routes"`
			} `json:"endpoints"`
			Cli struct {
				Install string            `json:"install"`
				Login   string            `json:"login"`
				Flags   map[string]string `json:"flags"`
			} `json:"cli"`
			Notes []string `json:"notes"`
		}
		json.Unmarshal(resp.ApiRef, &ref)

		bold := color.New(color.Bold)
		cyan := color.New(color.FgCyan)
		green := color.New(color.FgGreen)
		blue := color.New(color.FgBlue)
		red := color.New(color.FgRed)

		bold.Println("  Battle Faeries — API Reference")
		fmt.Println(strings.Repeat("─", 50))
		fmt.Println()

		// Auth
		cyan.Println("  Authentication")
		fmt.Printf("    %s\n", ref.Authentication.Method)
		fmt.Printf("    Example: %s\n", ref.Authentication.Example)
		fmt.Println()

		// Endpoints
		for _, group := range ref.Endpoints {
			cyan.Printf("  %s\n", group.Category)
			for _, r := range group.Routes {
				methodColor := green
				if r.Method == "GET" {
					methodColor = blue
				} else if r.Method == "DELETE" {
					methodColor = red
				}
				methodColor.Printf("    %-6s", r.Method)
				fmt.Printf(" %-40s %s\n", r.Path, r.Description)
				if r.Body != "" {
					fmt.Printf("           Body: %s\n", r.Body)
				}
			}
			fmt.Println()
		}

		// CLI flags
		cyan.Println("  Global CLI Flags")
		for flag, desc := range ref.Cli.Flags {
			green.Printf("    %-20s", flag)
			fmt.Printf(" %s\n", desc)
		}
		fmt.Println()

		// Notes
		cyan.Println("  Notes")
		for _, note := range ref.Notes {
			fmt.Printf("    • %s\n", note)
		}
		fmt.Println()

		return nil
	},
}
