package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/one-scope/discord-time-tracker/internal/app"
)

func main() {
	// App初期化
	tApp, tError := app.New()
	if tError != nil {
		log.Fatal(tError)
	}
	defer func() {
		if tError := tApp.Close(); tError != nil {
			log.Panic(tError)
		}
	}()

	// DiscordBot起動
	if tError := tApp.DiscordBot.Start(); tError != nil {
		log.Fatal(tError)
	}

	// 常駐化
	tBot := make(chan os.Signal, 1)
	signal.Notify(tBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-tBot
}
