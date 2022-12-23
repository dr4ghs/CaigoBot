package bot

import (
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

const (
	DISCORD_BOT_TOKEN_VAR_NAME   = "DISCORD_BOT_TOKEN"
	DISCORD_BOT_INTENTS_VAR_NAME = "DISCORD_BOT_INTENTS"
)

var (
	// Discord bot token
	botToken = os.Getenv(DISCORD_BOT_TOKEN_VAR_NAME)

	// Discord bot intents
	botIntents = 0

	// Discord bot session
	Session *discordgo.Session
)

func init() {
	var err error
	Session, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Cannot open the session: %e", err)
	}

	intents, err := strconv.ParseInt(os.Getenv(DISCORD_BOT_INTENTS_VAR_NAME), 0, 32)
	if err != nil {
		log.Printf("Cannot parse variable '%s', using default value '0'", DISCORD_BOT_INTENTS_VAR_NAME)
	} else {
		botIntents = int(intents)
	}

	Session.Identify.Intents = discordgo.Intent(botIntents)
}
