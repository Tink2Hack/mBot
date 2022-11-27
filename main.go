package main

import (
	"log"
	"time"

	"github.com/un4gi/mBot/config"
	"github.com/un4gi/mBot/env"
	"github.com/un4gi/mBot/mission"
	"github.com/un4gi/mBot/requests"
	"github.com/un4gi/mBot/targets"
)

func main() {

	targets.CheckTargets(requests.Urls[0])
	mission.CheckClaimed()
	for {
		log.Printf(env.InfoColor, "Checking in...")
		targets.CheckTargets(requests.Urls[0])
		if config.LoggedIn {
			targets.CheckForQR(requests.Urls[2])
			if mission.CheckWallet(requests.Urls[6]) {
				mission.CheckMissions(requests.Urls[1])
			}
		}

		secs := time.Duration(config.Delay) * time.Second
		time.Sleep(secs)
	}
}
