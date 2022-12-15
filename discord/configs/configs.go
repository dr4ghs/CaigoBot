package configs

import (
	"flag"
	"os"

	"github.com/bwmarrin/discordgo"
)

const (
	// Discord bot token environment variable name
	BOT_TOKEN_VAR_NAME = "CAIGOBOT_TOKEN"

	// Discord guild ID environment variable name
	DISCORD_GUILD_ID_VAR_NAME = "DISCORD_GUILD_ID"
)

var (
	// Discord bot intents
	Intents []discordgo.Intent = []discordgo.Intent{
		discordgo.IntentGuildMessages,
		discordgo.IntentDirectMessages,
	}

	// Discord bot connection token
	Token string

	// Discord guild ID
	GuildID string

	// Discord bot command flags
	RemoveCommand bool
)

func init() {
	flag.StringVar(&Token, "token", os.Getenv(BOT_TOKEN_VAR_NAME), "Bot token")
	flag.StringVar(&GuildID, "guild_id", os.Getenv(DISCORD_GUILD_ID_VAR_NAME), "ID of the Discord guild")
	flag.BoolVar(&RemoveCommand, "rmcmd", true, "Remove all commands after shutting or not")
	flag.Parse()
}
