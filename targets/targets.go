package targets

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Tink2Hack/mBot/auth"
	"github.com/Tink2Hack/mBot/config"
	"github.com/Tink2Hack/mBot/data"
	"github.com/Tink2Hack/mBot/discord"
	"github.com/Tink2Hack/mBot/env"
	"github.com/Tink2Hack/mBot/requests"
)

var notified []string

func CheckTargets(target string) {
	statuscode, body := requests.DoGetRequest(target)

	defer body.Close()
	if statuscode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil && err != context.Canceled {
			log.Printf(env.ErrorColor, err)
		}

		var newTargets data.TargetData
		err = json.Unmarshal(bodyBytes, &newTargets)
		bodyString := string(bodyBytes)

		if len(bodyString) > 3 { // Check to see if there is a JSON response
			for i := 0; i < len(newTargets); i++ { // Iterate over the number of unregistered targets
				if newTargets[i].Category.ID >= 3 { // Check to see if new target is a web/host target
				} else {
					if !newTargets[i].Registered {
						log.Printf(env.InfoColor, "Onboarding to "+newTargets[i].Codename+" - "+newTargets[i].Slug)
						OnboardTarget("https://platform.synack.com/api/targets/" + newTargets[i].Slug + "/signup")
						discord.ConnectDiscord("Onboarded to " + newTargets[i].Codename + "! Go find some vulns!!")
					}
				}
			}
		}
	} else if statuscode == http.StatusUnauthorized {
		config.LoggedIn = false
		auth.RenewSession()
	}
}

func OnboardTarget(target string) {
	var jsonStr = []byte(`{"ResearcherListing":{"terms":1}}`)
	requests.DoPostRequest(target, jsonStr)
}

func CheckForQR(target string) {
	statuscode, body := requests.DoGetRequest(target)

	defer body.Close()
	if statuscode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil && err != context.Canceled {
			log.Printf(env.ErrorColor, err)
		}

		var newTargets data.TargetData
		err = json.Unmarshal(bodyBytes, &newTargets)
		bodyString := string(bodyBytes)

		if len(bodyString) > 3 { // Check to see if there is a JSON response
			for i := 0; i < len(newTargets); i++ { // Iterate over the number of unregistered targets
				if newTargets[i].Category.ID < 3 && !checkIfNotified(notified, newTargets[i].Codename) { // Check to see if new target is a web/host target
					log.Printf(env.SuccessColor, "[!] Target "+newTargets[i].Codename+" is in QR!")
					discord.ConnectDiscord("Target " + newTargets[i].Codename + " is in QR! Go find some vulns!!")
					notified = append(notified, newTargets[i].Codename)
				}
			}
		}
	} else if statuscode == http.StatusUnauthorized {
		config.LoggedIn = false
		auth.RenewSession()
	}
}

func checkIfNotified(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
