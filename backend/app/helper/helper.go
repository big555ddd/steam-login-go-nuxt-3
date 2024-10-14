package helper

import (
	"context"
	"log"

	"github.com/softsrv/steamapi/steamapi"
	"github.com/spf13/viper"
)

// PlayerProfile โครงสร้างข้อมูลผู้เล่น
type PlayerProfile struct {
	SteamID      string `json:"steam_id"`
	PersonaName  string `json:"persona_name"`
	AvatarSmall  string `json:"avatar_small"`
	AvatarMedium string `json:"avatar_medium"`
	AvatarFull   string `json:"avatar_full"`
}

// ฟังก์ชันสร้าง Steam API Client
func NewSteamClient() *steamapi.Client {
	apiKey := viper.GetString("STEAM_API_KEY")
	return steamapi.NewClient(apiKey) // ใช้ NewClient ของ steamapi
}

// ฟังก์ชันดึงข้อมูลโปรไฟล์ผู้เล่น
func GetSteamProfile(steamID string) (*PlayerProfile, error) {
	client := NewSteamClient()

	player, err := client.Player(context.Background(), steamID)
	if err != nil {
		log.Printf("Error fetching Steam profile: %v", err)
		return nil, err
	}

	return &PlayerProfile{
		SteamID:      player.SteamID,
		PersonaName:  player.PersonaName,
		AvatarSmall:  player.AvatarSmall,
		AvatarFull:   player.AvatarFull,
		AvatarMedium: player.AvatarMedium,
	}, nil
}
