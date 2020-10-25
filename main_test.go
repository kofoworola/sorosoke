package main

import (
	"encoding/json"
	"os"
	"testing"
)

// this test is really just to make sure the split happens at a space on or before the breakpoint,
// and that the block character (and others like emoji) are counted as two characters,
// the split point of 76 is because the 76th character is the block.
func TestSplitTweets(t *testing.T) {
	dummyText := `Lorem Ipsum is simply dummy text of the printing and typesetting industry░. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen.`
	split := splitTweets(dummyText, 76)
	if len(split) != 4 {
		t.Fatalf("expected len 4 got %d", len(split))
	}
	if firstChar := string(split[1][0]); firstChar != "i" {
		t.Fatalf("expected first letter of second sentence to begin with i got %s", firstChar)
	}
}

func TestJsonContent(t *testing.T) {
	contentFile := os.Getenv("TWEET_CONTENT")
	if contentFile == "" {
		contentFile = "./content.json"
	}

	file, err := os.OpenFile(contentFile, os.O_RDONLY, 775)
	if err != nil {
		t.Fatalf("error opening file %s: %v", contentFile, err)
	}

	var tweets []string
	if err := json.NewDecoder(file).Decode(&tweets); err != nil {
		t.Fatalf("error decoding json: %v", err)
	}
}
