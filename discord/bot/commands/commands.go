package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dr4ghs/caigobot-discord/bot"
	"github.com/dr4ghs/caigobot-discord/configs"
)

var (
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
	bot.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := handlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

// RegisterCommands registers a list of commands for every specified guild
func RegisterCommands(s *discordgo.Session) {
	fmt.Println("Registering bot commands...")
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))

	for i, v := range commands {
		fmt.Printf("Adding command %v... ", v.Name)

		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, configs.GuildID, v)
		if err != nil {
			fmt.Println("ERROR")
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd

		fmt.Println("OK!")
	}
}

// RemoveCommands removes all registered commands for every specified guild
func RemoveCommands(s *discordgo.Session) {
	fmt.Println("Deleting bot commands...")

	for _, v := range registeredCommands {
		fmt.Printf("Removing command %v... ", v.Name)

		err := s.ApplicationCommandDelete(s.State.User.ID, configs.GuildID, v.ID)
		if err != nil {
			fmt.Println("ERROR")
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}

		fmt.Println("OK!")
	}
}
