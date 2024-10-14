package helper

import (
	"app/internal/logger"
	"context"
	"encoding/json"
	"fmt"

	"github.com/diamondburned/arikawa/discord"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

// ตั้งค่า OAuth2 Configuration สำหรับ Discord
func GetOAuth2Config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     viper.GetString("DISCORD_CLIENT_ID"),
		ClientSecret: viper.GetString("DISCORD_CLIENT_SECRET"),
		RedirectURL:  viper.GetString("DISCORD_REDIRECT_URI"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/api/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
		Scopes: []string{"identify"},
	}
}

// ฟังก์ชันแลกเปลี่ยน OAuth2 Code เป็น Access Token
func ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	oauth2Config := GetOAuth2Config()

	logger.Infof("code: %v", code)
	logger.Infof("oauth2Config.redir: %v", oauth2Config.RedirectURL)
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %v", err)
	}
	return token, nil
}

// ฟังก์ชันดึงข้อมูลผู้ใช้จาก Discord
func GetDiscordUser(token *oauth2.Token) (*discord.User, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(token))
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user info: %v", err)
	}
	defer resp.Body.Close()

	var user discord.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}
	return &user, nil
}
