package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(guideCmd)
}

var guideCmd = &cobra.Command{
	Use:   "guide",
	Short: "Show game mechanics and strategy guide",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		data, err := c.Get("/guide")
		if err != nil {
			return err
		}

		if jsonOutput {
			fmt.Println(string(data))
			return nil
		}

		var resp struct {
			Guide json.RawMessage `json:"guide"`
		}
		json.Unmarshal(data, &resp)

		var guide map[string]json.RawMessage
		json.Unmarshal(resp.Guide, &guide)

		bold := color.New(color.Bold)
		cyan := color.New(color.FgCyan)
		green := color.New(color.FgGreen)

		bold.Println("  Battle Faeries — Game Guide")
		fmt.Println(strings.Repeat("─", 50))
		fmt.Println()

		// Stamina
		var stamina struct {
			MaxStamina    int `json:"maxStamina"`
			CostPerBattle int `json:"costPerBattle"`
		}
		json.Unmarshal(guide["stamina"], &stamina)
		cyan.Println("  Stamina")
		fmt.Printf("    Max: %d  |  Cost/battle: %d\n", stamina.MaxStamina, stamina.CostPerBattle)
		fmt.Println("    Regen (free):    1 per 3 min (20/hr)")
		fmt.Println("    Regen (premium): 1 per 2 min (30/hr)")
		fmt.Println()

		// Stats
		cyan.Println("  Stat Formulas")
		fmt.Println("    Attack  = strength + floor(magic * 0.5) + equipment")
		fmt.Println("    Defense = floor(str * 0.7) + floor(agi * 0.7) + equipment")
		fmt.Println("    Speed   = agility + equipment")
		fmt.Println("    Crit    = 3 + floor(magic / 5) + equipment  (percentage)")
		fmt.Println("    HP      = hp + equipment  (1 allocated point = +10 HP)")
		fmt.Println("    Skill Power = baseValue * (1 + magic * 0.015)")
		fmt.Println("    Base stats: HP 200 | STR 5 | AGI 5 | MAG 4")
		fmt.Printf("    Stat points per level: %d\n\n", 3)

		// Leveling
		cyan.Println("  Leveling")
		fmt.Println("    XP to next level = 50 × level + 50  (1→2: 100, 5→6: 300, 10→11: 550)")
		fmt.Println("    XP per faerie: 1 XP per round survived + share of team bonus")
		fmt.Println("    Win bonus: kills × 50 ÷ team size  (kills = 5 - enemy survivors)")
		fmt.Println("    Loss bonus: 10 ÷ team size  |  Draw bonus: 20 ÷ team size")
		fmt.Println("    Level up: +3 stat points to allocate")
		fmt.Println("    Stat reset: bf stats <id> --reset (returns to base, recovers all points)")
		fmt.Println()

		// Battle
		cyan.Println("  Battle")
		fmt.Println("    Team size: 5 faeries  |  Max turns: 200")
		fmt.Println("    Damage = attack * (100/(100+defense)) * variance(±10%)")
		fmt.Println("    Crit: 2x  |  Element advantage: 1.5x  |  Weakness: 0.67x")
		fmt.Println("    Speed soft-cap at 50 (sqrt scaling above)")
		fmt.Println()

		// Elements
		cyan.Println("  Elements")
		fmt.Println("    Fire    → strong vs Nature, Shadow  | weak vs Water, Void")
		fmt.Println("    Water   → strong vs Fire, Void      | weak vs Nature, Light")
		fmt.Println("    Nature  → strong vs Water, Light     | weak vs Fire, Shadow")
		fmt.Println("    Light   → strong vs Shadow, Water    | weak vs Void, Nature")
		fmt.Println("    Shadow  → strong vs Nature, Void     | weak vs Light, Fire")
		fmt.Println("    Void    → strong vs Light, Fire      | weak vs Shadow, Water")
		fmt.Println()

		// Gold
		cyan.Println("  Gold (victory bonus by Fae Court bracket)")
		fmt.Println("    Wisp: 250  |  Sprite: 500  |  Sylph: 900  |  Archfae: 1500  |  Mythic: 2200")
		fmt.Println("    Kill bonus: 50g per enemy faerie defeated")
		fmt.Println("    Steal on win: min(loser_gold * 2%, 300g)")
		fmt.Println()

		// Trophies
		cyan.Println("  Trophies (percentile brackets — Fae Court)")
		fmt.Println("    Wisp(bottom 60%): 8-35   |  Sprite(top 40%): 10-40  |  Sylph(top 15%): 12-45")
		fmt.Println("    Archfae(top 5%): 15-50   |  Mythic(top 1%): 18-60")
		fmt.Println("    Under 1000 players: everyone is Wisp")
		fmt.Println("    Win streak: +5%/win (max +25%)  |  Underdog: 1.5x if outpowered")
		fmt.Println("    Loser loses 80% of winner's gain")
		fmt.Println()

		// Equipment
		cyan.Println("  Equipment Tiers")
		fmt.Println("    T1 Common:    1-3k gold     T2 Uncommon: 5-13k gold")
		fmt.Println("    T3 Rare:      25-50k gold   T4 Epic:     100-200k gold")
		fmt.Println("    T5 Legendary: 400-880k gold")
		fmt.Println("    Slots: weapon, armor, accessory (1 each per faerie)")
		fmt.Println()

		// Rate Limits
		cyan.Println("  Rate Limits")
		fmt.Println("    Global: 200 req/min  |  Reads: 120/min  |  Writes: 30/min  |  Battles: 10/min")
		fmt.Println()

		// Premium Tiers
		cyan.Println("  Free vs Premium ($4.99/mo)")
		fmt.Println("    Agent API is FREE for all players.")
		fmt.Println("    Feature            Free              Premium")
		fmt.Println("    ─────────────────  ────────────────  ────────────────")
		fmt.Println("    API Keys           1                 5")
		fmt.Println("    AI Faeries Chat    5 msgs/day        50 msgs/day")
		fmt.Println("    AI Commentary      1/day             10/day")
		fmt.Println("    Stamina Regen      1 per 3 min       1 per 2 min")
		fmt.Println("    Replay Retention   1 day             10 days")
		fmt.Println()

		// Leaderboard
		cyan.Println("  Leaderboard")
		fmt.Println("    Rank = floor(totalPower * 0.4 + trophies * 0.6)")
		fmt.Println("    totalPower = sum of all faerie base stats (no equipment)")
		fmt.Println()

		// Strategy
		green.Println("  Strategy Tips")
		fmt.Println("    1. Farm weakest opponent for gold + trophies")
		fmt.Println("    3. Buy T1 weapons for all 5 faeries (biggest power spike)")
		fmt.Println("    4. Allocate stat points: str + agi are strong early")
		fmt.Println("    5. Buy armor after all weapons equipped")
		fmt.Println("    6. Save for T2 gear mid-game")
		fmt.Println("    7. Win streaks compound — chain battles for +25% trophies")
		fmt.Println()
		fmt.Println("    Use: bf guide --json  for full machine-readable guide")

		return nil
	},
}
