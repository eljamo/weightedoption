package main

import (
	"fmt"
	"log"

	"github.com/eljamo/weightedoption/v2"
)

// Simulates 100 chances for dropping a raid exotic weapon from a Destiny which has a 5% drop chance when a player completes the raid
func main() {
	s, err := weightedoption.NewSelector(
		weightedoption.NewOption('🔫', 5),
		weightedoption.NewOption('❌', 95),
	)
	if err != nil {
		log.Fatal(err)
	}

	chances := make([]rune, 100)
	for i := 0; i < len(chances); i++ {
		chances[i] = s.Select()
	}
	fmt.Println(string(chances))

	tally := make(map[rune]int)
	for _, c := range chances {
		tally[c]++
	}

	_, err = fmt.Printf("\n🔫: %d\t❌ %d\n", tally['🔫'], tally['❌'])
	if err != nil {
		log.Fatal(err)
	}
}
