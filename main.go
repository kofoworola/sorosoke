package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	// twitter api credentials
	apiKey           = os.Getenv("API_KEY")
	apiSecret        = os.Getenv("API_SECRET")
	oauthToken       = os.Getenv("OAUTH_TOKEN")
	oauthTokenSecret = os.Getenv("OAUTH_TOKEN_SECRET")

	// contentFile is a json file with a list of text content that will be sent alongside each countdown.
	// it will be picked at random
	contentFile = os.Getenv("TWEET_CONTENT")

	emptyBlock   rune
	checkedBlock rune
)

const (
	parseFormat = "02/01/2006 15:04"
	// startDate is the starting reference date used to generate the percentage counter
	startDay = "17/10/2020 12:00"
	// electionDay is February 18 2023 as announced by INEC
	electionDay   = "18/02/2023 12:00"
	maxChars      = 280
	handle        = "sorosokeywere"
	twitterAPIURL = "https://api.twitter.com/1.1/statuses/update.json"
)

func init() {
	// this byte arrangement respresents the `░` symbol
	emptyBlock, _ = utf8.DecodeRune([]byte{226, 150, 145})
	// this byte arrangement represents the `▓` symbol
	checkedBlock, _ = utf8.DecodeRune([]byte{226, 150, 147})
}

func main() {
	// validation
	switch "" {
	case apiKey:
		log.Fatal("API_KEY environment variable not set")
	case apiSecret:
		log.Fatal("API_SECRET environment variable not set")
	case oauthToken:
		log.Fatal("OAUTH_TOKEN environment variable not set")
	case oauthTokenSecret:
		log.Fatalf("OAUTH_TOKEN_SECRET environment variable not set")
	}
	var b strings.Builder
	generateReminder(&b)
	b.WriteRune('\n')
	readContent(&b)
	tweets := splitTweets(b.String(), maxChars)
	if err := TweetAll(tweets); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	fmt.Println("Done")

}

func generateReminder(builder *strings.Builder) {
	electionDate, err := time.Parse(parseFormat, electionDay)
	if err != nil {
		log.Fatalf("error parising election date: %v", err)
	}
	startDate, err := time.Parse(parseFormat, startDay)
	if err != nil {
		log.Fatalf("error parsing start date: %v", err)
	}

	daysFromStart := int(electionDate.Sub(startDate).Hours() / 24)
	daysFromToday := int(electionDate.Sub(time.Now()).Hours() / 24)
	builder.WriteString(fmt.Sprintf("Remember there are %d day(s) to the 2023 election \n", daysFromToday))

	difference := daysFromStart - daysFromToday
	percentAway := (float32(difference) / float32(daysFromStart)) * 100
	numberChecked := int(percentAway) / 5
	for i := 0; i < numberChecked; i++ {
		builder.WriteRune(checkedBlock)
	}
	for i := 0; i < 20-numberChecked; i++ {
		builder.WriteRune(emptyBlock)
	}
}

func readContent(builder *strings.Builder) {
	if contentFile == "" {
		contentFile = "./content.json"
	}

	file, err := os.OpenFile(contentFile, os.O_RDONLY, 775)
	if err != nil {
		log.Fatalf("error opening file %s: %v", contentFile, err)
	}

	var decodedContent []string
	if err := json.NewDecoder(file).Decode(&decodedContent); err != nil {
		log.Fatalf("error decoding content: %v", err)
	}
	rand.Seed(time.Now().UnixNano())
	itemToUse := decodedContent[rand.Intn(len(decodedContent))]

	builder.WriteString(itemToUse)
}

// since some characters count as two characters on twitter,
// from my checks, my assumption is if it is more than 2 bytes it is counted as two characters.
// so we need to split sequentially to support all characters available.
func splitTweets(content string, tweetLength int) []string {
	tweets := make([]string, 0)
	count := 0
	byteCount := 0
	runes := []rune(content)
	lastSpaceIndex := 0
	for _, r := range runes {
		bytes := make([]byte, 4)
		rLen := utf8.EncodeRune(bytes, r)
		if rLen > 2 {
			count += 2
		} else {
			count++
		}
		if r == 32 || r == 10 {
			lastSpaceIndex = byteCount
		}
		byteCount += rLen
		if count >= tweetLength {
			if lastSpaceIndex == 0 {
				tweets = append(tweets, content)
				return tweets
			} else {
				tweets = append(tweets, content[:lastSpaceIndex])
				remainder := content[lastSpaceIndex+1:]
				tweets = append(tweets, splitTweets(remainder, tweetLength)...)
				return tweets
			}
		}
	}
	// less than tweet length
	tweets = append(tweets, content)
	return tweets

}
