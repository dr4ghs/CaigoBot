package api

import (
	"log"

	cfg "github.com/dr4ghs/caigobot-twitch/configs"
	"github.com/nicklaw5/helix"
)

var Twitch *helix.Client

func init() {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     cfg.Config.ClientID,
		ClientSecret: cfg.Config.AppAccessToken,
	})
	if err != nil {
		log.Fatalf("Cannot create a Twitch client: %v", err)
	}

	resp, _ := client.RequestAppAccessToken([]string{"user:read:email"})
	client.SetAppAccessToken(resp.Data.AccessToken)

	Twitch = client
}
