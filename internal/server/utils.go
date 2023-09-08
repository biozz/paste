package server

import (
	"math/rand"
	"strings"
	"time"
)

var (
	consonants = []string{"b", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "s", "t", "v", "w", "z"}
	vowels     = []string{"a", "e", "i", "o", "u", "r"}
	random     = rand.New(rand.NewSource(time.Now().Unix()))
)

func generateName(length int) string {
	var word []string
	for j := 0; j < length; j++ {
		var letter string
		if j%2 != 0 {
			letter = vowels[random.Intn(len(vowels))]
		} else {
			letter = consonants[random.Intn(len(consonants))]
		}

		word = append(word, letter)
	}
	return strings.Join(word, "")
}
