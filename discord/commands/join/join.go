package join

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dr4ghs/caigobot-discord/bot"
	"github.com/dr4ghs/caigobot-discord/utils"
	"github.com/google/uuid"
)

// Join discordgo application command
var Command = &discordgo.ApplicationCommand{
	Name:        "join",
	Description: "Ask to join a private streaming voice channel",
	DescriptionLocalizations: &map[discordgo.Locale]string{
		discordgo.Italian: "Chiedi di unirti ad un canale vocale di streaming",
	},
}

var (
	// Dummy configurations
	jrConfig = &JoinRequestConfig{ // TODO move to database
		NotificationChannel: "1057600703666851893",
		WaitingRoom:         "1057600551623344128",
		StreamRooms: []*JoinRequestStreamRoom{
			// Akhet's Room
			{
				ID:           "1052978287284064366",
				OwnerID:      "406957393675681803",
				AcceptPolicy: JoinRequestAcceptPolicyStaff,
				JoinPolicy:   JoinRequestJoinPolicyEveryone,
			},
			// dr4ghs' Room
			{
				ID:           "1043110020994314351",
				OwnerID:      "189090327825612801",
				AcceptPolicy: JoinRequestAcceptPolicyStaff,
				JoinPolicy:   JoinRequestJoinPolicyEveryone,
			},
			// Eliia's Room
			{
				ID:           "1044626362058678302",
				OwnerID:      "361841924627234816",
				AcceptPolicy: JoinRequestAcceptPolicyStaff,
				JoinPolicy:   JoinRequestJoinPolicyEveryone,
			},
			// Jeezy's Room
			{
				ID:           "1035311306913370263",
				OwnerID:      "313337874764398593",
				AcceptPolicy: JoinRequestAcceptPolicyStaff,
				JoinPolicy:   JoinRequestJoinPolicyEveryone,
			},
		},
	}

	// Pending join requests
	pendingRequests = make(map[string]*JoinRequest, 0)

	// Warning message if the user is not connected
	notConnectedWarningResponse = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Join Request",
					Description: "⚠️ Devi essere collegato ad un canale vocale per usare questo comando",
				},
			},
		},
	}

	// Request notification message components
	joinRequestAcceptDeclineComponents = []discordgo.MessageComponent{
		discordgo.Button{
			Label:    "Accetta",
			Style:    discordgo.SuccessButton,
			CustomID: "jr_accept",
		},
		discordgo.Button{
			Label:    "Rifiuta",
			Style:    discordgo.DangerButton,
			CustomID: "jr_decline",
		},
	}

	statusUpdateIntentResponse = &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Join Request",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Status",
						Inline: true,
					},
				},
			},
		},
	}
)

/*
Activates the join command

Parameters:
	- cmdHandlers: command handlers
	- compHandlers: component handlers
*/
func Activate(cmdHandlers map[string]func(*discordgo.Session, *discordgo.InteractionCreate),
	compHandlers map[string]func(*discordgo.Session, *discordgo.InteractionCreate)) {
	// Command handler
	cmdHandlers[JOIN_CMD_KEYWORD] = handler

	// Interaction handlers
	compHandlers[JOIN_REQUEST_CHOICE_KEYWORD] = choiceHandler
	compHandlers[JOIN_REQUEST_ACCEPT_KEYWORD] = acceptHandler
	compHandlers[JOIN_REQUEST_DECLINE_KEYWORD] = declineHandler
}

/*
Main command handler

Parameters:
	- s: pointer to discordgo.Session
	- i: pointer to discordgo.InteractionCreate
*/
func handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if utils.NotSentInCommandChannel(s, i) {
		log.Printf("The user %s did not used the join command in the command channel", i.Member.User.Username)

		return
	}

	if !joinRequestMemberConnected(s, i, i.Member.User.ID) {
		log.Printf("The user %s used the join command but is not connected to any voice channel", i.Member.User.Username)

		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Scegli il canale a cui vuoi unirti",
			Flags:   discordgo.MessageFlagsEphemeral,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "jr_choice",
							Placeholder: "Scegli un canale",
							Options:     jrConfig.getStreamRoomsSelectOptions(),
						},
					},
				},
			},
		},
	})
}

/*
Command channel choice handler

Parameters:
	- s: pointer to discordgo.Session
	- i: pointer to discordgo.InteractionCreate
*/
func choiceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !joinRequestMemberConnected(s, i, i.Member.User.ID) {
		log.Printf("The user %s choose a stream room but is not connected to any voice channel", i.Member.User.Username)

		return
	}

	choice := i.MessageComponentData().Values[0]
	channel, _ := bot.GetChannel(&choice, nil)
	responseMsg := fmt.Sprintf("%s ha chiesto di unirsi al canale \"%s\"", i.Member.User.Username, channel.Name)

	log.Println(responseMsg)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})

	s.ChannelMessageDelete(i.ChannelID, i.Message.ID)

	updateMsg := statusUpdateIntentResponse
	updateMsg.Embeds[0].Description = getUpdateStatusMessage(i.Member.User.Username, channel.Name)
	updateMsg.Embeds[0].Fields[0].Value = "🔃 Pending..."

	msg, _ := s.ChannelMessageSendComplex(i.ChannelID, updateMsg)
	reqId := uuid.NewString()
	pendingRequests[reqId] = &JoinRequest{
		ID:        reqId,
		Member:    i.Member.User.ID,
		Channel:   channel.ID,
		StatusMsg: msg.ID,
	}

	joinRequestSendNotification(s, i, reqId)
}

/*
Accepted request handler

Parameters:
	- s: pointer to discordgo.Session
	- i: pointer to discordgo.InteractionCreate
*/
func acceptHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	joinRequestResolve(s, i, true)
}

/*
Declined request handler

Parameters:
	- s: pointer to discordgo.Session
	- i: pointer to discordgo.InteractionCreate
*/
func declineHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	joinRequestResolve(s, i, false)
}
