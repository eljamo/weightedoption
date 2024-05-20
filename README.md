# weightedoption

A Go package for weighted random option selection

## Example Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/eljamo/weightedoption/v2"
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
```

```bash
go run ./example/example.go
❌❌❌❌❌❌❌❌❌❌❌❌🔫❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌🔫❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌❌🔫❌❌❌❌❌❌❌❌❌❌❌❌

🔫: 3	❌ 97
```

## Contributing

If you'd like to contribute, please fork the repository and work your magic. Open a pull request to the `main` branch if it is a `bugfix` or `feature` branch. If a `hotfix`, open a pull request to the respective `release` branch.

### Run the tests

```bash
go test --race --shuffle on ./...
```

### Run the example

```bash
go run ./example/example.go
```

### Branching Strategy

- `bugfix/*`
- `feature/*`
- `main`
- `release/v*`
- `hotfix/*`

### Workflow Breakdown

- **`feature` and `bugfix` branches:**

  - `feature` Branches: Create these for new features you are developing.
  - `bugfix` Branches: Create these for fixing bugs identified in the `main` branch.
  - Once the feature or fix is complete, merge these branches back into the `main` branch.

- **`main` branch:**

  - The `main` branch serves as the central integration branch where all `feature` and `bugfix` branches are merged.

- **`release` branch:**

  - When you are ready to make a release, create a `release` branch from the `main` branch.
  - Perform any final testing and adjustments on the `release` branch.
  - Once the release is stable, it can be deployed from this branch.

- **`hotfix` branch:**

  - If a critical issue is found after the release, create a `hotfix` branch from the `release` branch.
  - Fix the issue on the `hotfix` branch and then merge it back into both the `release` and `main` branches if applicable.
  - This ensures that the fix is included in the current release(s) and the `main` branch.

### Example Scenarios

- **Developing a feature:**

  1. Create a `feature` branch from `main`.
  2. Develop and test the feature on the `feature` branch.
  3. Merge the `feature` branch into `main`.

- **Fixing a bug:**

  1. Create a `bugfix` branch from `main`.
  2. Fix and test the bug on the `bugfix` branch.
  3. Merge the `bugfix` branch into `main`.

- **Making a release:**

  1. Create a `release` branch from `main` or another `release` branch.
  2. Perform testing on `release` branch.
  3. Deploy the `release` branch.

- **Applying a hotfix:**
  1. Create a `hotfix` branch from a `release` branch.
  2. Fix the critical issue on `hotfix` branch.
  3. Merge `hotfix` branch into the `release` and `main` branches.
