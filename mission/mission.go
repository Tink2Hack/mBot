package mission

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Tink2Hack/mBot/auth"
	"github.com/Tink2Hack/mBot/config"
	"github.com/Tink2Hack/mBot/data"
	"github.com/Tink2Hack/mBot/discord"
	"github.com/Tink2Hack/mBot/env"
	"github.com/Tink2Hack/mBot/requests"
)

func CheckMissions(target string) {
	statuscode, body := requests.DoGetRequest(target)
	defer body.Close()

	if statuscode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil && err != context.Canceled && err != io.EOF {
			log.Println(err)
			return
		}

		var newMissions data.MissionDataV2
		err = json.Unmarshal(bodyBytes, &newMissions)
		bodyString := string(bodyBytes)
		lenResp := len(newMissions)
		if len(bodyString) > 3 { // Check to see if there is a JSON response
			var wg sync.WaitGroup
			wg.Add(lenResp)

			fmt.Println("Starting loop to attempt to grab missions...")
			for i := 0; i < lenResp; i++ {
				go func(i int) {
					defer wg.Done()

					log.Println("Saw", newMissions[i].Payout.Amount, newMissions[i].ListingCodename, "mission -", newMissions[i].Title)
					if !CheckBlacklist(newMissions[i].ListingCodename) {
						GrabMission("https://platform.synack.com/api/tasks/v1/organizations/"+newMissions[i].OrganizationUid+"/listings/"+newMissions[i].ListingUid+"/campaigns/"+newMissions[i].CampaignUid+"/tasks/"+newMissions[i].ID+"/transitions", newMissions[i].Payout.Amount, newMissions[i].ListingCodename, newMissions[i].Title, 10, newMissions[i].ID)
						log.Println("Attempted to grab", newMissions[i].Payout.Amount, newMissions[i].ListingCodename, "mission -", newMissions[i].Title)
					} else {
						log.Println("Maximum number of missions already claimed for", newMissions[i].ListingCodename+". Skipping mission.")
					}
				}(i)
			}
			wg.Wait()
			fmt.Println("Finished attempt to grab new missions...")
		}
	} else if statuscode == http.StatusUnauthorized {
		config.LoggedIn = false
		auth.RenewSession()
	}
}

func GrabMission(target string, pay int, project string, title string, n int, id string) {
	var jsonStr = []byte(`{"type": "CLAIM"}`)
	statuscode, body := requests.DoPostRequest(target, jsonStr)
	defer body.Close()
	if statuscode == http.StatusCreated {
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil && err != context.Canceled && err != io.EOF {
			log.Println(err)
			return
		}

		bodyString := string(bodyBytes)
		if len(bodyString) > 0 { // Check to see if there is a JSON response
			log.Printf(env.SuccessColor, "Grabbed $ "+fmt.Sprint(pay)+" "+project+" mission - "+title)
			discord.ConnectDiscord("Grabbed $ " + fmt.Sprint(pay) + " " + project + " mission - " + title)
			file, found := missionMap[title]
			if found {
				missionEditURL := requests.Urls[4] + id + "/evidences"
				EditMission(file, missionEditURL)
			} else {
				log.Printf(env.WarningColor, "Unable to locate template for "+project+" mission - "+title+". Don't forget to create a new template for next time!")
			}
		} else {
			log.Printf(env.WarningColor, "Error grabbing $"+fmt.Sprint(pay)+" "+project+" mission - "+title)
		}
	} else if statuscode == http.StatusUnauthorized {
		log.Printf(env.WarningColor, "Bearer token expired.")
		os.Exit(0)
	} else if statuscode == http.StatusInternalServerError {
		if n > 0 {
			log.Printf(env.WarningColor, "500 Error. Retrying...")
			n = n - 1
			GrabMission(target, pay, project, title, n, id)
		} else {
			log.Printf(env.ErrorColor, "500 Error. 10 attempts have been made. Giving up!")
		}
	} else if statuscode == http.StatusForbidden {
		log.Printf(env.ErrorColor, "Too slow! Someone grabbed it.")
	} else if statuscode == http.StatusPreconditionFailed {
		log.Printf(env.ErrorColor, "You already have the maximum number of missions for this target.")
		blacklist = append(blacklist, project)
	} else {
		log.Printf(env.WarningColor, "Status Code: "+fmt.Sprint(statuscode))
	}
	body.Close()
}
