package weightedoption

import (
	"cmp"
	"errors"
	"math"
	"math/rand/v2"
	"slices"
)

var (
	ErrNoValidOptions = errors.New("no Options found with Weight >= 1")
	ErrWeightOverflow = errors.New("Option weight exceeds max integer value for this system's architecture")
)

// WeightIntegerConstraint is a type constraint for the Weight field of the Option struct.
type WeightIntegerConstraint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Option is a struct that holds a data value and its associated weight.
type Option[DataType any, WeightIntegerType WeightIntegerConstraint] struct {
	Data   DataType
	Weight WeightIntegerType
}

// NewOption creates a new Option with the provided data and weight.
func NewOption[DataType any, WeightIntegerType WeightIntegerConstraint](
	data DataType,
	weight WeightIntegerType,
) Option[DataType, WeightIntegerType] {
	return Option[DataType, WeightIntegerType]{Data: data, Weight: weight}
}

// Selector is a struct that holds a slice of Options, their running total weights, and the total weight.
type Selector[DataType any, WeightIntegerType WeightIntegerConstraint] struct {
	options             []Option[DataType, WeightIntegerType]
	runningTotalWeights []uint
	totalWeight         uint
}

// NewSelector creates a new Selector for selecting provided Options.
func NewSelector[DataType any, WeightIntegerType WeightIntegerConstraint](
	options ...Option[DataType, WeightIntegerType],
) (*Selector[DataType, WeightIntegerType], error) {
	var filteredOptions []Option[DataType, WeightIntegerType]
	for _, opt := range options {
		if opt.Weight > 0 {
			filteredOptions = append(filteredOptions, opt)
		}
	}

	slices.SortFunc(filteredOptions, func(a, b Option[DataType, WeightIntegerType]) int {
		return cmp.Compare(a.Weight, b.Weight)
	})

	var totalWeight uint
	runningTotalWeights := make([]uint, len(filteredOptions))

	for i, opt := range filteredOptions {
		weight := uint(opt.Weight)
		if weight >= math.MaxInt {
			return nil, ErrWeightOverflow
		}

		if (math.MaxInt - totalWeight) <= weight {
			return nil, ErrWeightOverflow
		}

		totalWeight += weight
		runningTotalWeights[i] = totalWeight
	}

	if totalWeight < 1 {
		return nil, ErrNoValidOptions
	}

	return &Selector[DataType, WeightIntegerType]{
		options:             filteredOptions,
		runningTotalWeights: runningTotalWeights,
		totalWeight:         totalWeight,
	}, nil
}

// Select returns a single Option.Data from Selector.options
func (s Selector[DataType, WeightIntegerType]) Select() DataType {
	r := rand.UintN(s.totalWeight) + 1
	i, _ := slices.BinarySearch(s.runningTotalWeights, r)
	return s.options[i].Data
}
