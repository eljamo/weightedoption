package weightedoption

import (
	"fmt"
	"math"
	"math/rand"
	randv2 "math/rand/v2"
	"slices"
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
			wantErr: ErrTotalWeightOverflow,
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
			t.Parallel()
			_, err := NewSelector(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewSelector() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewSelectorDoesNotMutateInput(t *testing.T) {
	t.Parallel()

	opts := []Option[rune, int]{
		{Data: 'c', Weight: 3},
		{Data: 'a', Weight: 1},
		{Data: 'b', Weight: 2},
	}
	want := []Option[rune, int]{
		{Data: 'c', Weight: 3},
		{Data: 'a', Weight: 1},
		{Data: 'b', Weight: 2},
	}

	if _, err := NewSelector(opts...); err != nil {
		t.Fatalf("NewSelector() unexpected error: %v", err)
	}

	if !slices.Equal(opts, want) {
		t.Errorf("NewSelector() mutated its input: got %v, want %v", opts, want)
	}
}

func TestCountFractionalDigits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   float64
		want    int
		wantErr bool
	}{
		{name: "one fractional digit", input: 0.5, want: 1},
		{name: "two fractional digits", input: 0.25, want: 2},
		{name: "not exactly representable 0.1", input: 0.1, want: 1},
		{name: "not exactly representable 0.3", input: 0.3, want: 1},
		{name: "not exactly representable 0.07", input: 0.07, want: 2},
		{name: "integral float", input: 1.0, want: 0},
		{name: "large integral float", input: 100.0, want: 0},
		{name: "three fractional digits", input: 1.125, want: 3},
		{name: "sixteen digits capped to nine", input: 1.0 / 3.0, want: 9},
		{name: "ten digits capped to nine", input: 0.1234567891, want: 9},
		{name: "NaN", input: math.NaN(), wantErr: true},
		{name: "positive infinity", input: math.Inf(1), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := countFractionalDigits(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("countFractionalDigits(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("countFractionalDigits(%v) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewSelectorFloat64Scaling(t *testing.T) {
	t.Parallel()

	s, err := NewSelector(NewOption('a', 0.5), NewOption('b', 0.25))
	if err != nil {
		t.Fatal("Failed to create Selector:", err)
	}

	if s.totalWeight != 75 {
		t.Errorf("totalWeight = %d, want 75", s.totalWeight)
	}
	if len(s.cumulativeWeightSums) != 2 {
		t.Fatalf("len(cumulativeWeightSums) = %d, want 2", len(s.cumulativeWeightSums))
	}
	if got := s.cumulativeWeightSums[len(s.cumulativeWeightSums)-1]; got != 75 {
		t.Errorf("cumulativeWeightSums[last] = %d, want 75", got)
	}
}

func TestSelectorSelectFloat64Distribution(t *testing.T) {
	t.Parallel()

	s, err := NewSelector(NewOption('a', 0.5), NewOption('b', 0.25))
	if err != nil {
		t.Fatal("Failed to create Selector:", err)
	}

	counts := make(map[rune]int)
	for range 100_000 {
		counts[s.Select()]++
	}

	if counts['a'] == 0 {
		t.Error("option 'a' (weight 0.5) was never selected")
	}
	if counts['b'] == 0 {
		t.Error("option 'b' (weight 0.25) was never selected")
	}
	if counts['a'] <= counts['b'] {
		t.Errorf("option 'a' (weight 0.5, count %d) should be selected more than option 'b' (weight 0.25, count %d)", counts['a'], counts['b'])
	}
}

func TestNewSelectorFloat64PrecisionCap(t *testing.T) {
	t.Parallel()

	s, err := NewSelector(NewOption('a', 1.0/3.0), NewOption('b', 1.0))
	if err != nil {
		t.Fatal("Failed to create Selector:", err)
	}

	// 1/3 capped to 9 digits scales to 333_333_333; 1.0 scales to 1_000_000_000.
	if s.totalWeight != 1_333_333_333 {
		t.Errorf("totalWeight = %d, want 1333333333", s.totalWeight)
	}
}

func TestNewSelectorFloat64PrecisionLoss(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cs      []Option[rune, float64]
		wantErr error
	}{
		{
			name:    "weight rounding to zero at capped precision errors",
			cs:      []Option[rune, float64]{{Data: 'a', Weight: 1e-10}, {Data: 'b', Weight: 0.5}},
			wantErr: ErrWeightPrecisionLoss,
		},
		{
			name:    "single tiny weight errors",
			cs:      []Option[rune, float64]{{Data: 'a', Weight: 1e-10}},
			wantErr: ErrWeightPrecisionLoss,
		},
		{
			name:    "smallest representable capped weight is accepted",
			cs:      []Option[rune, float64]{{Data: 'a', Weight: 1e-9}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewSelector(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewSelector() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewSelectorFloat64Overflow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cs      []Option[rune, float64]
		wantErr error
	}{
		{
			name:    "huge integral float weight overflows",
			cs:      []Option[rune, float64]{{Data: 'a', Weight: 1e300}},
			wantErr: ErrSingleWeightOverflow,
		},
		{
			name:    "max float64 weight overflows",
			cs:      []Option[rune, float64]{{Data: 'a', Weight: math.MaxFloat64}},
			wantErr: ErrSingleWeightOverflow,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewSelector(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewSelector() error = %v, wantErr %v", err, tt.wantErr)
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
	for range testIterations {
		c := picker.Select()
		counts[c]++
	}

	verifyFrequencyCounts(t, counts, options)
}

func TestSelector_SelectFromUsesProvidedRand(t *testing.T) {
	t.Parallel()

	selector, err := NewSelector(
		NewOption("a", 1),
		NewOption("b", 2),
		NewOption("c", 3),
	)
	if err != nil {
		t.Fatal("Failed to create Selector:", err)
	}

	rng1 := randv2.New(randv2.NewPCG(1, 2))
	rng2 := randv2.New(randv2.NewPCG(1, 2))

	for range 20 {
		got := selector.SelectFrom(rng1)
		want := selector.SelectFrom(rng2)
		if got != want {
			t.Fatalf("SelectFrom() with identical generators got %q, want %q", got, want)
		}
	}
}

func mockFrequencyOptions(t *testing.T, n int) []Option[int, int] {
	t.Helper()
	options := make([]Option[int, int], 0, n)
	for i := 1; i <= n; i++ {
		c := NewOption(i, i)
		options = append(options, c)
	}
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

const (
	BMMinOptions int = 10
	BMMaxOptions int = 10_000_000
)

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

func mockOptions(n int) []Option[rune, int] {
	options := make([]Option[rune, int], 0, n)
	for range n {
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
