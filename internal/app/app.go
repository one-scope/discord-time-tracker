package app

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/one-scope/discord-time-tracker/internal/db"
	"github.com/one-scope/discord-time-tracker/internal/discord"
)

type App struct {
	DiscordBot *discord.Bot
	LogFile    *os.File
}

func New() (*App, error) {
	var tError error

	// App初期化
	tApp := &App{}

	// ログ環境変数読み込み
	tLogConfig := LoadLogConfig()
	// ログファイル初期化
	tApp.LogFile, tError = logSettings(tLogConfig.FilePath)
	if tError != nil {
		return nil, tError
	}

	// DB環境変数読み込み
	tDBConfig := db.LoadDBConfig()
	// データベース初期化
	tDB, tError := db.New(&tDBConfig)
	if tError != nil {
		return nil, fmt.Errorf("DB init error: %w", tError)
	}

	// Discord環境変数読み込み
	tDiscordConfig := discord.LoadDiscordBotConfig()
	// DiscordBot初期化
	tApp.DiscordBot, tError = discord.New(&tDiscordConfig, tDB)
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
