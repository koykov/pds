package victoria

import (
	"github.com/koykov/pbtk/cardinality"
	"github.com/koykov/vmchain"
)

type mwCardinality struct {
	name string
}

func NewCardinality(name string) cardinality.MetricsWriter {
	return &mwCardinality{name: name}
}

func (mw *mwCardinality) Add(err error) error {
	result := "success"
	if err != nil {
		result = "fail"
	}
	vmchain.Counter("cardinality_add").WithLabel("name", mw.name).
		WithLabel("result", result).Inc()
	return err
}

func (mw *mwCardinality) Estimate(value uint64) uint64 {
	vmchain.Histogram("cardinality_unique").WithLabel("name", mw.name).Update(float64(value))
	return value
}

var _ = NewCardinality
