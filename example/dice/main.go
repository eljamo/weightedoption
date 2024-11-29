package main

import "github.com/eljamo/weightedoption/v3"

func getOptions(sides int) []weightedoption.Option[int, int] {
	options := make([]weightedoption.Option[int, int], sides)
	for i := 0; i < sides; i++ {
		options[i] = weightedoption.NewOption(i, 1)
	}

	return options
}

func main() {
	sides := 20

	// Create a new selector with 20 options, each with a weight of 1
	selector, err := weightedoption.NewSelector(
		getOptions(sides)...,
	)
	if err != nil {
		panic(err)
	}

	// Select a random option
	selected := selector.Select()

	// Print the selected option
	println(selected)
}
