package utils

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dr4ghs/caigobot-discord/bot"
)

var NotCommandChannelInteractionResponse = &discordgo.InteractionResponse{
	Type: discordgo.InteractionResponseChannelMessageWithSource,
	Data: &discordgo.InteractionResponseData{
		Flags:   discordgo.MessageFlagsEphemeral,
		Content: "Non puoi usare i comandi di CaigoBot in questo canale.",
	},
}

func NotSentInCommandChannel(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if i.ChannelID == bot.ActiveGuild.CommandChannelID {
		return false
	}

	s.InteractionRespond(i.Interaction, NotCommandChannelInteractionResponse)

	return true
}
