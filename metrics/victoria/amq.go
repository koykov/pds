package victoria

import (
	"github.com/koykov/pbtk/amq"
	"github.com/koykov/vmchain"
)

type mwAMQ struct {
	name string
}

func NewAMQ(name string) amq.MetricsWriter {
	return &mwAMQ{name: name}
}

func (mw *mwAMQ) Capacity(cap uint64) {
	vmchain.Gauge("amq_capacity", nil).WithLabel("name", mw.name).Set(float64(cap))
}

func (mw *mwAMQ) Set(err error) error {
	result := "success"
	if err != nil {
		result = "fail"
	}
	vmchain.Counter("amq_set").WithLabel("name", mw.name).
		WithLabel("result", result).Inc()
	return err
}

func (mw *mwAMQ) Unset(err error) error {
	result := "success"
	if err != nil {
		result = "fail"
	}
	vmchain.Counter("amq_unset").WithLabel("name", mw.name).
		WithLabel("result", result).Inc()
	return err
}

func (mw *mwAMQ) Contains(positive bool) bool {
	result := "positive"
	if !positive {
		result = "negative"
	}
	vmchain.Counter("amq_contains").WithLabel("name", mw.name).
		WithLabel("result", result).Inc()
	return positive
}

func (mw *mwAMQ) Reset() {}

var _ = NewAMQ
