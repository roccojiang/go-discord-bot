package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	gt "github.com/bas24/googletranslatefree"
	dg "github.com/bwmarrin/discordgo"
	xkcd "github.com/nishanths/go-xkcd"
)

// Initialise token variable
var (
	Token = os.Getenv("BOT_SECRET")
)

func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := dg.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *dg.Session, m *dg.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	// Translate Andrei's weird German text
	if strings.HasPrefix(m.Content, "!whatisandreisaying") {
		en := strings.TrimPrefix(m.Content, "!whatisandreisaying ")
		de, _ := gt.Translate(en, "de", "en")
		s.ChannelMessageSend(m.ChannelID, de)
	}

	if strings.HasPrefix(m.Content, "!xkcd") {
		xkcdClient := xkcd.NewClient()
		comic, _ := xkcdClient.Random()

		img := &dg.MessageEmbedImage{URL: comic.ImageURL}
		author := &dg.MessageEmbedAuthor{
			URL:     "https://xkcd.com",
			Name:    "xkcd",
			IconURL: "https://is2-ssl.mzstatic.com/image/thumb/Purple123/v4/b5/f6/a2/b5f6a20c-5e4e-fb72-2592-2841784bc48c/AppIcon-0-0-1x_U007emarketing-0-0-0-5-0-0-sRGB-0-0-0-GLES2_U002c0-512MB-85-220-0-0.jpeg/320x0w.jpg",
		}

		msg := &dg.MessageEmbed{
			Author: author,
			Color:  0x96aac8,
			URL:    fmt.Sprintf("https://xkcd.com/%d", comic.Number),
			Title:  fmt.Sprintf("#%d - %s", comic.Number, comic.Title),
			Image:  img,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, msg)
	}
}
