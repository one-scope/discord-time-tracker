package main

import (
	"os"

	"github.com/one-scope/discord-time-tracker/internal/discordbot"
	"gopkg.in/yaml.v2"
)

type App struct {
	DiscordBot *discordbot.DiscordBot
}

func newApp(aConfigPath string) (*App, error) {
	// 設定ファイル読み込み
	tConfig := Config{}
	tFile, tError := os.OpenFile(aConfigPath, os.O_RDONLY, 0)
	if tError != nil {
		return nil, tError
	}
	if tError := yaml.NewDecoder(tFile).Decode(&tConfig); tError != nil {
		return nil, tError
	}
	tFile.Close()

	// App初期化
	tApp := &App{}

	// DiscordBot初期化
	tApp.DiscordBot, tError = discordbot.New(tConfig.DiscordBot.DiscordBotToken)
	if tError != nil {
		return nil, tError
	}

	return tApp, nil
}
