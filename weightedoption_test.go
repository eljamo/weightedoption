package weightedoption

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

const (
	testOptions    int = 10
	testIterations int = 1_000_000
)

func TestNewSelector(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cs      []Option[rune, int]
		wantErr error
	}{
		{
			name:    "no options",
			cs:      []Option[rune, int]{},
			wantErr: ErrNoValidOptions,
		},
		{
			name:    "no options with weight greater than 0",
			cs:      []Option[rune, int]{{Data: 'a', Weight: 0}, {Data: 'b', Weight: 0}},
			wantErr: ErrNoValidOptions,
		},
		{
			name:    "one option with weight greater than 0",
			cs:      []Option[rune, int]{{Data: 'a', Weight: 1}},
			wantErr: nil,
		},
		{
			name:    "weight overflow",
			cs:      []Option[rune, int]{{Data: 'a', Weight: math.MaxInt/2 + 1}, {Data: 'b', Weight: math.MaxInt/2 + 1}},
			wantErr: ErrWeightOverflow,
		},
		{
			name:    "nominal case",
			cs:      []Option[rune, int]{{Data: 'a', Weight: 1}, {Data: 'b', Weight: 2}},
			wantErr: nil,
		},
		{
			name:    "one valid option and one invalid option with negative weight",
			cs:      []Option[rune, int]{{Data: 'a', Weight: 3}, {Data: 'b', Weight: -2}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSelector(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewSelector() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewSelectorUsingSortSearchInts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cs      []Option[rune, int]
		wantErr error
	}{
		{
			name:    "no options",
			cs:      []Option[rune, int]{},
			wantErr: ErrNoValidOptions,
		},
		{
			name:    "no options with weight greater than 0",
			cs:      []Option[rune, int]{{Data: 'a', Weight: 0}, {Data: 'b', Weight: 0}},
			wantErr: ErrNoValidOptions,
		},
		{
			name:    "one option with weight greater than 0",
			cs:      []Option[rune, int]{{Data: 'a', Weight: 1}},
			wantErr: nil,
		},
		{
			name:    "weight overflow",
			cs:      []Option[rune, int]{{Data: 'a', Weight: math.MaxInt/2 + 1}, {Data: 'b', Weight: math.MaxInt/2 + 1}},
			wantErr: ErrWeightOverflow,
		},
		{
			name:    "nominal case",
			cs:      []Option[rune, int]{{Data: 'a', Weight: 1}, {Data: 'b', Weight: 2}},
			wantErr: nil,
		},
		{
			name:    "one valid option and one invalid option with negative weight",
			cs:      []Option[rune, int]{{Data: 'a', Weight: 3}, {Data: 'b', Weight: -2}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSelectorUsingSortSearchInts(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewSelectorUsingSortSearchInts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSelector_Select(t *testing.T) {
	t.Parallel()

	options := mockFrequencyOptions(t, testOptions)
	picker, err := NewSelector(options...)
	if err != nil {
		t.Fatal("Failed to create Selector:", err)
	}

	counts := make(map[int]int)
	for i := 0; i < testIterations; i++ {
		c := picker.Select()
		counts[c]++
	}

	verifyFrequencyCounts(t, counts, options)
}

func mockFrequencyOptions(t *testing.T, n int) []Option[int, int] {
	t.Helper()
	options := make([]Option[int, int], 0, n)
	for i := 1; i <= n; i++ {
		c := NewOption(i, i)
		options = append(options, c)
	}
	t.Log("Mocked options:", options)
	return options
}

func verifyFrequencyCounts(t *testing.T, counts map[int]int, options []Option[int, int]) {
	t.Helper()

	for i := 0; i < len(options)-1; i++ {
		if counts[options[i].Data] > counts[options[i+1].Data] {
			t.Errorf(
				"Option with lower weight %d (count: %d) was selected more than option with higher weight %d (count: %d)",
				options[i].Weight, counts[options[i].Data], options[i+1].Weight, counts[options[i+1].Data],
			)
		}
	}
}

const BMMinOptions int = 10
const BMMaxOptions int = 10_000_000

func BenchmarkNewSelector(b *testing.B) {
	for n := BMMinOptions; n <= BMMaxOptions; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			options := mockOptions(n)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = NewSelector(options...)
			}
		})
	}
}

func BenchmarkSelect(b *testing.B) {
	for n := BMMinOptions; n <= BMMaxOptions; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			options := mockOptions(n)
			selector, err := NewSelector(options...)
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = selector.Select()
			}
		})
	}
}

func BenchmarkSelectParallel(b *testing.B) {
	for n := BMMinOptions; n <= BMMaxOptions; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			options := mockOptions(n)
			selector, err := NewSelector(options...)
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = selector.Select()
				}
			})
		})
	}
}

func BenchmarkNewSelectorUsingSortSearchInts(b *testing.B) {
	for n := BMMinOptions; n <= BMMaxOptions; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			options := mockOptions(n)
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = NewSelector(options...)
			}
		})
	}
}

func BenchmarkSortSearchIntsSelect(b *testing.B) {
	for n := BMMinOptions; n <= BMMaxOptions; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			options := mockOptions(n)
			selector, err := NewSelector(options...)
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_ = selector.Select()
			}
		})
	}
}

func BenchmarkSortSearchIntsSelectParallel(b *testing.B) {
	for n := BMMinOptions; n <= BMMaxOptions; n *= 10 {
		b.Run(fmt.Sprintf("size=%s", fmt1eN(n)), func(b *testing.B) {
			options := mockOptions(n)
			selector, err := NewSelector(options...)
			if err != nil {
				b.Fatal(err)
			}
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = selector.Select()
				}
			})
		})
	}
}

func mockOptions(n int) []Option[rune, int] {
	options := make([]Option[rune, int], 0, n)
	for i := 0; i < n; i++ {
		s := 'O'
		w := rand.Intn(10)
		c := NewOption(s, w)
		options = append(options, c)
	}
	return options
}

// fmt1eN returns simplified order of magnitude scientific notation for n,
// e.g. "1e2" for 100, "1e7" for 10 million.
func fmt1eN(n int) string {
	return fmt.Sprintf("1e%d", int(math.Log10(float64(n))))
}
