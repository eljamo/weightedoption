package weightedoption

import (
	"errors"
	"math"
	"math/rand/v2"
	"sort"
)

var (
	ErrWeightOverflow = errors.New("sum of Option weights exceeds total")
	ErrNoValidOptions = errors.New("0 Option(s) with Weight >= 1")
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

// NewOption creates a new Option.
func NewOption[DataType any, WeightIntegerType WeightIntegerConstraint](
	data DataType,
	weight WeightIntegerType,
) Option[DataType, WeightIntegerType] {
	return Option[DataType, WeightIntegerType]{Data: data, Weight: weight}
}

// SearchIntsFuncSignature is the signature of the function used to search for an integer in a sorted slice of integers.
type SearchIntsFuncSignature func(runningTotalWeights []int, randInt int) int

// Selector is a struct that holds a slice of Options and their cumulative weights.
type Selector[DataType any, WeightIntegerType WeightIntegerConstraint] struct {
	options             []Option[DataType, WeightIntegerType]
	runningTotalWeights []int
	totalWeight         int
	searchIntsFunc      SearchIntsFuncSignature
}

// NewSelector creates a new Selector.
func NewSelector[DataType any, WeightIntegerType WeightIntegerConstraint](
	options ...Option[DataType, WeightIntegerType],
) (*Selector[DataType, WeightIntegerType], error) {
	var filteredOptions []Option[DataType, WeightIntegerType]
	for _, opt := range options {
		if opt.Weight > 0 {
			filteredOptions = append(filteredOptions, opt)
		}
	}

	sort.Slice(filteredOptions, func(i, j int) bool {
		return filteredOptions[i].Weight < filteredOptions[j].Weight
	})

	runningTotalWeights := make([]int, len(filteredOptions))
	totalWeight := 0

	for i, opt := range filteredOptions {
		if uint(opt.Weight) >= math.MaxInt {
			return nil, ErrWeightOverflow
		}

		weight := int(opt.Weight)
		if weight > math.MaxInt-totalWeight {
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
		searchIntsFunc:      searchInts,
	}, nil
}

// NewSelectorWithCustomSearchIntsFunc creates a new Selector with a custom searchIntsFunc.
func NewSelectorWithCustomSearchIntsFunc[DataType any, WeightIntegerType WeightIntegerConstraint](
	searchIntsFunc SearchIntsFuncSignature,
	options ...Option[DataType, WeightIntegerType],
) (*Selector[DataType, WeightIntegerType], error) {
	selector, err := NewSelector(options...)
	if err != nil {
		return nil, err
	}

	selector.searchIntsFunc = searchIntsFunc
	return selector, nil
}

// NewSelectorUsingSortSearchInts creates a new Selector using the sort.SearchInts function.
func NewSelectorUsingSortSearchInts[DataType any, WeightIntegerType WeightIntegerConstraint](
	options ...Option[DataType, WeightIntegerType],
) (*Selector[DataType, WeightIntegerType], error) {
	return NewSelectorWithCustomSearchIntsFunc(sort.SearchInts, options...)
}

// Select returns a single option from the Selector.
func (s Selector[DataType, WeightIntegerType]) Select() DataType {
	r := rand.IntN(s.totalWeight) + 1
	i := s.searchIntsFunc(s.runningTotalWeights, r)
	return s.options[i].Data
}

// searchInts searches for the index of the first element in runningTotalWeights
// that is greater than or equal to randInt. The slice must be sorted in
// ascending order.
func searchInts(runningTotalWeights []int, randInt int) int {
	start, end := 0, len(runningTotalWeights)
	for start < end {
		mid := int(uint(start+end) >> 1)
		if runningTotalWeights[mid] < randInt {
			start = mid + 1
		} else {
			end = mid
		}
	}

	return start
}
