package weightedoption

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
)

var (
	// ErrNoValidOptions is used and returned when no Options are found with a Weight >= 1.
	ErrNoValidOptions = errors.New("no Options found with Weight >= 1")
	// Error for individual weight exceeding system max integer value.
	ErrSingleWeightOverflow = errors.New("Option weight exceeds max integer value for this system's architecture")

	// Error for total weight exceeding system max integer value when summing.
	ErrTotalWeightOverflow = errors.New("total weight exceeds max integer value for this system's architecture")
)

// WeightConstraint is a type constraint for the Weight field of the Option struct.
type WeightConstraint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float64
}

// Option is a struct that holds a data value and its associated weight.
// If using floating point weights, the weights will be multiplied by 100 to convert them to integers.
// So 0.5 will become 50, 0.75 will become 75, etc. It's only precise to 2 decimal places. If you need
// more precision you should convert the weights to integers yourself.
type Option[DataType any, WeightType WeightConstraint] struct {
	Data   DataType
	Weight WeightType
}

// NewOption creates a new Option with the provided data and weight.
// If using floating point weights, the weights will be multiplied by 100 to convert them to integers.
// So 0.5 will become 50, 0.75 will become 75, etc. It's only precise to 2 decimal places. If you need
// more precision you should convert the weights to integers yourself.
func NewOption[DataType any, WeightType WeightConstraint](
	data DataType,
	weight WeightType,
) Option[DataType, WeightType] {
	return Option[DataType, WeightType]{Data: data, Weight: weight}
}

// Selector is a struct that holds a slice of Options, their running total weights, and the total weight.
type Selector[DataType any, WeightType WeightConstraint] struct {
	options              []DataType
	cumulativeWeightSums []uint
	totalWeight          uint
}

func isFloat64[WeightType WeightConstraint](val WeightType) bool {
	_, is := any(val).(float64)
	return is
}

// countFractionalDigits counts the number of digits after the decimal point, returning an error for invalid floats.
func countFractionalDigits(f float64) (int, error) {
	// Handle invalid cases early
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("invalid float: %v", f)
	}

	// If the number is already an integer, no fractional digits
	if f == math.Trunc(f) {
		return 0, nil
	}

	// Count the number of fractional digits by multiplying until the fractional part becomes zero
	count := 0
	for math.Trunc(f) != f {
		f *= 10
		count++
	}
	return count, nil
}

func processWeight(weight float64) (int, error) {
	digits, err := countFractionalDigits(weight)
	if err != nil {
		return 0, fmt.Errorf("invalid weight: weight=%v", weight)
	}
	return digits, nil
}

// maxFractionalDigits finds the maximum number of fractional digits in a list of floats and returns an error if any float is invalid.
func maxFractionalDigits[DataType any, WeightType WeightConstraint](options []Option[DataType, WeightType]) (int, error) {
	maxDigits := 0

	for _, opt := range options {
		var digits int

		switch weight := any(opt.Weight).(type) {
		case float64:
			_, err := processWeight(weight)
			if err != nil {
				return 0, fmt.Errorf("invalid float64 weight for option: option=%v", opt.Data)
			}
		default:
			return 0, fmt.Errorf("weight failed type assertion, not a valid float64: option=%v, weight=%v", opt.Data, opt.Weight)
		}

		if digits > maxDigits {
			maxDigits = digits
		}
	}

	return maxDigits, nil
}

const decimalBase = 10

func scaleFloat64ToInt[DataType any, WeightType WeightConstraint](maxPrecision int, options []Option[DataType, WeightType]) ([]Option[DataType, WeightType], error) {
	scaleFactor := math.Pow(decimalBase, float64(maxPrecision))
	for i, opt := range options {
		switch weight := any(opt.Weight).(type) {
		case float64:
			scaledWeight := int(math.Round(weight * scaleFactor))
			options[i].Weight = WeightType(scaledWeight)
		default:
			return nil, fmt.Errorf("weight failed type assertion, not a valid float64: option=%v, weight=%v", opt.Data, opt.Weight)
		}
	}
	return options, nil
}

func prepareOptions[DataType any, WeightType WeightConstraint](
	options ...Option[DataType, WeightType],
) ([]Option[DataType, WeightType], error) {
	var filteredOptions []Option[DataType, WeightType]

	// Filter out options with non-positive weights
	for _, opt := range options {
		if opt.Weight > 0 {
			filteredOptions = append(filteredOptions, opt)
		}
	}

	if len(filteredOptions) == 0 {
		return nil, ErrNoValidOptions
	}

	// Sort options by weight in ascending order
	slices.SortFunc(options, func(a, b Option[DataType, WeightType]) int {
		return cmp.Compare(a.Weight, b.Weight)
	})

	// If integers just return the list
	if !isFloat64(filteredOptions[0].Weight) {
		return filteredOptions, nil
	}

	// Find the maximum number of fractional digits, returning an error if any float is invalid.
	maxDigits, err := maxFractionalDigits(filteredOptions)
	if err != nil {
		return nil, err
	}

	return scaleFloat64ToInt(maxDigits, filteredOptions)
}

// NewSelector creates a new Selector for selecting provided Options.
func NewSelector[DataType any, WeightType WeightConstraint](
	opts ...Option[DataType, WeightType],
) (*Selector[DataType, WeightType], error) {
	opts, err := prepareOptions(opts...)
	if err != nil {
		return nil, err
	}

	var totalWeight uint
	cumulativeWeightSums := make([]uint, len(opts))
	options := make([]DataType, len(opts))
	for i, opt := range opts {
		weight := uint(opt.Weight)
		// Check for overflow
		if weight > math.MaxInt {
			return nil, ErrSingleWeightOverflow
		}

		if (math.MaxInt - totalWeight) < weight {
			return nil, ErrTotalWeightOverflow
		}

		totalWeight += weight
		options[i] = opt.Data
		cumulativeWeightSums[i] = totalWeight
	}

	if totalWeight < 1 {
		return nil, ErrNoValidOptions
	}

	return &Selector[DataType, WeightType]{
		options:              options,
		cumulativeWeightSums: cumulativeWeightSums,
		totalWeight:          totalWeight,
	}, nil
}

// Select returns a single DataType from Selector.Options
func (s Selector[DataType, WeightType]) Select() DataType {
	r := rand.UintN(s.totalWeight) + 1
	i, _ := slices.BinarySearch(s.cumulativeWeightSums, r)
	return s.options[i]
}
