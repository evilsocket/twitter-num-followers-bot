package main

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
