# weightedoption

A Go package for weighted random option selection

## Example Usage

```go
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
	for i := range chances {
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
```

```bash
go run ./example/exotic/main.go
❌❌❌❌❌❌❌❌❌❌❌❌🔫❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌🔫❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌🔫❌❌❌❌❌❌❌❌❌❌❌❌

🔫: 3	❌ 97
```

### Float64 Weights

```go
// Weights may also be float64; they are scaled to integers internally.
// Precision is capped at 9 fractional digits; weights are rounded to that
// precision, and a weight that rounds to zero returns an error.
s, err := weightedoption.NewSelector(
    weightedoption.NewOption("rare", 0.5),
    weightedoption.NewOption("common", 99.5),
)
```

### Deterministic Selection

`SelectFrom` selects using your own `math/rand/v2` generator, which is useful when you need deterministic or reproducible behavior:

```go
rng := rand.New(rand.NewPCG(1, 2))
data := s.SelectFrom(rng)
```

## Contributing

If you'd like to contribute, please fork the repository and work your magic. Open a pull request to the `main` branch if it is a `bugfix` or `feature` branch. If it is a `hotfix` branch, open a pull request to the respective `release` branch.

### Run the tests

```bash
go test --race --shuffle=on ./...
```

### Run the examples

```bash
go run ./example/advanced/main.go
go run ./example/dice/main.go
go run ./example/exotic/main.go
go run ./example/exotic_pity/main.go
go run ./example/gacha/main.go
go run ./example/gacha_pity/main.go
```
