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
			name:    "weight overflow from single uint64 exceeding system math.MaxInt",
			cs:      []Option[rune, uint64]{{Data: 'a', Weight: uint64(math.MaxInt) + 1}},
			wantErr: ErrWeightOverflow,
		},
	}

	for _, tt := range u64tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSelector(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewSelector() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
