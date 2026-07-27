package victoria

import (
	"github.com/koykov/pbtk/frequency"
	"github.com/koykov/vmchain"
)

type mwSignedFrequency struct {
	mwFrequency
}

func NewSignedFrequency(name string) frequency.SignedMetricsWriter {
	return &mwSignedFrequency{mwFrequency{name: name}}
}

func (mw *mwSignedFrequency) Estimate(value int64) int64 {
	vmchain.Histogram("frequency_estimation").WithLabel("name", mw.name).Update(float64(value))
	return value
}

var _ = NewSignedFrequency
