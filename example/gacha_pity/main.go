package main

import (
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/eljamo/weightedoption/v3"
)

func oneOrTen() int {
	if rand.IntN(2) == 0 {
		return 1
	}
	return 10
}

type GachaBanner struct {
	// The pool of items that can be dropped
	pool []weightedoption.Option[string, float64]
	// The pity threshold for the banner
	pityThreshold int
	// The item that will be dropped when the pity threshold is reached
	pityDrop string
	// The pity counter for the banner
	pityCounterMap map[string]int
	// selector
	selector *weightedoption.Selector[string, float64]
}

func (b *GachaBanner) pull(userId string) string {
	// Pull an item from the selector
	drop := b.selector.Select()
	pityCount := b.pityCounterMap[userId] + 1

	if pityCount >= b.pityThreshold {
		// Reset the pity counter
		b.pityCounterMap[userId] = 0
		return b.pityDrop
	}

	// Update the pity counter
	b.pityCounterMap[userId] = pityCount

	return drop
}

func (b *GachaBanner) PullN(n int, userId string) []string {
	drops := make([]string, n)

	for i := 0; i < n; i++ {
		drop := b.pull(userId)
		drops[i] = drop
	}

	return drops

}

// Simulates a mobile game gacha pull which drops 10 items and they have floating point weights
func main() {
	userId := "123456"
	pityThreshold := 90
	pityDrop := "5★ Character"
	pityCounterMap := make(map[string]int)

	pool := []weightedoption.Option[string, float64]{
		{Data: pityDrop, Weight: 0.6},
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

	banner := GachaBanner{
		pool:           pool,
		pityThreshold:  pityThreshold,
		pityDrop:       pityDrop,
		pityCounterMap: pityCounterMap,
		selector:       s,
	}

	count := 0
	pityPulled := false
	allDrops := make([]string, 0)
	tally := make(map[string]int)

	// Run until the pity drop is pulled, om the 90th it'll be guaranteed
	for !pityPulled {
		timesToPull := oneOrTen()
		drops := banner.PullN(timesToPull, userId)

		for _, drop := range drops {
			allDrops = append(allDrops, drop)

			count++

			if drop == pityDrop {
				count = 0
				pityPulled = true
			}
		}

		if pityPulled {
			break
		}
	}

	for _, drop := range allDrops {
		tally[drop]++
	}

	// Print the tally for each item
	fmt.Println("Tally:")
	for item, count := range tally {
		fmt.Printf("%s: %d\n", item, count)
	}
}
