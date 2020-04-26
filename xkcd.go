package main

import (
	"fmt"

	d "github.com/bwmarrin/discordgo"
	xkcd "github.com/nishanths/go-xkcd"
)

// createComicEmbed creates a pointer to a Discord MessageEmbed object with a random xkcd comic
func createComicEmbed() *d.MessageEmbed {
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

	return embed
}
