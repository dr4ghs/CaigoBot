package join

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/dr4ghs/caigobot-discord/bot"
)

// ============================================================================
// JoinRequestConfig
// ============================================================================

// Guild join requests configurations
type JoinRequestConfig struct {
	// Notification channel Discord ID
	NotificationChannel string

	// Waiting room voice channel Discord ID
	WaitingRoom string

	// Available stream rooms
	StreamRooms []*JoinRequestStreamRoom
}

// === METHODS ================================================================

/*
Composes the discordgo.SelectMenuOptions with the regestered stream rooms

Returns:
	- []discordgo.SelectMenuOption
*/
func (cfg *JoinRequestConfig) getStreamRoomsSelectOptions() []discordgo.SelectMenuOption {
	options := make([]discordgo.SelectMenuOption, 0)

	for _, c := range cfg.StreamRooms {
		channel, _ := bot.GetChannel(&c.ID, nil)
		options = append(options, discordgo.SelectMenuOption{
			Label: channel.Name,
			Value: c.ID,
		})
	}

	return options
}

// ============================================================================
// JoinRequestStreamRoom
// ============================================================================

// Stream room data
type JoinRequestStreamRoom struct {
	// Stream Room ID as a string UUID
	ID string

	// Discord ID of the Stream Roomd
	OwnerID string

	// Discord role of staff members
	StaffRole *discordgo.Role

	// Discord role of the subscriber role
	SubscriberRole *discordgo.Role

	// Discord role of the follower role
	FollowerRole *discordgo.Role

	// Stream Room join requests accept policy
	AcceptPolicy JoinRequestAcceptPolicy

	// Stream Room join request join policy
	JoinPolicy JoinRequestJoinPolicy
}

// ============================================================================
// JoinRequest
// ============================================================================

// Join request data
type JoinRequest struct {
	// ID of the join request as UUID string
	ID string

	// Discord ID of the member who created the join request
	Member string

	// Discord ID of the requested stream room voice channel
	Channel string

	// Whether or not the request has been accepted
	Accepted bool

	// ID of the notification message
	NotificationMsg string

	// ID of the status message
	StatusMsg string
}

// === METHODS ================================================================

/*
Retrieves the discordgo instance of the member and channel registered in the join
request.

Returns:
	- *discordgo.Member
	- *discordgo.Channel
	- error
*/
func (jr *JoinRequest) getMemberAndChannelDetails() (*discordgo.Member, *discordgo.Channel, error) {
	member, err := bot.GetMember(&jr.Member, nil)
	if err != nil {
		return nil, nil, err
	}

	channel, err := bot.GetChannel(&jr.Channel, nil)
	if err != nil {
		return nil, nil, err
	}

	return member, channel, nil
}

/*
Resolves the join request and updates the status message accordingly

Parameters:
	- s: the discordgo Session
*/
func (jr *JoinRequest) resolve(s *discordgo.Session) {
	member, channel, err := jr.getMemberAndChannelDetails()

	statusUpdate := statusUpdateIntentResponse
	statusUpdate.Embeds[0].Fields[0].Value = "⛔ Rifiutata"

	if err != nil {
		log.Printf("Error while resolving join request (ID: %s): %e", jr.ID, err)

		statusUpdate.Embeds[0].Fields[0].Value = "❗ Ops, qualcosa è andato storto"
	} else {
		if jr.Accepted {
			statusUpdate.Embeds[0].Fields[0].Value = "✅ Accettata"

			s.GuildMemberMove(bot.ActiveGuild.ID, member.User.ID, &channel.ID)
		}
	}

	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      jr.StatusMsg,
		Channel: bot.ActiveGuild.CommandChannelID,
		Embeds:  statusUpdateIntentResponse.Embeds,
	})
}

/*
Updates the join request status message.

Parameters:
	- s: pointer to the discordgo.Session
	- i: pointer to the discordgo.InteractionCreate
	- accepted: whether or not the join request has been accepted
*/
func (jr *JoinRequest) updateRequestStatusMessage(s *discordgo.Session, i *discordgo.InteractionCreate, accepted bool) {
	member, channel, err := jr.getMemberAndChannelDetails()
	if err != nil {
		log.Printf("Error while updating join request (%s) status message: %e", jr.ID, err)
	}

	result := "accettata"
	acceptedLog := "accepted"
	emoji := "✅ Accettata"
	if !accepted {
		result = "rifiutata"
		acceptedLog = "declined"
		emoji = "⛔ Rifiutata"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s *La richiesta di %s per il canale %s è stata %s*", emoji, member.User.Username, channel.Name, result),
		},
	})

	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      jr.StatusMsg,
		Channel: bot.ActiveGuild.CommandChannelID,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Join Request",
				Description: fmt.Sprintf("%s ha chiesto di unirsi al canale \"%s\"", member.User.Username, channel.Name),
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Status",
						Value:  emoji,
						Inline: true,
					},
				},
			},
		},
	})

	log.Printf("%s request to join channel \"%s\" has been %s by %s.", member.User.Username, channel.Name, acceptedLog, i.Member.User.Username)
}
