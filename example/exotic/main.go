package main

import (
	"fmt"
	"log"

	"github.com/eljamo/weightedoption/v3"
)

// Simulates 100 chances for dropping a raid exotic weapon from Destiny, which has a 5% drop chance when a player completes the raid
func main() {
	// Create a new selector with options and their weights
	s, err := weightedoption.NewSelector(
		weightedoption.NewOption('🔫', 5),  // 5% chance for the exotic weapon
		weightedoption.NewOption('❌', 95), // 95% chance for no drop
	)
	if err != nil {
		log.Fatal(err)
	}

	chances := make([]rune, 100) // Array to store the results of 100 attempts
	for i := 0; i < len(chances); i++ {
		chances[i] = s.Select() // Select an option based on their weights
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
