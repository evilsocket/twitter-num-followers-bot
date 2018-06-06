package main

import (
	"log"
	"os/user"
	"path/filepath"
	"strings"
	"sync"

	"github.com/dghubble/go-twitter/twitter"
)

var (
	cache = make(map[int64]int)
	lock  = sync.Mutex{}
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

func didChange(u twitter.User) bool {
	lock.Lock()
	defer lock.Unlock()

	if prev, found := cache[u.ID]; found {
		cache[u.ID] = u.FollowersCount
		return prev != u.FollowersCount
	} else {
		cache[u.ID] = u.FollowersCount
		return true
	}
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
