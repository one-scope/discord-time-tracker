package config

import (
	"os"
	"path"
)

type Config struct {
	DiscordBot DiscordBotConfig
	Log        LogConfig
	DB         DBConfig
}

type DiscordBotConfig struct {
	DiscordBotToken       string
	FlushTimingCron       string
	DiscordErrorChannelID string
}
type LogConfig struct {
	FilePath string
}
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (aConfig *Config) LoadEnv() {

	tDiscordBotConfig := DiscordBotConfig{
		DiscordBotToken:       os.Getenv("DISCORD_BOT_TOKEN"),
		FlushTimingCron:       os.Getenv("DISCORD_FLUSH_TIMING_CRON"),
		DiscordErrorChannelID: os.Getenv("DISCORD_ERROR_CHANNEL_ID"),
	}
	tLogConfig := LogConfig{
		FilePath: path.Join(os.Getenv("LOG_FILE_BASE_PATH"), os.Getenv("LOG_FILE_NAME")),
	}
	tDBConfig := DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}

	aConfig.DiscordBot = tDiscordBotConfig
	aConfig.Log = tLogConfig
	aConfig.DB = tDBConfig

	return
}
