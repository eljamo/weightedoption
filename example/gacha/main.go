package main

import (
	"fmt"
	"log"

	"github.com/eljamo/weightedoption/v3"
)

// Simulates a mobile game gacha pull which drops 8 items and they have floating point weights
func main() {
	pool := []weightedoption.Option[string, float64]{
		{Data: "5★ Character", Weight: 0.6},
		{Data: "4★ Character", Weight: 3.3},
		{Data: "4★ Weapon (Sword)", Weight: 1.77},
		{Data: "3★ Weapon (Sword)", Weight: 18.86},
		{Data: "3★ Weapon (Polearm)", Weight: 18.86},
		{Data: "3★ Weapon (Bow)", Weight: 18.86},
		{Data: "3★ Weapon (Claymore)", Weight: 18.86},
		{Data: "3★ Weapon (Staff)", Weight: 18.86},
	}

	// Create a new selector with options and their weights
	s, err := weightedoption.NewSelector(pool...)
	if err != nil {
		log.Fatal(err)
	}

	tally := make(map[string]int) // Tally map to count occurrences
	drops := make([]string, 8)    // Array to store the results of the gacha pull
	for i := 0; i < len(drops); i++ {
		result := s.Select() // Select an option based on their weights
		drops[i] = result
		tally[result]++
	}

	// Print the tally for each item
	fmt.Println("Tally:")
	for item, count := range tally {
		fmt.Printf("%s: %d\n", item, count)
	}
}
