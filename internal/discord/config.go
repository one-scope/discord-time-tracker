package discord

import "os"

type DiscordBotConfig struct {
	DiscordBotToken       string
	FlushTimingCron       string
	DiscordErrorChannelID string
}

func LoadDiscordBotConfig() DiscordBotConfig {
	return DiscordBotConfig{
		DiscordBotToken:       os.Getenv("DISCORD_BOT_TOKEN"),
		FlushTimingCron:       os.Getenv("DISCORD_FLUSH_TIMING_CRON"),
		DiscordErrorChannelID: os.Getenv("DISCORD_ERROR_CHANNEL_ID"),
	}
}
