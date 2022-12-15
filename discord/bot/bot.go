package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dr4ghs/caigobot-discord/configs"
)

var (
	// Discord bot session
	Session *discordgo.Session
)

func init() {
	var err error
	Session, err = discordgo.New("Bot " + configs.Token)
	if err != nil {
		log.Fatalf("Cannot open the session: %e", err)
	}

	defineIntents()
}

func defineIntents() {
	var intents discordgo.Intent
	for _, element := range configs.Intents {
		intents |= element
	}

	Session.Identify.Intents = intents
}
