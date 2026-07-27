package victoria

import (
	"github.com/koykov/pbtk/frequency"
	"github.com/koykov/vmchain"
)

type mwPreciseFrequency struct {
	mwFrequency
}

func NewPreciseFrequency(name string) frequency.PreciseMetricsWriter {
	return &mwPreciseFrequency{mwFrequency{name: name}}
}

func (mw *mwPreciseFrequency) Estimate(value float64) float64 {
	vmchain.Histogram("frequency_estimation").WithLabel("name", mw.name).Update(value)
	return value
}

var _ = NewPreciseFrequency
