package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix"
	"golang.org/x/crypto/acme/autocert"
)

var (
	Router  *gin.Engine
	Manager *autocert.Manager
)

func init() {
	gin.SetMode(gin.ReleaseMode)

	Router = gin.Default()

	Router.GET("/twitch", getAlerts)
	Router.PUT("/twitch/:user", newAlert)
	Router.DELETE("/twitch/:id", removeAlert)

	Manager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("caigobot.agsystem.dev"),
		Cache:      autocert.DirCache(".cache"),
	}
}

func getAlerts(ctx *gin.Context) {
	resp, err := Twitch.GetEventSubSubscriptions(&helix.EventSubSubscriptionsParams{})
	if err != nil {
		log.Printf("Cannot retrieve EventSub subcriptions: %v\n", err)
	}

	ctx.IndentedJSON(http.StatusOK, resp.Data.EventSubSubscriptions)
}

func newAlert(ctx *gin.Context) {
	user := ctx.Param("user")
	getUserResp, err := Twitch.GetUsers(&helix.UsersParams{
		Logins: []string{user},
	})
	if err != nil {
		log.Printf("Cannot find user %s: %v\n", user, err)
	}
	userId := getUserResp.Data.Users[0].ID

	createEventSubResp, err := Twitch.CreateEventSubSubscription(&helix.EventSubSubscription{
		Type:    helix.EventSubTypeStreamOnline,
		Version: "1",
		Condition: helix.EventSubCondition{
			BroadcasterUserID: userId,
		},
		Transport: helix.EventSubTransport{
			Method:   "webhook",
			Callback: "https://caigobot.agsystem.dev/twitch/notify",
			Secret:   "mynewalertsecret",
		},
	})
	if err != nil {
		log.Printf("Cannot subscribe user %s: %v\n", user, err)
	}

	data, _ := json.MarshalIndent(createEventSubResp, "", "  ")
	fmt.Printf("%s\n", data)
}

func removeAlert(ctx *gin.Context) {
	id := ctx.Param("id")

	resp, _ := Twitch.RemoveEventSubSubscription(id)

	data, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("%s\n", data)
}
