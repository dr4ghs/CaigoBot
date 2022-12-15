package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	cfg "github.com/dr4ghs/caigobot-discord/configs"
)

var (
	// Discord bot session
	Session *discordgo.Session
)

func init() {
	var err error
	Session, err = discordgo.New("Bot " + cfg.Config.BotToken)
	if err != nil {
		log.Fatalf("Cannot open the session: %e", err)
	}

	Session.Identify.Intents = discordgo.Intent(cfg.Config.Intents)
}
