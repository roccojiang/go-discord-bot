package main

import d "github.com/bwmarrin/discordgo"

// createComicEmbed creates a pointer to a Discord MessageEmbed object with a help box
func createHelpEmbed(s *d.Session) *d.MessageEmbed {
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
				Name:  "!tweet",
				Value: "Shows a random tweet from our favourite news source.",
			},
			&d.MessageEmbedField{
				Name:  "!translate [text] [input lang] [output lang]",
				Value: "Uses Google Translate. Keep text in quotations, and use two-letter language codes.",
			},
			&d.MessageEmbedField{
				Name:  "!xkcd [num]",
				Value: "Shows an xkcd comic. If the comic number is not provided, a random one will be chosen.",
			},
		},
	}

	return embed
}
