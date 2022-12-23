package bot

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Guild struct {
	ID               string
	Admins           []string
	CommandChannelID string
	Commands         []string
}

var (
	ActiveGuild = &Guild{
		ID: os.Getenv("DISCORD_GUILD_ID"),
		Admins: []string{
			"189090327825612801",
		},
		CommandChannelID: "1045645143371423807",
		Commands: []string{
			"join",
		},
	}
)

func GetGuildName(id string) (*string, error) {
	for _, g := range Session.State.Guilds {
		if g.ID == id {
			return &g.Name, nil
		}
	}

	return nil, fmt.Errorf("cannot find any guild with ID %s", id)
}

func GetChannel(channelId *string, channelName *string) (*discordgo.Channel, error) {
	if channelId == nil && channelName == nil {
		return nil, fmt.Errorf("no channel ID nor name given")
	}

	guild, _ := Session.State.Guild(ActiveGuild.ID) // TODO handle error
	for _, c := range guild.Channels {
		if (channelId != nil && c.ID == *channelId) ||
			(channelName != nil && c.Name == *channelName) {
			return c, nil
		}
	}

	return nil, fmt.Errorf("cannot find any channel with ID %s or name %s", *channelId, *channelName)
}

func GetMember(memberId *string, memberName *string) (*discordgo.Member, error) {
	if memberId == nil && memberName == nil {
		return nil, fmt.Errorf("no member ID nor name given")
	}

	guild, err := Session.State.Guild(ActiveGuild.ID)
	if err != nil {
		return nil, err
	}

	for _, m := range guild.Members {
		if (memberId != nil && m.User.ID == *memberId) ||
			(memberName != nil && m.User.Username == *memberName) {
			return m, nil
		}
	}

	if memberId == nil {
		return nil, fmt.Errorf("cannot find any member of name %s", *memberName)
	}

	return nil, fmt.Errorf("cannot find any member with ID %s", *memberId)
}
