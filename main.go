package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/one-scope/discord-time-tracker/internal/app"
	"github.com/one-scope/discord-time-tracker/internal/discordbot"
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
	defer tApp.LogFile.Close()

	// DiscordBot並列起動
	if tError := discordbot.Start(tApp.DiscordBot); tError != nil {
		log.Fatal(tError)
	}

	// 常駐化
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot

	if tError := tApp.DiscordBot.Session.Close(); tError != nil {
		log.Fatal(tError)
	}
}
