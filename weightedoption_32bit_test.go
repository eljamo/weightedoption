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
			name:    "weight overflow from single option exceeding system's math.MaxInt",
			cs:      []Option[rune, uint32]{{Data: 'a', Weight: uint32(math.MaxInt32) + 1}},
			wantErr: ErrSingleWeightOverflow,
		},
		{
			name:    "weight doesn't overflow if a single option is equal to the system's math.MaxInt",
			cs:      []Option[rune, uint32]{{Data: 'a', Weight: uint32(math.MaxInt32)}},
			wantErr: nil,
		},
		{
			name: "weight overflow from three options exceeding the system's math.MaxInt",
			cs: []Option[rune, uint32]{
				{Data: 'a', Weight: uint32(math.MaxInt32)/3 + 1},
				{Data: 'b', Weight: uint32(math.MaxInt32)/3 + 1},
				{Data: 'c', Weight: uint32(math.MaxInt32)/3 + 1},
			},
			wantErr: ErrTotalWeightOverflow,
		},
		{
			name: "weight doesn't overflow from three options close to but not exceeding the system's math.MaxInt",
			cs: []Option[rune, uint32]{
				{Data: 'a', Weight: uint32(math.MaxInt32) / 3},
				{Data: 'b', Weight: uint32(math.MaxInt32) / 3},
				{Data: 'c', Weight: uint32(math.MaxInt32)/3 + 1},
			},
			wantErr: nil,
		},
	}

	for _, tt := range u32tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewSelector(tt.cs...)
			if err != tt.wantErr {
				t.Errorf("NewSelector() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
