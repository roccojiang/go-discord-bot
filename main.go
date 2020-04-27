package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	d "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	p "github.com/mattn/go-shellwords"
)

var (
	botToken       string
	twitterAccount string
)

// Load .env file and variables
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken = os.Getenv("BOT_SECRET")
	twitterAccount = os.Getenv("TWITTER_ACCOUNT")
}

func main() {
	// Create a new Discord session using the provided bot token
	dg, err := d.New("Bot " + botToken)
	if err != nil {
		log.Println("Error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		log.Println("Error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received
	fmt.Printf("Bot is running since %s. Press CTRL-C to exit.\n", time.Now().Format(time.ANSIC))
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to
func messageCreate(s *d.Session, m *d.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Parse commands
	input, err := p.Parse(m.Content)
	if err != nil {
		log.Println("Error parsing input,", err)
		return
	}

	// Commands
	if len(input) > 0 {
		command := input[0]
		commandNum := len(input)
		switch command {
		// Ping
		case "!ping":
			s.ChannelMessageSend(m.ChannelID, "I'm alive... *sadly*...")

		// Get tweet from Twitter
		case "!tweet":
			switch commandNum {
			// Default Twitter account
			case 1:
				s.ChannelMessageSendEmbed(m.ChannelID, createTweetEmbed(twitterAccount))
			// Custom specified Twitter account
			case 2:
				customTwitterAccount := input[1]
				s.ChannelMessageSendEmbed(m.ChannelID, createTweetEmbed(customTwitterAccount))
			default:
				s.ChannelMessageSend(m.ChannelID, "Invalid input for !tweet. Check !help for an overview.")
				return
			}

		// Google Translate
		// Mostly for Andrei's weird German stuff
		case "!translate":
			switch commandNum {
			case 4:
				s.ChannelMessageSend(m.ChannelID, createTranslateMessage(input[1], input[2], input[3]))
			default:
				s.ChannelMessageSend(m.ChannelID, "Invalid input for !translate. Check !help for an overview.")
				return
			}

		// Get xkcd comic
		case "!xkcd":
			switch commandNum {
			// Random comic
			case 1:
				s.ChannelMessageSendEmbed(m.ChannelID, createComicEmbed(getRandomComic()))
			// Specified comic number
			case 2:
				comicNum, err := strconv.Atoi(input[1])
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Invalid input for !xkcd. Check !help for an overview.")
					return
				}
				s.ChannelMessageSendEmbed(m.ChannelID, createComicEmbed(getComic(comicNum)))
			default:
				s.ChannelMessageSend(m.ChannelID, "Invalid input for !xkcd. Check !help for an overview.")
				return
			}

		// Help command
		case "!help":
			s.ChannelMessageSendEmbed(m.ChannelID, createHelpEmbed(s))
		}
	}
}
