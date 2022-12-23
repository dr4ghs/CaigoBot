package bot

import "os"

type Guild struct {
	ID       string
	Commands []string
}

var (
	Guilds = []*Guild{
		{
			ID: os.Getenv("DISCORD_GUILD_ID"),
			Commands: []string{
				"greets",
			},
		},
	}
)
