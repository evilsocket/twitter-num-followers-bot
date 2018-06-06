package main

import (
	"math"
	"strconv"
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
			if len(sn) < 3 {
				return false
			}

			first := rune(sn[0])
			for _, c := range sn {
				if c != first {
					return false
				}
			}

			return true
		},
		Template: "Check this out, @{username} just reached {followers} followers!",
	},

	// tnx to https://www.thepolyglotdeveloper.com/2016/12/determine-number-prime-using-golang/
	Check{
		Name: "isprime",
		Checker: func(n int) bool {
			sn := strconv.Itoa(n)
			if len(sn) < 3 {
				return false
			}

			for i := 2; i <= int(math.Floor(math.Sqrt(float64(n)))); i++ {
				if n%i == 0 {
					return false
				}
			}

			return n > 1
		},
		Template: "Check this out, @{username} just reached {followers} followers, which is a prime number!",
	},
}
