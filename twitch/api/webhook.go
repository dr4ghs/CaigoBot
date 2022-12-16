package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix"
)

const (
	MESSAGE_TYPE              = "twitch-eventsub-message-type"
	MESSAGE_TYPE_VERIFICATION = "webhook_callback_verification"
)

type StreamOnlineBody struct {
	Subscription helix.EventSubSubscription       `json:"subscription"`
	Event        helix.EventSubStreamOfflineEvent `json:"event"`
	Challenge    string                           `json:"challenge"`
}

func init() {
	Router.POST("/twitch/notify", streamOnlineNotify)
}

func streamOnlineNotify(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
	}
	defer ctx.Request.Body.Close()

	req := &StreamOnlineBody{}
	json.Unmarshal(body, req)
	fmt.Printf("%+v\n", req)

	if !helix.VerifyEventSubNotification("mynewalertsecret", ctx.Request.Header, string(body)) {
		log.Println("No valid signature on subscription")
		return
	} else {
		log.Println("Verified signature for subscription")
	}

	ctx.Writer.WriteString(req.Challenge)
	ctx.Writer.Flush()
}
