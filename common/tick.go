package common

import (
	"encoding/json"
	"math/rand"
)

// Represent a transaction
type TickValue struct {
	TickName string  `json:"TickName"`
	Value    float64 `json:"Value"`
}

var deviations []float64

func init() {
	mapped := map[int]float64{
		1:     10,
		2:     8,
		3:     6,
		5:     4,
		80:    2,
		200:   1,
		1000:  0.5,
		10000: 0.2,
	}

	count := 0
	for k := range mapped {
		count = count + k
	}

	count = count * 2
	deviations = make([]float64, count)

	for k, v := range mapped {
		for i := 0; i < k; i++ {
			deviations[i] = -v
			deviations[count-i-1] = v
		}
	}
}

func (tv TickValue) SerializeJson() ([]byte, error) {
	return json.Marshal(tv)
}

func (tv TickValue) GetNextTransaction() (TickValue, error) {
	// Value change from -10% to +10% randomly
	return tv.getRandomValue(), nil
}

func (tv TickValue) getRandomValue() TickValue {
	rndDeviation := deviations[rand.Intn(len(deviations))-1]
	rndValue := tv.Value + (tv.Value * rndDeviation / 100)
	return TickValue{
		TickName: tv.TickName,
		Value:    rndValue,
	}
}
