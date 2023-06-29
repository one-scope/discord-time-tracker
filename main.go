package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/one-scope/discord-time-tracker/internal/app"
)

func main() {
	// 設定ファイルパス
	tConfigPath := *flag.String("c", "config.yml", "set config file")
	flag.Parse()

	// App初期化
	tApp, tError := app.New(tConfigPath)
	if tError != nil {
		log.Fatal(tError)
	}
	defer func() {
		// logファイル閉じる
		if tError := tApp.LogFile.Close(); tError != nil {
			log.Fatal(tError)
		}
	}()
	log.Println("Log Start")

	// DiscordBot起動
	if tError := tApp.DiscordBot.Start(); tError != nil {
		log.Fatal(tError)
	}
	defer func() {
		// DiscordBot停止
		if tError := tApp.DiscordBot.Close(); tError != nil {
			log.Fatal(tError)
		}
	}()

	// 常駐化
	tBot := make(chan os.Signal, 1)
	signal.Notify(tBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-tBot
}