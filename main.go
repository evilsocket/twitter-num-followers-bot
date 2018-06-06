package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/yaml.v2"
)

type Keys struct {
	Username       string `yaml:"username"`
	ConsumerKey    string `yaml:"consumer_key"`
	ConsumerSecret string `yaml:"consumer_secret"`
	AccessToken    string `yaml:"access_token"`
	AccessSecret   string `yaml:"access_secret"`
}

func (k Keys) Valid() bool {
	return k.ConsumerKey != "" && k.ConsumerSecret != "" && k.AccessToken != "" && k.AccessSecret != ""
}

type Config struct {
	Keys Keys `yaml:"twitter"`
}

var (
	keysPath = flag.String("keys", "~/.tnfbot.config.yml", "Yaml file containing Twitter auth keys.")
	period   = flag.Duration("period", 5*time.Minute, "Period to wait between one check and the next one.")

	config = Config{}
	err    = (error)(nil)
	client = (*twitter.Client)(nil)
)

func expandPath(path string) (string, error) {
	// Check if path is empty
	if path != "" {
		if strings.HasPrefix(path, "~") {
			usr, err := user.Current()
			if err != nil {
				return "", err
			} else {
				// Replace only the first occurrence of ~
				path = strings.Replace(path, "~", usr.HomeDir, 1)
			}
		}
		return filepath.Abs(path)
	}
	return "", nil
}

func getFollowers() (followers []twitter.User, err error) {
	log.Printf("fetching list of followers ...")

	nextCursor := int64(0)
	followers = make([]twitter.User, 0)

	for {
		params := &twitter.FollowerListParams{
			ScreenName: config.Keys.Username,
			Count:      -1,
			Cursor:     nextCursor,
		}

		users, _, err := client.Followers.List(params)
		if err != nil {
			return nil, err
		}

		for _, u := range users.Users {
			followers = append(followers, u)
		}

		if users.NextCursor == 0 {
			break
		} else {
			nextCursor = users.NextCursor
		}
	}

	return
}

func main() {
	if *keysPath, err = expandPath(*keysPath); err != nil {
		log.Fatal(err)
	} else if data, err := ioutil.ReadFile(*keysPath); err != nil {
		log.Fatal(err)
	} else if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatal(err)
	} else if !config.Keys.Valid() {
		log.Fatalf("keys file %s is empty or not valid", *keysPath)
	}

	tconfig := oauth1.NewConfig(config.Keys.ConsumerKey, config.Keys.ConsumerSecret)
	token := oauth1.NewToken(config.Keys.AccessToken, config.Keys.AccessSecret)
	httpClient := tconfig.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)

	log.Printf("bot started with a period of %v", *period)

	for {
		followers, err := getFollowers()
		if err == nil {
			log.Printf("%v", followers)
		} else {
			log.Printf("error: %v", err)
		}

		for _, u := range followers {
			log.Printf("%s followers: %d", u.Name, u.FollowersCount)
		}

		time.Sleep(*period)
	}
}
