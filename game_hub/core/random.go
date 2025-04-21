package core

import (
	"math"
	"math/rand"
	"time"
)

type RandomGenerator struct {
	rand           *rand.Rand
	minRangeNumber int
	maxRangeNumber int
}

func NewRandomGenerator() *RandomGenerator {
	return &RandomGenerator{
		rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
		minRangeNumber: 0,
		maxRangeNumber: math.MaxInt32,
	}
}

func (r *RandomGenerator) Generate(from, to int) (int, error) {
	if from > to || from < r.minRangeNumber || to > r.maxRangeNumber {
		return 0, NewAppError(ErrInvalidRange, "invalid_range", map[string]any{
			"min": from,
			"max": to,
		})
	}
	return r.rand.Intn(to-from+1) + from, nil
}
