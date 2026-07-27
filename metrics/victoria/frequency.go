package victoria

import (
	"github.com/koykov/pbtk/frequency"
	"github.com/koykov/vmchain"
)

type mwFrequency struct {
	name string
}

func NewFrequency(name string) frequency.MetricsWriter {
	return &mwFrequency{name: name}
}

func (mw *mwFrequency) Add(err error) error {
	result := "success"
	if err != nil {
		result = "fail"
	}
	vmchain.Counter("frequency_add").WithLabel("name", mw.name).
		WithLabel("result", result).Inc()
	return err
}

func (mw *mwFrequency) Estimate(value uint64) uint64 {
	vmchain.Histogram("frequency_estimation").WithLabel("name", mw.name).Update(float64(value))
	return value
}

var _ = NewFrequency
