package mission

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Tink2Hack/mBot/auth"
	"github.com/Tink2Hack/mBot/config"
	"github.com/Tink2Hack/mBot/data"
	"github.com/Tink2Hack/mBot/env"
	"github.com/Tink2Hack/mBot/requests"
)

func CheckClaimed() {
	statuscode, body := requests.DoGetRequest(requests.Urls[3])
	defer body.Close()

	if statuscode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil && err != context.Canceled && err != io.EOF {
			log.Println(err)
			return
		}
		var claimedMissions data.MissionDataV2
		err = json.Unmarshal(bodyBytes, &claimedMissions)
		if err != nil && err != context.Canceled && err != io.EOF {
			fmt.Println("Error unmarshalling data:", err)
		}
		bodyString := string(bodyBytes)
		//fmt.Println(bodyString)
		lenResp := len(claimedMissions)
		if len(bodyString) > 3 { // Check to see if there is a JSON response
			fmt.Println("Checking claimed missions...")
			for i := 0; i < lenResp; i++ {
				file, found := missionMap[claimedMissions[i].Title]
				if found {
					log.Printf(env.SuccessColor, "Populating mission data for "+claimedMissions[i].ListingCodename+" mission - "+claimedMissions[i].Title+"!")
					missionEditURL := requests.Urls[4] + claimedMissions[i].ID + "/evidences"
					EditMission(file, missionEditURL)
				} else {
					log.Printf(env.WarningColor, "Unable to locate template for "+claimedMissions[i].Title+". Don't forget to create a new template for next time!")
				}
			}
		}
	} else if statuscode == http.StatusUnauthorized {
		config.LoggedIn = false
		auth.RenewSession()
	}
}

func CheckWallet(url string) bool {
	isRoom := false
	statuscode, body := requests.DoGetRequest(url)
	defer body.Close()

	if statuscode == http.StatusOK {

		bodyBytes, err := ioutil.ReadAll(body)
		if err != nil && err != context.Canceled && err != io.EOF {
			log.Println(err)
		}

		type ClaimedAmount struct {
			ClaimedAmount int `json:"claimedAmount"`
		}

		var wallet ClaimedAmount
		err = json.Unmarshal(bodyBytes, &wallet)
		if err != nil && err != context.Canceled && err != io.EOF {
			fmt.Println("Error checking amount claimed:", err)
		}

		if len(string(bodyBytes)) > 3 {
			if wallet.ClaimedAmount < 75 { // Check to see if there is room in wallet
				isRoom = true
			}
		} else {
			isRoom = false
		}
	} else if statuscode == http.StatusUnauthorized {
		config.LoggedIn = false
		auth.RenewSession()
	}
	return isRoom
}

func EditMission(file, url string) {
	payload, _ := os.Open(file)
	defer payload.Close()
	b, _ := ioutil.ReadAll(payload)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(b))
	requests.SetHeaders(req)
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil && err != context.Canceled && err != io.EOF {
		log.Fatal(err)
	}
}
