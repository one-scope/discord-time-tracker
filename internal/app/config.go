package app

type Config struct {
	DiscordBot DiscordBotConfig `yaml:"discordbot"`
	Log        LogConfig        `yaml:"log"`
}

type DiscordBotConfig struct {
	DiscordBotToken string `yaml:"discord_bot_token"`
	DataDir         string `yaml:"data_dir"`
	ExecutionTiming string `yaml:"execution_timing"`
}
type LogConfig struct {
	FilePath string `yaml:"file_path"`
}
