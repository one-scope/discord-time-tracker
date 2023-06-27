package config

type Config struct {
	DiscordBot DiscordBotConfig `yaml:"discordbot"`
	Log        LogConfig        `yaml:"log"`
}

type DiscordBotConfig struct {
	DiscordBotToken string `yaml:"discord_bot_token"`
	DataPathBase    string `yaml:"data_path_base"`
	ExecutionTiming string `yaml:"execution_timing"`
}
type LogConfig struct {
	FilePath string `yaml:"file_path"`
}
