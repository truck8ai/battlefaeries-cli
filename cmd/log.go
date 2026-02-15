package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/truck8ai/battlefaeries-cli/internal/client"
	"github.com/truck8ai/battlefaeries-cli/internal/config"
)

func init() {
	rootCmd.AddCommand(logCmd)
	logCmd.AddCommand(logEnableCmd)
	logCmd.AddCommand(logDisableCmd)
	logCmd.AddCommand(logClearCmd)
	logCmd.AddCommand(logPathCmd)
	logCmd.AddCommand(logStatsCmd)

	logCmd.Flags().IntVarP(&logTailN, "tail", "n", 20, "Number of recent entries to show")
	logCmd.Flags().StringVar(&logFilterMethod, "method", "", "Filter by HTTP method (GET, POST, DELETE)")
	logCmd.Flags().StringVar(&logFilterPath, "path", "", "Filter by path substring")
	logCmd.Flags().BoolVar(&logFilterErrors, "errors", false, "Show only errors")
}

var (
	logTailN       int
	logFilterMethod string
	logFilterPath  string
	logFilterErrors bool
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View local activity log",
	Long: `View, manage, and analyze local API activity logs.

Logs are stored at ~/.battlefaeries/logs/activity.jsonl when enabled.
Enable with 'bf log enable' or use the --log flag on any command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logPath := config.LogPath()

		f, err := os.Open(logPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No log file found. Enable logging with: bf log enable")
				return nil
			}
			return err
		}
		defer f.Close()

		var entries []client.LogEntry
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB line buffer
		for scanner.Scan() {
			var entry client.LogEntry
			if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
				continue
			}

			// Apply filters
			if logFilterMethod != "" && !strings.EqualFold(entry.Method, logFilterMethod) {
				continue
			}
			if logFilterPath != "" && !strings.Contains(entry.Path, logFilterPath) {
				continue
			}
			if logFilterErrors && entry.Success {
				continue
			}

			entries = append(entries, entry)
		}

		if len(entries) == 0 {
			fmt.Println("No matching log entries.")
			return nil
		}

		// Show last N entries
		start := 0
		if len(entries) > logTailN {
			start = len(entries) - logTailN
		}

		methodColor := map[string]*color.Color{
			"GET":    color.New(color.FgCyan),
			"POST":   color.New(color.FgGreen),
			"DELETE": color.New(color.FgRed),
		}

		for _, entry := range entries[start:] {
			ts := entry.Timestamp
			if len(ts) > 19 {
				ts = ts[:19] // trim to seconds
			}

			mc := methodColor[entry.Method]
			if mc == nil {
				mc = color.New(color.FgWhite)
			}

			status := color.GreenString("OK")
			if !entry.Success {
				status = color.RedString("FAIL")
			}

			fmt.Printf("%s  %s  %-6s %-40s  %4dms  %s",
				color.HiBlackString(ts),
				status,
				mc.Sprint(entry.Method),
				entry.Path,
				entry.DurationMs,
				"",
			)
			if entry.Error != "" {
				fmt.Printf(" %s", color.RedString(entry.Error))
			}
			fmt.Println()
		}

		fmt.Printf("\n%s Showing %d of %d entries. Log: %s\n",
			color.HiBlackString("—"),
			min(logTailN, len(entries)),
			len(entries),
			logPath,
		)

		return nil
	},
}

var logEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable persistent activity logging",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.LogEnabled = true
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Printf("Logging enabled. Logs will be written to:\n  %s\n", config.LogPath())
		return nil
	},
}

var logDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable persistent activity logging",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		cfg.LogEnabled = false
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Println("Logging disabled. Existing logs are preserved.")
		return nil
	},
}

var logClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Delete all local logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		logPath := config.LogPath()
		if err := os.Remove(logPath); err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No log file to clear.")
				return nil
			}
			return err
		}
		fmt.Println("Logs cleared.")
		return nil
	},
}

var logPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the log file path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.LogPath())
	},
}

var logStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show aggregate statistics from local logs",
	RunE: func(cmd *cobra.Command, args []string) error {
		logPath := config.LogPath()

		f, err := os.Open(logPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("No log file found. Enable logging with: bf log enable")
				return nil
			}
			return err
		}
		defer f.Close()

		var (
			total     int
			successes int
			failures  int
			battles   int
			wins      int
			purchases int
			totalMs   int64
		)

		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
		for scanner.Scan() {
			var entry client.LogEntry
			if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
				continue
			}

			total++
			totalMs += entry.DurationMs
			if entry.Success {
				successes++
			} else {
				failures++
			}

			// Count battles
			if entry.Method == "POST" && entry.Path == "/battle" && entry.Success {
				battles++
				// Check if we won by looking at the response
				var resp struct {
					Winner string `json:"winner"`
				}
				if err := json.Unmarshal(entry.Response, &resp); err == nil && resp.Winner == "attacker" {
					wins++
				}
			}

			// Count purchases
			if entry.Method == "POST" && (strings.Contains(entry.Path, "/purchase") || strings.Contains(entry.Path, "/shop/purchase")) && entry.Success {
				purchases++
			}
		}

		if total == 0 {
			fmt.Println("No log entries.")
			return nil
		}

		avgMs := totalMs / int64(total)

		bold := color.New(color.Bold)

		bold.Println("Activity Statistics")
		fmt.Println(strings.Repeat("─", 40))
		fmt.Printf("  Total requests:    %d\n", total)
		fmt.Printf("  Successes:         %s\n", color.GreenString("%d", successes))
		fmt.Printf("  Failures:          %s\n", color.RedString("%d", failures))
		fmt.Printf("  Success rate:      %.1f%%\n", float64(successes)/float64(total)*100)
		fmt.Printf("  Avg latency:       %dms\n", avgMs)
		fmt.Println()
		if battles > 0 {
			fmt.Printf("  Battles fought:    %d\n", battles)
			fmt.Printf("  Battles won:       %d\n", wins)
			fmt.Printf("  Win rate:          %.1f%%\n", float64(wins)/float64(battles)*100)
		}
		if purchases > 0 {
			fmt.Printf("  Purchases:         %d\n", purchases)
		}
		fmt.Printf("\n  Log file: %s\n", config.LogPath())

		return nil
	},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
