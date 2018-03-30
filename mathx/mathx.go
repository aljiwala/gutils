package mathx

import (
	"fmt"
	"math"
)

// Byte is a float64 where the String() method prints out a human-redable description.
type Byte float64

var units = []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB"}

// Round will round the value to the nearest natural number.
// .5 will be rounded up.
func Round(f float64) float64 {
	if f < 0 {
		return math.Ceil(f - 0.5)
	}
	return math.Floor(f + 0.5)
}

// RoundPlus will round the value to the given precision.
// e.g. RoundPlus(7.258, 2) will return 7.26
func RoundPlus(f float64, precision int) float64 {
	shift := math.Pow(10, float64(precision))
	return Round(f*shift) / shift
}

// Min gets the lowest of two numbers.
func Min(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

// Max gets the highest of two numbers.
func Max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

// MinMax should return min and max int value from given container.
func MinMax(container []int) (min int, max int) {
	if len(container) == 0 {
		return
	}

	min = container[0]
	max = container[0]
	for _, value := range container {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}

	return
}

// NextMin should return min int value from given container.
func NextMin(v []int, cur int) (int, bool) {
	minElems := []int{}
	for _, e := range v {
		if e < cur {
			minElems = append(minElems, e)
		}
	}

	if len(minElems) == 0 {
		return cur, true
	}

	_, max := MinMax(minElems)
	return max, false
}

// NextMax should return max int value from given container.
func NextMax(v []int, cur int) (int, bool) {
	maxElems := []int{}
	for _, e := range v {
		if e > cur {
			maxElems = append(maxElems, e)
		}
	}

	if len(maxElems) == 0 {
		return cur, true
	}

	min, _ := MinMax(maxElems)
	return min, false
}

// Limit a value between a lower and upper limit.
func Limit(v, lower, upper float64) float64 {
	return math.Max(math.Min(v, upper), lower)
}

// DivideCeil divides two integers and rounds up, rather than down (which is
// what happens when you do int64/int64).
func DivideCeil(count int64, pageSize int64) int64 {
	return int64(math.Ceil(float64(count) / float64(pageSize)))
}

// IsSignedZero checks if this number is a signed zero (i.e. -0, instead of +0).
func IsSignedZero(f float64) bool {
	return math.Float64bits(f)^uint64(1<<63) == 0
}

// String is the string representation of byte.
func (b Byte) String() string {
	i := 0
	for ; i < len(units); i++ {
		if b < 1024 {
			return fmt.Sprintf("%.1f%s", b, units[i])
		}
		b /= 1024
	}
	return fmt.Sprintf("%.1f%s", b*1024, units[i-1])
}
