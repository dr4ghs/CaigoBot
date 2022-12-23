package commands

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/dr4ghs/caigobot-discord/bot"
	"github.com/dr4ghs/caigobot-discord/commands/join"
)

const (
	REMOVE_CMD_VAR_NAME = "REMOVE_CMD"
)

var (
	// Indicates if the bot has to remove all registered command of
	// a guild before closing.
	RemoveCommandsCfg = false

	// Commands registered for every guild
	registeredCommands map[string][]*discordgo.ApplicationCommand

	// All avaiable commands
	Commands = map[string]*discordgo.ApplicationCommand{
		"join": join.Command,
	}

	// All avaiable command handlers
	CmdHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	CompHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}

	guildCommands = map[string][]string{
		"656203337955541031": {
			"join",
		},
	}
)

func init() {
	rm_cmd, err := strconv.ParseBool(os.Getenv(REMOVE_CMD_VAR_NAME))
	if err != nil {
		log.Printf("Cannot parse variable '%s', using default value 'false'", REMOVE_CMD_VAR_NAME)
	} else {
		RemoveCommandsCfg = rm_cmd
	}

	for k, _ := range Commands {
		switch k {
		case "join":
			join.Activate(CmdHandlers, CompHandlers)
		}
	}

	bot.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := CmdHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := CompHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	registeredCommands = make(map[string][]*discordgo.ApplicationCommand)
}

// RegisterCommands registers a list of commands for every specified guild
func RegisterCommands(s *discordgo.Session) {
	fmt.Printf("Registering bot commands for guild %s...\n", bot.ActiveGuild.ID)
	registeredCommands[bot.ActiveGuild.ID] = make([]*discordgo.ApplicationCommand, len(Commands))

	for i, cmdName := range bot.ActiveGuild.Commands {
		fmt.Printf("Adding command %v... ", cmdName)

		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, bot.ActiveGuild.ID, Commands[cmdName])
		if err != nil {
			fmt.Println("ERROR")
			log.Panicf("Cannot create '%v' command: %v", cmdName, err)
		}
		registeredCommands[bot.ActiveGuild.ID][i] = cmd

		fmt.Println("OK!")
	}
}

// RemoveCommands removes all registered commands for every specified guild
func RemoveCommands(s *discordgo.Session) {
	fmt.Printf("Deleting bot commands for guild %s...\n", bot.ActiveGuild.ID)

	for _, cmd := range registeredCommands[bot.ActiveGuild.ID] {
		fmt.Printf("Removing command %v... ", cmd.Name)

		err := s.ApplicationCommandDelete(s.State.User.ID, bot.ActiveGuild.ID, cmd.ID)
		if err != nil {
			fmt.Println("ERROR")
			log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
		}

		fmt.Println("OK!")
	}
}
