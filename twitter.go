package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dghubble/oauth1"
)

type Tweet struct {
	ID     string `json:"id_str"`
	Handle string `json:"screen_name"`
}

func TweetAll(tweets []string) error {
	config := oauth1.NewConfig(apiKey, apiSecret)
	token := oauth1.NewToken(oauthToken, oauthTokenSecret)

	client := config.Client(context.Background(), token)
	var err error
	lastSent := &Tweet{}
	for _, t := range tweets {
		lastSent, err = SendTweet(t, lastSent.ID, lastSent.Handle, client)
		if err != nil {
			return err
		}
	}
	return nil
}

func SendTweet(content, reply_id, replyHandle string, client *http.Client) (*Tweet, error) {
	if reply_id != "" {
		content = fmt.Sprintf("@%s %s", replyHandle, content)
	}
	values := url.Values{
		"status":                {content},
		"in_reply_to_status_id": {reply_id},
	}
	resp, err := client.PostForm(twitterAPIURL, values)
	if err != nil {
		return nil, fmt.Errorf("error sending tweet: %w", err)
	}
	var tweet Tweet
	if err := json.NewDecoder(resp.Body).Decode(&tweet); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return &tweet, nil
}
