package main

import (
	"log"

	"github.com/dr4ghs/caigobot-twitch/api"
	"github.com/gin-gonic/autotls"
)

func main() {
	log.Fatal(autotls.RunWithManager(api.Router, api.Manager))
}
