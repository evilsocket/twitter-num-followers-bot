package main

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"text/template"

	"github.com/dghubble/go-twitter/twitter"
)

type Check struct {
	Name     string
	Checker  func(n int) bool
	Template string
}

var Checks = []Check{
	Check{
		Name: "samesame",
		Checker: func(n int) bool {
			sn := strconv.Itoa(n)
			slen := len(sn)
			return slen >= 3 && strings.Count(sn, string(sn[0])) == slen
		},
		Template: "Check this out, @{{.ScreenName}} just reached {{.FollowersCount}} followers!",
	},

	// tnx to https://www.thepolyglotdeveloper.com/2016/12/determine-number-prime-using-golang/
	Check{
		Name: "isprime",
		Checker: func(n int) bool {
			if n < 100 {
				return false
			}

			for i := 2; i <= int(math.Floor(math.Sqrt(float64(n)))); i++ {
				if n%i == 0 {
					return false
				}
			}

			return true
		},
		Template: "Check this out, @{{.ScreenName}} just reached {{.FollowersCount}} followers, which is a prime number!",
	},
}

func (c Check) Text(u twitter.User) (out string, err error) {
	parsed, err := template.New("bot").Parse(c.Template)
	if err != nil {
		return
	}

	var buff bytes.Buffer
	err = parsed.Execute(&buff, u)
	if err != nil {
		return
	}

	return buff.String(), nil
}
