package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dr4ghs/caigobot-discord/bot"
	"github.com/dr4ghs/caigobot-discord/bot/commands"
	"github.com/dr4ghs/caigobot-discord/configs"
)

func main() {
	if err := bot.Session.Open(); err != nil {
		log.Fatalf("Cannot enstablish connection: %e", err)
	}
	commands.RegisterCommands(bot.Session)

	defer bot.Session.Close()

	fmt.Println("Bot is running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if configs.RemoveCommand {
		commands.RemoveCommands(bot.Session)
	}

	fmt.Println("Gracefully closed the application.")
	os.Exit(0)
}
