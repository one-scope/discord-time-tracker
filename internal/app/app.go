package app

import (
	"io"
	"log"
	"os"

	"github.com/one-scope/discord-time-tracker/internal/config"
	"github.com/one-scope/discord-time-tracker/internal/discord"
	"gopkg.in/yaml.v3"
)

type App struct {
	DiscordBot *discord.Bot
	LogFile    *os.File
}

func New(aConfigPath string) (*App, error) {
	// 設定ファイル読み込み
	tConfig := config.Config{}
	tFile, tError := os.OpenFile(aConfigPath, os.O_RDONLY, 0)
	if tError != nil {
		return nil, tError
	}
	defer func() {
		if tError := tFile.Close(); tError != nil {
			log.Println(tError)
		}
	}()
	if tError := yaml.NewDecoder(tFile).Decode(&tConfig); tError != nil {
		return nil, tError
	}

	// App初期化
	tApp := &App{}

	// DiscordBot初期化
	tApp.DiscordBot, tError = discord.New(&tConfig.DiscordBot)
	if tError != nil {
		return nil, tError
	}

	// ログファイル初期化
	tApp.LogFile, tError = logSettings(tConfig.Log.FilePath)
	if tError != nil {
		return nil, tError
	}

	return tApp, nil
}

func logSettings(aFilePath string) (*os.File, error) {
	tFile, tError := os.OpenFile(aFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if tError != nil {
		return nil, tError
	}
	tMultiLogFile := io.MultiWriter(os.Stdout, tFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetOutput(tMultiLogFile)

	return tFile, nil
}

func (aApp *App) Close() error {
	var tError error
	tError = aApp.DiscordBot.Close()
	tError = aApp.LogFile.Close()
	return tError
}
