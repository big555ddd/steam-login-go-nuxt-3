package model

type SteamUser struct {
	ID        int64  `bun:"id,pk,autoincrement" json:"id"`
	SteamID   string `bun:"steam_id" json:"steam_id"`
	DiscordID string `bun:"discord_id" json:"discord_id"`
	Points    int64  `bun:"points" json:"points"`
	VipPoints int64  `bun:"vip_points" json:"vip_points"`

	CreateUpdateUnixTimestamp
	SoftDelete
}

type SteamInfo struct {
	SteamID      string `bun:"steam_id" json:"steam_id"`
	PersonaName  string `bun:"persona_name" json:"persona_name"`
	AvatarSmall  string `bun:"avatar_small" json:"avatar_small"`
	AvatarMedium string `bun:"avatar_medium" json:"avatar_medium"`
	AvatarFull   string `bun:"avatar_full" json:"avatar_full"`
}
