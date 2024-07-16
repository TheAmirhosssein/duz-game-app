package match

import "math/rand"

func playerIcon() string {
	randNumber := rand.Int()
	if randNumber%2 == 0 {
		return "X"
	} else {
		return "O"
	}
}
