package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/eljamo/weightedoption/v3"
)

type Player struct {
	RollModifierForColumnThree int
	RollModifierForColumnFour  int
}

type Config struct {
	Length        int
	NumberOfRolls int
}

const (
	MIN_OPTIONS = 2
	MAX_OPTIONS = 12
	MIN_ROLLS   = 1

	NUM_OF_SELECTORS = MAX_OPTIONS - MIN_OPTIONS
)

func generateEquallyWeightedSelectors() ([]*weightedoption.Selector[int, int], error) {
	selectors := make([]*weightedoption.Selector[int, int], 0, NUM_OF_SELECTORS)
	optionSlices := make([][]weightedoption.Option[int, int], 0, NUM_OF_SELECTORS)

	for i := MIN_OPTIONS; i <= MAX_OPTIONS; i++ {
		options := make([]weightedoption.Option[int, int], i)
		for j := 0; j < i; j++ {
			// All options are equally weighted with a value of 1
			options[j] = weightedoption.NewOption(j, 1)
		}

		optionSlices = append(optionSlices, options)
	}

	for _, options := range optionSlices {
		selector, err := weightedoption.NewSelector(options...)
		if err != nil {
			return nil, err
		}

		selectors = append(selectors, selector)
	}

	return selectors, nil
}

func getIndexes(selectors []*weightedoption.Selector[int, int], selections []Config) ([][]int, error) {
	perkIndexes := make([][]int, len(selections))
	for i, config := range selections {
		if config.Length < MIN_OPTIONS {
			return nil, fmt.Errorf("length value (%d) must be equal or greater than to %d", config.Length, MIN_OPTIONS)
		}

		if config.Length > MAX_OPTIONS {
			return nil, fmt.Errorf("length value (%d) must be equal or less than %d", config.Length, MAX_OPTIONS)
		}

		perks, err := selectIndexes(selectors, config.Length, config.NumberOfRolls)
		if err != nil {
			return nil, err
		}
		perkIndexes[i] = perks
	}

	return perkIndexes, nil
}

func selectIndexes(selectors []*weightedoption.Selector[int, int], columnLength, numRolls int) ([]int, error) {
	if numRolls < 1 {
		return nil, fmt.Errorf("numRolls value (%d) must be equal or greater than %d", numRolls, MIN_ROLLS)
	}

	selectorIndex := columnLength - MIN_OPTIONS
	selectedPerks := make(map[int]bool)
	perks := make([]int, numRolls)
	for i := 0; i < numRolls; i++ {
		var perk int
		for {
			perk = selectors[selectorIndex].Select()
			if !selectedPerks[perk] {
				selectedPerks[perk] = true
				break
			}
		}
		perks[i] = perk
	}

	return perks, nil
}

// Simulates a player rolling for weapon perks in a game, the system doesn't need to know about the perk
// themselves, just the indexes of the perks. The player may get a bonus roll(s) for columns three and four.
func main() {
	// Generate selectors for the weapon perks
	selectors, err := generateEquallyWeightedSelectors()
	if err != nil {
		panic(err)
	}

	player := Player{
		// Simulate player's may get a bonus roll(s) for columns three and four
		RollModifierForColumnThree: rand.IntN(3),
		RollModifierForColumnFour:  rand.IntN(3),
	}

	weapon := []Config{
		{Length: 9, NumberOfRolls: 2},
		{Length: 8, NumberOfRolls: 2},
		{Length: 11, NumberOfRolls: 1 + player.RollModifierForColumnThree},
		{Length: 12, NumberOfRolls: 1 + player.RollModifierForColumnFour},
	}

	weaponIndexes, err := getIndexes(selectors, weapon)
	if err != nil {
		panic(err)
	}

	fmt.Println("Weapon Perk Indexes:")
	for i, column := range weaponIndexes {
		fmt.Printf("Column %d: %v\n", i+1, column)
	}
}
