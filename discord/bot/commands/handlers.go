package commands

import (
	"github.com/bwmarrin/discordgo"
)

// greets test command handler
func greets(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Hello! Welcome to CaigoBoyz!",
		},
	})
}
