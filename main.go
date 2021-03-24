package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const altText = "The image is a meme featuring Tintin and Captain Haddock. With an exhausted body language, Captain Haddock comments to Tintin: “What a week, huh?” To which Tintin replies, “Captain, it’s Wednesday!”"

//go:embed img.jpg
var img []byte

func getClient() (*twitter.Client, error) {
	ck := os.Getenv("WEDNESDAY_CONSUMER_KEY")
	cs := os.Getenv("WEDNESDAY_CONSUMER_SECRET")
	at := os.Getenv("WEDNESDAY_ACCESS_TOKEN")
	ats := os.Getenv("WEDNESDAY_ACCESS_TOKEN_SECRET")

	errs := []string{}
	if ck == "" {
		errs = append(errs, "WEDNESDAY_CONSUMER_KEY")
	}
	if cs == "" {
		errs = append(errs, "WEDNESDAY_CONSUMER_SECRET")
	}
	if at == "" {
		errs = append(errs, "WEDNESDAY_ACCESS_TOKEN")
	}
	if ats == "" {
		errs = append(errs, "WEDNESDAY_ACCESS_TOKEN_SECRET")
	}
	if len(errs) != 0 {
		return nil, fmt.Errorf("Missing environment variable(s): %s", strings.Join(errs, "\n\t"))
	}

	config := oauth1.NewConfig(ck, cs)
	token := oauth1.NewToken(at, ats)
	return twitter.NewClient(config.Client(oauth1.NoContext, token)), nil
}

func isWednesday() bool {
	return time.Now().Weekday() == time.Wednesday
}

func shouldTweet() bool {
	r := bufio.NewReader(os.Stdin)
	fmt.Println("It's Wednesday. Should we tweet the meme? [y/n]")
	i, err := r.ReadString('\n')
	if err != nil {
		os.Exit(1)
	}
	return i == "y\n"
}

func main() {
	if !isWednesday() || !shouldTweet() {
		return
	}

	c, err := getClient()
	if err != nil {
		log.Fatal(err)
	}

	m, _, err := c.Media.Upload(img, "tweet_image")
	if err != nil {
		log.Fatal(err)
	}

	t, _, err := c.Statuses.Update("", &twitter.StatusUpdateParams{MediaIds: []int64{m.MediaID}})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("https://twitter.com/%s/status/%d\n", t.User.ScreenName, t.ID)
}
