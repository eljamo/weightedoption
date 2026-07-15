//go:build amd64 || arm64 || mips64 || mips64le || ppc64 || ppc64le || riscv64 || s390x || wasm

package weightedoption

import (
	"math"
	"testing"
)

func TestNewSelector64Bit(t *testing.T) {
	t.Parallel()

	u64tests := []struct {
		name    string
		cs      []Option[rune, uint64]
		wantErr error
	}{
		{
			name:    "weight overflow from single option exceeding the system's math.MaxInt",
			cs:      []Option[rune, uint64]{{Data: 'a', Weight: uint64(math.MaxInt64) + 1}},
			wantErr: ErrSingleWeightOverflow,
		},
		{
			name:    "weight doesn't overflow if a single option is equal to the system's math.MaxInt",
			cs:      []Option[rune, uint64]{{Data: 'a', Weight: uint64(math.MaxInt64)}},
			wantErr: nil,
		},
		{
			name: "weight overflow from three options exceeding the system's math.MaxInt",
			cs: []Option[rune, uint64]{
				{Data: 'a', Weight: uint64(math.MaxInt64)/3 + 1},
				{Data: 'b', Weight: uint64(math.MaxInt64)/3 + 1},
				{Data: 'c', Weight: uint64(math.MaxInt64)/3 + 1},
			},
			wantErr: ErrTotalWeightOverflow,
		},
		{
			name: "weight doesn't overflow from three options close to but not exceeding the system's math.MaxInt",
			cs: []Option[rune, uint64]{
				{Data: 'a', Weight: uint64(math.MaxInt64) / 3},
				{Data: 'b', Weight: uint64(math.MaxInt64) / 3},
				{Data: 'c', Weight: uint64(math.MaxInt64)/3 + 1},
			},
			wantErr: nil,
		},
	}

	for _, tt := range u64tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewSelector(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewSelector() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewSelectorFloat64Large64Bit(t *testing.T) {
	t.Parallel()

	// 1e15 fits a 64-bit int; on 32-bit platforms it correctly overflows
	// (see the mirror case in weightedoption_32bit_test.go).
	if _, err := NewSelector(NewOption('a', 1e15)); err != nil {
		t.Errorf("NewSelector() error = %v, wantErr nil", err)
	}
}

func TestNewSelectorFloat64PrecisionCap64Bit(t *testing.T) {
	t.Parallel()

	// Regression: before the precision cap, 1/3 carried 16 fractional
	// digits, inflating the scale factor to 1e16 and overflowing on the
	// 1000.0 co-weight.
	s, err := NewSelector(NewOption('a', 1.0/3.0), NewOption('b', 1000.0))
	if err != nil {
		t.Fatal("Failed to create Selector:", err)
	}

	if s.totalWeight != 1_000_333_333_333 {
		t.Errorf("totalWeight = %d, want 1000333333333", s.totalWeight)
	}
}
