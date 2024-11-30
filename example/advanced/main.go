package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/eljamo/weightedoption/v3"
)

type Player struct {
	RollBonusSocketThree int
	RollBonusSocketFour  int
}

type ItemSocketConfig struct {
	SocketPlugCount int
	Rolls           int
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

func getIndexes(selectors []*weightedoption.Selector[int, int], selections []ItemSocketConfig) ([][]int, error) {
	plugIndexes := make([][]int, len(selections))
	for i, config := range selections {
		if config.SocketPlugCount < MIN_OPTIONS {
			return nil, fmt.Errorf("length value (%d) must be equal or greater than to %d", config.SocketPlugCount, MIN_OPTIONS)
		}

		if config.SocketPlugCount > MAX_OPTIONS {
			return nil, fmt.Errorf("length value (%d) must be equal or less than %d", config.SocketPlugCount, MAX_OPTIONS)
		}

		plugs, err := selectIndexes(selectors, config.SocketPlugCount, config.Rolls)
		if err != nil {
			return nil, err
		}
		plugIndexes[i] = plugs
	}

	return plugIndexes, nil
}

func selectIndexes(selectors []*weightedoption.Selector[int, int], columnLength, numRolls int) ([]int, error) {
	if numRolls < 1 {
		return nil, fmt.Errorf("numRolls value (%d) must be equal or greater than %d", numRolls, MIN_ROLLS)
	}

	selectorIndex := columnLength - MIN_OPTIONS
	selectedPerks := make(map[int]bool)
	plugs := make([]int, numRolls)
	for i := 0; i < numRolls; i++ {
		var plug int
		for {
			plug = selectors[selectorIndex].Select()
			if !selectedPerks[plug] {
				selectedPerks[plug] = true
				break
			}
		}
		plugs[i] = plug
	}

	return plugs, nil
}

// This is using Destiny 2's weapon perk system as an example. Items have sockets which have plugs.
// Destiny 2's uses this system for it's weapons perks. Each weapon has a certain number of sockets
// with each socket having at least one plug (perk) to choose from a pool of plugs. The number of
// plugs in each socket can vary. The player may get bonus rolls for certain sockets.
func main() {
	// Generate selectors for the weapon perks
	selectors, err := generateEquallyWeightedSelectors()
	if err != nil {
		panic(err)
	}

	player := Player{
		// Simulate that players may get bonus rolls for sockets three and four
		RollBonusSocketThree: rand.IntN(3),
		RollBonusSocketFour:  rand.IntN(3),
	}

	weapon := []ItemSocketConfig{
		{SocketPlugCount: 9, Rolls: 2},
		{SocketPlugCount: 8, Rolls: 2},
		{SocketPlugCount: 11, Rolls: 1 + player.RollBonusSocketThree},
		{SocketPlugCount: 12, Rolls: 1 + player.RollBonusSocketFour},
	}

	weaponSocketPlugIndexes, err := getIndexes(selectors, weapon)
	if err != nil {
		panic(err)
	}

	fmt.Println("Plug Indexes for Item Sockets:")
	for i, column := range weaponSocketPlugIndexes {
		fmt.Printf("Socket %d: %v\n", i+1, column)
	}
}
