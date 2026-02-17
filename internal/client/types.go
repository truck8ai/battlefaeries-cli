package client

type PlayerStatus struct {
	ID                   string  `json:"id"`
	DisplayName          string  `json:"display_name"`
	Gold                 int     `json:"gold"`
	Stamina              int     `json:"stamina"`
	Trophies             int     `json:"trophies"`
	TotalPower           int     `json:"total_power"`
	CombatPower          int     `json:"combat_power"`
	WinStreak            int     `json:"win_streak"`
	BestWinStreak        int     `json:"best_win_streak"`
	FaerieCount          int     `json:"faerie_count"`
	IsAgentControlled    bool    `json:"is_agent_controlled"`
	IsPremium            bool    `json:"is_premium"`
	PremiumExpiresAt     *string `json:"premium_expires_at"`
}

type Faerie struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Element           string      `json:"element"`
	Level             int         `json:"level"`
	TotalXP           int         `json:"total_xp"`
	HP                int         `json:"hp"`
	Strength          int         `json:"strength"`
	Agility           int         `json:"agility"`
	Magic             int         `json:"magic"`
	UnallocatedPoints int         `json:"unallocated_points"`
	WeaponName        *string     `json:"weapon_name"`
	ArmorName         *string     `json:"armor_name"`
	AccessoryName     *string     `json:"accessory_name"`
	Skills            []FaerieSkill `json:"skills"`
}

type FaerieSkill struct {
	PlayerSkillID string `json:"player_skill_id"`
	Name          string `json:"name"`
	SkillType     string `json:"skill_type"`
	Element       string `json:"element"`
	Power         int    `json:"power"`
	SkillSlot     int    `json:"skill_slot"`
}

type OpponentFaerie struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Level   int    `json:"level"`
	Element string `json:"element"`
	Power   int    `json:"power"`
}

type Opponent struct {
	ID                string           `json:"id"`
	DisplayName       string           `json:"display_name"`
	Trophies          int              `json:"trophies"`
	TotalPower        int              `json:"total_power"`
	CombatPower       int              `json:"combat_power"`
	FaerieCount       int              `json:"faerie_count"`
	IsAgentControlled bool             `json:"is_agent_controlled"`
	Team              []OpponentFaerie `json:"team"`
}

type ShopItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	EquipmentType string `json:"equipment_type"`
	Tier          int    `json:"tier"`
	Price         int    `json:"price"`
	AttackBonus   int    `json:"attack_bonus"`
	DefenseBonus  int    `json:"defense_bonus"`
	HPBonus       int    `json:"hp_bonus"`
	SpeedBonus    int    `json:"speed_bonus"`
	CritBonus     int    `json:"crit_bonus"`
}

type Equipment struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	EquipmentType string  `json:"equipment_type"`
	Tier          int     `json:"tier"`
	AttackBonus   int     `json:"attack_bonus"`
	DefenseBonus  int     `json:"defense_bonus"`
	HPBonus       int     `json:"hp_bonus"`
	SpeedBonus    int     `json:"speed_bonus"`
	CritBonus     int     `json:"crit_bonus"`
	EquippedOn    *string `json:"equipped_on"`
}

type BattleHistory struct {
	ID                   string  `json:"id"`
	AttackerID           string  `json:"attacker_id"`
	DefenderID           string  `json:"defender_id"`
	WinnerID             *string `json:"winner_id"`
	AttackerName         string  `json:"attacker_name"`
	DefenderName         string  `json:"defender_name"`
	AttackerGoldChange   int     `json:"attacker_gold_change"`
	AttackerTrophyChange int     `json:"attacker_trophy_change"`
	BattleLogID          string  `json:"battle_log_id"`
	CreatedAt            string  `json:"created_at"`
}

type Tournament struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Status           string `json:"status"`
	MaxPlayers       int    `json:"max_players"`
	PrizeTrophies    int    `json:"prize_trophies"`
	ParticipantType  string `json:"participant_type"`
	ParticipantCount int    `json:"participant_count"`
	IsRegistered     bool   `json:"is_registered"`
}
