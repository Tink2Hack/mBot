[![License](https://img.shields.io/badge/license-MIT-_red.svg)](https://opensource.org/licenses/MIT) [![Twitter Follow](https://img.shields.io/twitter/follow/un4gi_io?label=%40un4gi_io&style=social)](https://twitter.com/un4gi_io)

<img src="img/mBot.png">

# mBot (A Go-Based Synack Mission Bot)

This mission bot contains functionality for onboarding, claiming missions, Discord notifications, and auto-relogin.

---

## Pre-Requisites

If you want to get the most use out of mBot, you should follow the steps below to get a Discord bot token and Authy secret. These will allow you to both receive notifications and automatically log back in to the platform should your bot session get disconnected.

### Discord Setup Instructions

1. To get your Discord bot token, you will need to follow the steps at <https://www.writebots.com/discord-bot-token/>.
2. Once you have your Discord bot token, add the `CHANNEL_ID` and `DISCORD_TOKEN` to `config.json`.

### Authy Configuration

1. To get your Authy Secret, you will need to follow the instructions at <https://github.com/alexzorin/authy>.
2. Once you have your secret, add the `AUTHY_SECRET` to `config.json` along with your `EMAIL_ADDRESS`, and `PASSWORD` for Synack.

---

## Installation

### Building From Source

If you prefer, you can build mBot straight from the source directory:

```bash
git clone https://github.com/Un4gi/mBot.git
cd mBot
go build .
```

## Usage

Usage has been streamlined and no longer needs arguments passed. As long as you have everything configured properly, you can now simply execute the binary from the command line.

Example:

```bash
mBot
```

## Mission Templates

Good news... mBot now allows auto-population of the Intro/Testing Methodology/Conclusion fields for each claimed mission!

*That's great... but how do I do it?*

__It's simple, really!__

1. Store your templates as JSON files in the `templates/` folder
2. Link the mission title to the template file in `mission/templateMap.go`.
