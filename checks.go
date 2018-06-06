package main

import (
	"bytes"
	"math"
	"strconv"
	"strings"
	"text/template"

	"github.com/dghubble/go-twitter/twitter"
)

const (
	prefix          = "✨ ACHIEVEMENT UNLOCKED!✨"
	genericTemplate = prefix + " @{{.ScreenName}} just reached {{.FollowersCount}} followers!"
)

type Check struct {
	Name     string
	Checker  func(n int) bool
	Template string
}

var Checks = []Check{
	Check{
		Name: "debug",
		Checker: func(n int) bool {
			return true
		},
		Template: "[DEBUG] " + genericTemplate + " [/DEBUG]",
	},
	Check{
		Name: "samesame",
		Checker: func(n int) bool {
			sn := strconv.Itoa(n)
			slen := len(sn)
			return slen >= 3 && strings.Count(sn, string(sn[0])) == slen
		},
		Template: genericTemplate,
	},

	// tnx to https://www.thepolyglotdeveloper.com/2016/12/determine-number-prime-using-golang/
	Check{
		Name: "prime",
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
		Template: genericTemplate + " (Neat...a prime number!)",
	},

	Check{
		Name: "pow2",
		Checker: func(n int) bool {
			if n < 100 {
				return false
			}
			return (n & (n - 1)) == 0
		},
		Template: genericTemplate + " (Neat...a power of 2!)",
	},

	Check{
		Name: "1337",
		Checker: func(n int) bool {
			// TODO: load more 1337 words from precomputed vocabulary
			return n == 1337
		},
		Template: genericTemplate,
	},

	Check{
		Name: "3zeros",
		Checker: func(n int) bool {
			return n%1000 == 0
		},
		Template: genericTemplate,
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
