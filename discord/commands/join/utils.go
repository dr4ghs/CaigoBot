package join

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dr4ghs/caigobot-discord/bot"
)

/*
Send a join request notification in the notification channel

Parameters:
	- s: pointer to discordgo.Session
	- i: pointer to discordgo.InteractionCreate
	- reqId: string UUID of the request
*/
func joinRequestSendNotification(s *discordgo.Session, i *discordgo.InteractionCreate, reqId string) error {
	jr := pendingRequests[reqId]
	member, channel, err := jr.getMemberAndChannelDetails()
	if err != nil {
		return err
	}

	msg, err := s.ChannelMessageSendComplex(jrConfig.NotificationChannel, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Join Request",
				Description: "Un utente vuole accedere ad un canale di streaming.",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "ID",
						Value: reqId,
					},
					{
						Name:  "Utente",
						Value: member.User.Username,
					},
					{
						Name:  "Canale",
						Value: channel.Name,
					},
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: joinRequestAcceptDeclineComponents,
			},
		},
	})
	if err != nil {
		return err
	}

	jr.NotificationMsg = msg.ID

	return nil
}

/*
Resolves the join request

Parameters:
	- s: pointer to discordgo.Session
	- i: pointer to discordgo.InteractionCreate
	- accepted: whether or not the request has been accepted
*/
func joinRequestResolve(s *discordgo.Session, i *discordgo.InteractionCreate, accepted bool) bool {
	reqId := i.Message.Embeds[0].Fields[0].Value
	joinReq := pendingRequests[reqId]

	joinReq.Accepted = accepted

	if !joinRequestMemberConnected(s, i, joinReq.Member) && accepted {
		log.Printf("The user %s is not connected to any voice channel", i.Member.User.Username)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "L'utente non è connesso ad un canale vocale",
			},
		})

		return false
	}

	joinReq.resolve(s)
	joinReq.updateRequestStatusMessage(s, i, accepted)
	delete(pendingRequests, joinReq.ID)

	s.ChannelMessageDelete(jrConfig.NotificationChannel, joinReq.NotificationMsg)

	return true
}

/*
Chekcs if the user who requested to join is already connected to a voice channel

Parameters:
	- s: pointer to discordgo.Session
	- i: pointer to discordgo.InteractionCreate
	- userId: Discord ID of the user tho requested to join a stream room

Returns:
	- bool: true if the user is already connected to a voice channel
*/
func joinRequestMemberConnected(s *discordgo.Session, i *discordgo.InteractionCreate, userId string) bool {
	if vs, _ := s.State.VoiceState(bot.ActiveGuild.ID, userId); vs == nil {
		s.InteractionRespond(i.Interaction, notConnectedWarningResponse)

		return false
	}

	return true
}

/*
Build the request update message

Parameters:
	- username: the usename of the one who requested to join
	- channelName: name of the channel to join

Returns:
	- string: the updated status message
*/
func getUpdateStatusMessage(username string, channelName string) string {
	return fmt.Sprintf("%s ha chiesto di unirsi al canale \"%s\"", username, channelName)
}
