//go:build amd64 || arm64 || mips64 || mips64le || ppc64 || ppc64le || riscv64 || s390x || wasm
// +build amd64 arm64 mips64 mips64le ppc64 ppc64le riscv64 s390x wasm

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
