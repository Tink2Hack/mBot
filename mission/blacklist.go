package mission

import (
	"log"
	"strings"

	"github.com/Tink2Hack/mBot/config"
	"github.com/Tink2Hack/mBot/env"
)

var blacklist []string

func CheckBlacklist(str string) bool {
	for _, v := range blacklist {
		if v == str {
			return true
		}
	}
	return false
}

func AddBlacklist(codename string) {
	for s := range config.DoNotGrab {
		if strings.Contains(codename, config.DoNotGrab[s]) {
			blacklist = append(blacklist, codename)
			log.Printf(env.WarningColor, "Added "+codename+" to the blacklist.")
		}
	}
}
