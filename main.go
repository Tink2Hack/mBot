package main

import (
	"log"
	"time"

	"github.com/Tink2Hack/mBot/config"
	"github.com/Tink2Hack/mBot/env"
	"github.com/Tink2Hack/mBot/mission"
	"github.com/Tink2Hack/mBot/requests"
	"github.com/Tink2Hack/mBot/targets"
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
