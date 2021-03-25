package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"
	"time"
)

const cacheFile = ".its-wednesday"
const cacheFormat = "2006-01-02 15:04:05"

type cache struct {
	lastAttempt time.Time
	lastTweet   time.Time
	now         time.Time
	loaded      bool
}

var value cache

func (c *cache) path() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Could not detect current user")
	}
	return strings.Join([]string{usr.HomeDir, cacheFile}, string(os.PathSeparator))
}

func (c *cache) bytes() []byte {
	return []byte(fmt.Sprintf("%s,%s", c.lastAttempt.Format(cacheFormat), c.lastTweet.Format(cacheFormat)))
}

func (c *cache) write() {
	err := os.WriteFile(c.path(), c.bytes(), 0666)
	if err != nil {
		log.Fatal("Could not write cache file")
	}
}

func (c *cache) tweetedToday() bool {
	return (c.lastTweet.Year() == c.now.Year() &&
		c.lastTweet.Month() == c.now.Month() &&
		c.lastTweet.Day() == c.now.Day())
}

func (c *cache) triedToday() bool {
	return (c.lastAttempt.Year() == c.now.Year() &&
		c.lastAttempt.Month() == c.now.Month() &&
		c.lastAttempt.Day() == c.now.Day())
}

func (c *cache) triedThisMorning() bool {
	return c.triedToday() && c.lastAttempt.Hour() >= 8 && c.lastAttempt.Hour() < 12
}

func (c *cache) triedThisAfternoon() bool {
	return c.triedToday() && c.lastAttempt.Hour() >= 12 && c.lastAttempt.Hour() < 18
}

func (c *cache) triedThisEvening() bool {
	return c.triedToday() && c.lastAttempt.Hour() >= 18
}

func loadCache() cache {
	if value.loaded {
		return value
	}

	value = cache{now: time.Now()}
	b, err := os.ReadFile(cacheFile)
	if err != nil {
		return value
	}
	vs := strings.Split(string(b), ",")

	a, err := time.Parse(cacheFormat, vs[0])
	if err != nil {
		return value
	}

	t, err := time.Parse(cacheFormat, vs[1])
	if err != nil {
		return value
	}

	value.lastAttempt = a
	value.lastTweet = t
	value.loaded = true
	return value
}
