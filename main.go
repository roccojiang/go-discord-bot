package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	gt "github.com/bas24/googletranslatefree"
	dg "github.com/bwmarrin/discordgo"
	parse "github.com/mattn/go-shellwords"
	xkcd "github.com/nishanths/go-xkcd"
)

// Initialise token variable
var (
	Token = os.Getenv("BOT_SECRET")
)

func main() {
	// Create a new Discord session using the provided bot token
	dg, err := dg.New("Bot " + Token)
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
func messageCreate(s *dg.Session, m *dg.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Parse commands
	input, err := parse.Parse(m.Content)
	if err != nil {
		fmt.Println("Error parsing input,", err)
		return
	}

	switch input[0] {
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

		text := input[1]
		in_lang := input[2]
		out_lang := input[3]

		translated, _ := gt.Translate(text, in_lang, out_lang)
		s.ChannelMessageSend(m.ChannelID, translated)

	// Random xkcd comic
	case "!xkcd":
		// Create xkcd client and get random comic
		xkcdClient := xkcd.NewClient()
		comic, _ := xkcdClient.Random()

		// Create embed
		embed := &dg.MessageEmbed{
			Author: &dg.MessageEmbedAuthor{
				Name:    "xkcd",
				URL:     "https://xkcd.com",
				IconURL: "https://is2-ssl.mzstatic.com/image/thumb/Purple123/v4/b5/f6/a2/b5f6a20c-5e4e-fb72-2592-2841784bc48c/AppIcon-0-0-1x_U007emarketing-0-0-0-5-0-0-sRGB-0-0-0-GLES2_U002c0-512MB-85-220-0-0.jpeg/320x0w.jpg",
			},
			Color: 0x96aac8,
			URL:   fmt.Sprintf("https://xkcd.com/%d", comic.Number),
			Title: fmt.Sprintf("#%d - %s", comic.Number, comic.Title),
			Image: &dg.MessageEmbedImage{URL: comic.ImageURL},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)

	// Help command
	case "!help":
		embed := &dg.MessageEmbed{
			Author: &dg.MessageEmbedAuthor{
				Name:    s.State.User.Username,
				IconURL: s.State.User.AvatarURL(""),
			},
			Title: "Commands",
			Fields: []*dg.MessageEmbedField{
				&dg.MessageEmbedField{
					Name:  "!help",
					Value: "Invokes this help box... but it's obvious you know that already.",
				},
				&dg.MessageEmbedField{
					Name:  "!ping",
					Value: "Check if the bot is alive.",
				},
				&dg.MessageEmbedField{
					Name:  "!translate [text] [input lang] [output lang]",
					Value: "Uses Google Translate. Keep text in quotations, and use two-letter language codes.",
				},
				&dg.MessageEmbedField{
					Name:  "!xkcd",
					Value: "Shows a random xkcd comic.",
				},
			},
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}
