package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	d "github.com/bwmarrin/discordgo"
	t "github.com/n0madic/twitter-scraper"
)

// Number of pages to look back through (20 tweets per page)
// The higher the number, the slower responses become
const pages = 5

func init() {
	// Random seed
	rand.Seed(time.Now().UnixNano())
}

// getTweet returns a Tweet object from the twitter-scraper library
func getTweet(account string) *t.Result {
	// Choose a random tweet from within the page limit
	var chosenTweet *t.Result
	limit := rand.Intn(pages * 20)

	count := 0
	// Iterate over each tweet from the channel
	for tweet := range t.GetTweets(account, pages) {
		if tweet.Error != nil {
			panic(tweet.Error)
		}

		if count == limit {
			chosenTweet = tweet
			break
		}
		count++
	}

	return chosenTweet
}

// getFormattedTweet formats useful data from Tweet objects into separate variables
func getFormattedTweet(account string) (text, url, photo, timestamp string, retweets, likes int) {
	tweet := getTweet(account)
	text = tweet.Text

	// Remove extra urls linking to retweets
	extraURLs := regexp.MustCompile(`((http|https):\/\/[\w\-]+(\.[\w\-]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?)`)
	text = extraURLs.ReplaceAllString(text, "\n\n")

	// Remove pic.twitter.com urls
	tweetURL := regexp.MustCompile(`(pic.twitter.com)(.*)`)
	text = tweetURL.ReplaceAllString(text, "")

	// Add markdown links to hashtags
	hashtag := regexp.MustCompile(`(^|)#([A-Za-z_][A-Za-z0-9_]*)`)
	formatHashtagLink := func(match string) string {
		return fmt.Sprintf("[%s](https://twitter.com/hashtag/%s)", match, match[1:])
	}
	text = hashtag.ReplaceAllStringFunc(text, formatHashtagLink)

	// Add markdown links to user mentions
	mention := regexp.MustCompile(`(^|)@([A-Za-z_][A-Za-z0-9_]*)`)
	formatMentionLink := func(match string) string {
		return fmt.Sprintf("[%s](https://twitter.com/%s)", match, match[1:])
	}
	text = mention.ReplaceAllStringFunc(text, formatMentionLink)

	// Format tweet URL
	url = fmt.Sprintf("https://twitter.com/%s/status/%s", account, tweet.ID)

	// Format tweet timestamp
	timestamp = time.Unix(tweet.Timestamp, 0).Format("2006-01-02T15:04:05.000Z")

	// Get photo URL
	// Video thumbnail is there is a video
	// First photo is there are photos
	// Empty string is returned if there are no photos
	if len(tweet.Videos) > 0 {
		photo = tweet.Videos[0].Preview
	} else if len(tweet.Photos) > 0 {
		photo = tweet.Photos[0]
	}

	return text, url, photo, timestamp, tweet.Retweets, tweet.Likes
}

// getFormattedProfile formats useful data from Profile objects from the twitter-scraper library
func getFormattedProfile(account string) (name, avatar string) {
	profile, err := t.GetProfile(account)
	if err != nil {
		panic(err)
	}

	return profile.Name, profile.Avatar
}

// createTweetEmbed creates a pointer to a Discord MessageEmbed object with the tweet data
func createTweetEmbed(account string) *d.MessageEmbed {
	tweet, url, photo, timestamp, retweets, likes := getFormattedTweet(account)
	profileName, profilePhoto := getFormattedProfile(account)

	embed := &d.MessageEmbed{
		Author: &d.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s (@%s)", profileName, account),
			URL:     url,
			IconURL: profilePhoto,
		},
		Color:       0x00acee,
		Description: tweet,
		Image:       &d.MessageEmbedImage{URL: photo},
		Footer: &d.MessageEmbedFooter{
			Text:    fmt.Sprintf("%d Retweets | %d Likes", retweets, likes),
			IconURL: "https://cdn2.iconfinder.com/data/icons/minimalism/512/twitter.png",
		},
		Timestamp: timestamp,
	}

	return embed
}
