package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	d "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	parse "github.com/mattn/go-shellwords"
	xkcd "github.com/nishanths/go-xkcd"
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
		fmt.Println("Error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
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
	input, err := parse.Parse(m.Content)
	if err != nil {
		fmt.Println("Error parsing input,", err)
		return
	}

	if len(input) > 0 {
		command := input[0]
		switch command {
		// Ping
		case "!ping":
			s.ChannelMessageSend(m.ChannelID, "Pong")

		// Google Translate
		// Mostly for Andrei's weird German text
		case "!translate":
			if len(input) != 4 {
				s.ChannelMessageSend(m.ChannelID, "Invalid input for !translate. Check !help for an overview.")
				return
			}
			s.ChannelMessageSend(m.ChannelID, createTranslateMessage(input[1], input[2], input[3]))

		// Random xkcd comic
		case "!xkcd":
			// Create xkcd client and get random comic
			xkcdClient := xkcd.NewClient()
			comic, _ := xkcdClient.Random()

			// Create embed
			embed := &d.MessageEmbed{
				Author: &d.MessageEmbedAuthor{
					Name:    "xkcd",
					URL:     "https://xkcd.com",
					IconURL: "https://is2-ssl.mzstatic.com/image/thumb/Purple123/v4/b5/f6/a2/b5f6a20c-5e4e-fb72-2592-2841784bc48c/AppIcon-0-0-1x_U007emarketing-0-0-0-5-0-0-sRGB-0-0-0-GLES2_U002c0-512MB-85-220-0-0.jpeg/320x0w.jpg",
				},
				Color: 0x96aac8,
				URL:   fmt.Sprintf("https://xkcd.com/%d", comic.Number),
				Title: fmt.Sprintf("#%d - %s", comic.Number, comic.Title),
				Image: &d.MessageEmbedImage{URL: comic.ImageURL},
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)

		// Tweet command
		case "!tweet":
			s.ChannelMessageSendEmbed(m.ChannelID, createTweetEmbed(twitterAccount))

		// Help command
		case "!help":
			embed := &d.MessageEmbed{
				Author: &d.MessageEmbedAuthor{
					Name:    s.State.User.Username,
					IconURL: s.State.User.AvatarURL(""),
				},
				Title: "Commands",
				Fields: []*d.MessageEmbedField{
					&d.MessageEmbedField{
						Name:  "!help",
						Value: "Invokes this help box... but it's obvious you know that already.",
					},
					&d.MessageEmbedField{
						Name:  "!ping",
						Value: "Check if the bot is alive.",
					},
					&d.MessageEmbedField{
						Name:  "!translate [text] [input lang] [output lang]",
						Value: "Uses Google Translate. Keep text in quotations, and use two-letter language codes.",
					},
					&d.MessageEmbedField{
						Name:  "!xkcd",
						Value: "Shows a random xkcd comic.",
					},
					&d.MessageEmbedField{
						Name:  "!tweet",
						Value: "Shows a random tweet from our favourite news source.",
					},
				},
			}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		}
	}
}
