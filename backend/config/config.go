package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func Init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	} else {
		log.Println(".env file loaded successfully")
	}

	Database()
	app()
	// Set up Viper to automatically use environment variables
	viper.AutomaticEnv()
}

func app() {
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")

	conf("TOKEN_SECRET_USER", "kkrg0d")
	conf("TOKEN_DURATION_USER", 24*time.Hour)
	conf("STEAM_API_KEY", "S")
	conf("DISCORD_CLIENT_ID", "ss")
	conf("DISCORD_CLIENT_SECRET", "")
	conf("DISCORD_REDIRECT_URI", "")
}
