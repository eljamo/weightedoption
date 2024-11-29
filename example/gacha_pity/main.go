package main

import (
	"fmt"
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
	pool           []weightedoption.Option[string, float64]
	pityThreshold  int
	pityDrop       string
	pityCounterMap map[string]int
	selector       *weightedoption.Selector[string, float64]
}

func (b *GachaBanner) pull(userId string) string {
	// Pull an item from the selector
	drop := b.selector.Select()
	pityCount := b.pityCounterMap[userId] + 1

	if pityCount >= b.pityThreshold {
		fmt.Println("Pity drop after reaching threshold")
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

func NewGachaBanner(pool []weightedoption.Option[string, float64], pityThreshold int, pityDrop string) (*GachaBanner, error) {
	s, err := weightedoption.NewSelector(pool...)
	if err != nil {
		return nil, err
	}

	return &GachaBanner{
		pool:           pool,
		pityThreshold:  pityThreshold,
		pityDrop:       pityDrop,
		pityCounterMap: make(map[string]int),
		selector:       s,
	}, nil
}

func generateOptions(mainItem, secondaryItem, tertiaryItem string, remainingItems []string) []weightedoption.Option[string, float64] {
	const totalWeight = 100.0
	weights := map[string]float64{
		"main":      0.6,
		"secondary": 3.3,
		"tertiary":  1.77,
	}

	remainingWeight := (totalWeight - weights["main"] - weights["secondary"] - weights["tertiary"]) / float64(len(remainingItems))

	options := []weightedoption.Option[string, float64]{
		{Data: mainItem, Weight: weights["main"]},
		{Data: secondaryItem, Weight: weights["secondary"]},
		{Data: tertiaryItem, Weight: weights["tertiary"]},
	}

	for _, item := range remainingItems {
		options = append(options, weightedoption.Option[string, float64]{Data: item, Weight: remainingWeight})
	}

	return options
}

// Simulates a mobile game gacha pull which drops 10 items and they have floating point weights
func main() {
	userId := "123456"
	pityThreshold := 90
	pityDrop := "5★ Character"

	pool := generateOptions("5★ Character", "4★ Character", "4★ Weapon (Sword)", []string{
		"3★ Weapon (Sword)",
		"3★ Weapon (Polearm)",
		"3★ Weapon (Bow)",
		"3★ Weapon (Claymore)",
		"3★ Weapon (Staff)",
	})

	banner, err := NewGachaBanner(pool, pityThreshold, pityDrop)
	if err != nil {
		panic(err)
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

	fmt.Println("Tally:")
	for item, count := range tally {
		fmt.Printf("%s: %d\n", item, count)
	}
}
