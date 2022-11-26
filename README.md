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

### Option 1: Download the Packaged Binary

You can download the latest release from <https://github.com/Un4gi/mBot/releases> (This is by far the quickest option if you don't want to deal with installing Go).

### Option 2: Install Using Go

Installing with Go is simple:

```bash
go get github.com/un4gi/mBot

# EDIT: This one doesn't seem to work at the moment... use the one above.
# go install github.com/un4gi/mBot@latest
```

### Option 3: Building From Source

If you prefer, you can build mBot straight from the source directory:

```bash
git clone https://github.com/Un4gi/mBot.git
cd mBot
go build .
```

## Usage

Example:

```bash
mBot -t auth_token_here -d 30
```

For help, use the `-h` flag:

```bash
mBot -h
```

| Flag | Description | Example |
|------|-------------|---------|
| `-d` | Specifies a delay between polls (to please Synack) | `mBot -d 30` |
| `-t` | Passes your Authorization: Bearer token to the bot | `mBot -t <auth_token_here>` |

## Mission Templates

Good news... mBot now allows auto-population of the Intro/Testing Methodology/Conclusion fields for each claimed mission!

*That's great... but how do I do it?*

__It's simple, really!__

1. Store your templates as JSON files in the `templates/` folder
2. Link the mission title to the template file in `mission/templateMap.go`.
