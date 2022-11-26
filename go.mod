module github.com/un4gi/mBot

go 1.17

require (
	github.com/bwmarrin/discordgo v0.23.2
	github.com/pquerna/otp v1.3.0
)

retract (
	v1.0.0 // Causing issues with install
)
