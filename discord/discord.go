package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/Tink2Hack/mBot/config"
	"github.com/Tink2Hack/mBot/env"
	"github.com/Tink2Hack/mBot/requests"
)

var BotID string

func ConnectDiscord(m string) {
	d, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Println("Error starting new bot session,", err)
	}

	user, err := d.User("@me")
	if err != nil {
		log.Printf(env.ErrorColor, "Unable to send Discord notification. Checking current target connection...")
		if requests.VerifyOptimusDownload() {
			ConnectDiscord(m)
		} else {
			log.Printf(env.WarningColor, "Your target connection may be blocking connections to Discord. Giving up on sending notification.")
		}
		return
	}
	BotID = user.ID

	err = d.Open()
	if err != nil {
		log.Printf(env.ErrorColor, "Error connecting to Discord.")
	}

	log.Println(fmt.Sprint(user), "has connected to Discord!")
	SendMessage(d, m)
}

func SendMessage(d *discordgo.Session, m string) {
	_, err := d.ChannelMessageSend(config.Channel, m)
	if err != nil {
		log.Println(err)
	}
}
