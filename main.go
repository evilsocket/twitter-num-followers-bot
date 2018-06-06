package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"gopkg.in/yaml.v2"
)

var (
	keysPath = flag.String("keys", "~/.tnfbot.config.yml", "Yaml file containing Twitter auth keys.")
	period   = flag.Duration("period", 5*time.Minute, "Period to wait between one check and the next one.")

	config = Config{}
	err    = (error)(nil)
	client = (*twitter.Client)(nil)
)

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
			for _, u := range followers {
				if didChange(u) {
					log.Printf("checking user %s (id=%d followers=%d) ...", u.Name, u.ID, u.FollowersCount)

					for _, check := range Checks {
						if check.Checker(u.FollowersCount) {
							if text, err := check.Text(u); err == nil {
								log.Printf("> %s", text)
								profileURL := fmt.Sprintf("https://twitter.com/%s", u.ScreenName)
								if err := takeScreenshot(profileURL); err == nil {
									log.Printf("TODO: tweet + screenshot.png")
								} else {
									log.Printf("error while taking profile screenshot: %v", err)
								}
							} else {
								log.Printf("error while creating tweet: %v", err)
							}

							break
						}
					}
				}
			}
		} else {
			log.Printf("error: %v", err)
		}

		log.Printf("sleeping for %v", *period)

		time.Sleep(*period)
	}
}
