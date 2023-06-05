package main

type Config struct {
	DiscordBot DiscordBotConfig `yaml:"discordbot"`
}

type DiscordBotConfig struct {
	DiscordBotToken string `yaml:"discord_bot_token"`
	DataDir         string `yaml:"data_dir"`
	ExecutionTiming string `yaml:"execution_timing"`
}
