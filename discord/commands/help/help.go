package help

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dr4ghs/caigobot-discord/commands/join"
)

// Help discordgo application command
var Command = &discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Lista dei comandi disponibili",
}

/*
Activates the help command

Parameters:
	- cmdHandlers: command handlers
	- compHandlers: component handlers
*/
func Activate(cmdHandlers map[string]func(*discordgo.Session, *discordgo.InteractionCreate),
	compHandlers map[string]func(*discordgo.Session, *discordgo.InteractionCreate)) {
	cmdHandlers["help"] = handler
	compHandlers["help_selector"] = helpSelectorHandler
}

/*
Main command handler

Parameters:
	- s: pointer to discordgo.Session
	- i: pointer to discordgo.InteractionCreate
*/
func handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Title: "Help",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "help_selector",
							Placeholder: "Scegli il comando",
							Options: []discordgo.SelectMenuOption{
								{
									Label: "help",
									Value: "help",
								},
								{
									Label: "join",
									Value: "join",
								},
							},
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Chiudi",
							Style:    discordgo.SecondaryButton,
							Disabled: false,
							CustomID: "delete",
						},
					},
				},
			},
		},
	})
}

/*
Help command selector handler

Parameters:
	- s: pointer to discordgo.Session
	- i: pointer to discordgo.InteractionCreate
*/
func helpSelectorHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd := i.MessageComponentData().Values[0]
	var cmdDetails *discordgo.MessageEmbed

	switch cmd {
	case "help":
		cmdDetails = &discordgo.MessageEmbed{
			Title:       Command.Name,
			Description: Command.Description,
		}
	case "join":
		cmdDetails = &discordgo.MessageEmbed{
			Title:       join.Command.Name,
			Description: join.Command.Description,
		}
	}

	s.ChannelMessageEditEmbed(i.ChannelID, i.Message.ID, cmdDetails)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})
}
