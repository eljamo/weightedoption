package main

import (
	"fmt"

	"github.com/eljamo/weightedoption/v3"
)

func getOptions(sides int) []weightedoption.Option[int, int] {
	options := make([]weightedoption.Option[int, int], sides)
	for i := 0; i < sides; i++ {
		options[i] = weightedoption.NewOption(i+1, 1)
	}

	return options
}

func rollDice(sides int) (int, error) {
	selector, err := weightedoption.NewSelector(
		getOptions(sides)...,
	)
	if err != nil {
		return 0, err
	}

	return selector.Select(), nil
}

func main() {
	diceSides := []int{4, 6, 8, 10, 12, 20}

	for _, sides := range diceSides {
		result, err := rollDice(sides)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Rolled a d%d: %d\n", sides, result)
	}
}
