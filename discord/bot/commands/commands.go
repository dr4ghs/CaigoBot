package commands

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/dr4ghs/caigobot-discord/bot"
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
	commands = map[string]*discordgo.ApplicationCommand{
		"greets": {
			Name:        "greets",
			Description: "Saluta l'utente che lo chiama",
		},
	}

	// All avaiable command handlers
	handlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"greets": greets,
	}
)

func init() {
	rm_cmd, err := strconv.ParseBool(os.Getenv(REMOVE_CMD_VAR_NAME))
	if err != nil {
		log.Printf("Cannot parse variable '%s', using default value 'false'", REMOVE_CMD_VAR_NAME)
	} else {
		RemoveCommandsCfg = rm_cmd
	}

	bot.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	registeredCommands = make(map[string][]*discordgo.ApplicationCommand)
}

// RegisterCommands registers a list of commands for every specified guild
func RegisterCommands(s *discordgo.Session) {
	for _, guild := range bot.Guilds {
		fmt.Printf("Registering bot commands for guild %s...\n", guild.ID)
		registeredCommands[guild.ID] = make([]*discordgo.ApplicationCommand, len(commands))

		for i, cmdName := range guild.Commands {
			fmt.Printf("Adding command %v... ", cmdName)

			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guild.ID, commands[cmdName])
			if err != nil {
				fmt.Println("ERROR")
				log.Panicf("Cannot create '%v' command: %v", cmdName, err)
			}
			registeredCommands[guild.ID][i] = cmd

			fmt.Println("OK!")
		}
	}
}

// RemoveCommands removes all registered commands for every specified guild
func RemoveCommands(s *discordgo.Session) {
	for _, guild := range bot.Guilds {
		fmt.Printf("Deleting bot commands for guild %s...\n", guild.ID)

		for _, cmd := range registeredCommands[guild.ID] {
			fmt.Printf("Removing command %v... ", cmd.Name)

			err := s.ApplicationCommandDelete(s.State.User.ID, guild.ID, cmd.ID)
			if err != nil {
				fmt.Println("ERROR")
				log.Panicf("Cannot delete '%v' command: %v", cmd.Name, err)
			}

			fmt.Println("OK!")
		}
	}
}
