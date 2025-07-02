package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/eljamo/weightedoption/v3"
)

type PlayerDataForActivity struct {
	NumOfCompletions  int
	NumOfAchievements int
}

func randomBool() bool {
	return rand.IntN(2) == 0
}

func generateOptionsFromPlayerData(playerData *PlayerDataForActivity) []weightedoption.Option[bool, int] {
	const baseChance = 5
	const maxChance = 100
	finalChance := baseChance + playerData.NumOfCompletions + playerData.NumOfAchievements
	noDropChance := 0

	if finalChance > maxChance {
		finalChance = maxChance
	} else {
		noDropChance = maxChance - finalChance
	}

	return []weightedoption.Option[bool, int]{
		weightedoption.NewOption(true, finalChance),
		weightedoption.NewOption(false, noDropChance),
	}
}

func (p *PlayerDataForActivity) Select() (bool, error) {
	options := generateOptionsFromPlayerData(p)
	s, err := weightedoption.NewSelector(options...)
	if err != nil {
		return false, err
	}

	return s.Select(), nil
}

func main() {
	playerData := &PlayerDataForActivity{
		NumOfCompletions:  0,
		NumOfAchievements: 0,
	}

	// Run the simulation until the player gets the exotic weapon
	dropped := false

	for !dropped {
		playerData.NumOfCompletions++

		if randomBool() {
			playerData.NumOfAchievements++
		}

		result, err := playerData.Select()
		if err != nil {
			panic(err)
		}

		if result {
			dropped = true
		}

		if dropped {
			fmt.Printf("Exotic weapon dropped after %d completions and %d achievements\n", playerData.NumOfCompletions, playerData.NumOfAchievements)
			break
		}
	}
}
