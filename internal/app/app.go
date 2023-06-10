package app

import (
	"io"
	"log"
	"os"

	"github.com/one-scope/discord-time-tracker/internal/discordbot"
	"gopkg.in/yaml.v2"
)

type App struct {
	DiscordBot *discordbot.DiscordBot
	LogFile    *os.File
}

func New(aConfigPath string) (*App, error) {
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

	// ログファイル初期化
	tApp.LogFile, tError = LogSettings(tConfig.Log.FilePath)
	if tError != nil {
		return nil, tError
	}

	return tApp, nil
}

func LogSettings(aFilePath string) (*os.File, error) {
	tFile, tError := os.OpenFile(aFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if tError != nil {
		return nil, tError
	}
	tMultiLogFile := io.MultiWriter(os.Stdout, tFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(tMultiLogFile)

	return tFile, nil
}
