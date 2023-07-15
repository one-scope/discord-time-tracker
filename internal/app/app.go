package app

import (
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/one-scope/discord-time-tracker/internal/config"
	"github.com/one-scope/discord-time-tracker/internal/db"
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

	// ログファイル初期化
	tApp.LogFile, tError = logSettings(tConfig.Log.FilePath)
	if tError != nil {
		return nil, tError
	}

	// データベース初期化
	tDB, tError := db.New()
	if tError != nil {
		return nil, fmt.Errorf("DB init error: %w", tError)
	}

	// DiscordBot初期化
	tApp.DiscordBot, tError = discord.New(&tConfig.DiscordBot, tDB)
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
	var tReturnError error
	if tError := aApp.DiscordBot.Close(); tError != nil {
		tReturnError = fmt.Errorf("DiscordBot close error: %w", tError)
	}
	if tError := aApp.LogFile.Close(); tError != nil {
		tReturnError = fmt.Errorf("LogFile close error: %w: %w", tError, tReturnError)
	}
	if tError := aApp.DiscordBot.DataManager.DB.Close(); tError != nil {
		tReturnError = fmt.Errorf("DB close error: %w: %w", tError, tReturnError)
	}
	return tReturnError
}
