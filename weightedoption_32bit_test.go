//go:build 386 || arm || mips || mipsle
// +build 386 arm mips mipsle

package weightedoption

import (
	"math"
	"testing"
)

func TestNewSelector32Bit(t *testing.T) {
	t.Parallel()

	u32tests := []struct {
		name    string
		cs      []Option[rune, uint32]
		wantErr error
	}{
		{
			name:    "weight overflow from single uint32 exceeding system math.MaxInt",
			cs:      []Option[rune, uint32]{{Data: 'a', Weight: uint32(math.MaxInt) + 1}},
			wantErr: ErrWeightOverflow,
		},
	}

	for _, tt := range u32tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSelector(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewSelector() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
