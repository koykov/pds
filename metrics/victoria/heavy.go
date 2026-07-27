package victoria

import (
	"math"

	"github.com/koykov/pbtk/heavy"
	"github.com/koykov/vmchain"
)

type mwHeavy struct {
	name string
}

func NewHeavy(name string) heavy.MetricsWriter {
	return &mwHeavy{name: name}
}

func (mw *mwHeavy) Add(err error) error {
	result := "success"
	if err != nil {
		result = "fail"
	}
	vmchain.Counter("heavy_items").
		WithLabel("name", mw.name).
		WithLabel("result", result).Inc()
	return err
}

func (mw *mwHeavy) Hits(hits []heavy.Freq) {
	n := len(hits)
	if n == 0 {
		return
	}

	// base metrics
	min_, max_ := hits[n-1].Freq(), hits[0].Freq()
	var sum float64
	for i := 0; i < n; i++ {
		f := hits[i].Freq()
		sum += f
		vmchain.Histogram("heavy_freq_distribution").WithLabel("name", mw.name).Update(f)
	}
	mean := sum / float64(n)
	vmchain.Gauge("heavy_freq_min", nil).WithLabel("name", mw.name).Set(min_)
	vmchain.Gauge("heavy_freq_max", nil).WithLabel("name", mw.name).Set(max_)
	vmchain.Gauge("heavy_freq_mean", nil).WithLabel("name", mw.name).Set(mean)
	vmchain.Gauge("heavy_freq_sum", nil).WithLabel("name", mw.name).Set(sum)

	// variance,  stddev and relative metrics
	var variance, skewness float64
	for i := 0; i < n; i++ {
		diff := hits[i].Freq() - mean
		variance += diff * diff
		skewness += math.Pow(diff, 3)
	}
	variance /= float64(n)
	stddev := math.Sqrt(variance)
	skewness /= float64(n) * math.Pow(stddev, 3)
	coefVariation := stddev / mean
	vmchain.Gauge("heavy_freq_stddev", nil).WithLabel("name", mw.name).Set(stddev)
	vmchain.Gauge("heavy_freq_variance", nil).WithLabel("name", mw.name).Set(variance)
	vmchain.Gauge("heavy_freq_skewness", nil).WithLabel("name", mw.name).Set(skewness)
	vmchain.Gauge("heavy_freq_cvar", nil).WithLabel("name", mw.name).Set(coefVariation)

	// percentiles
	vmchain.Gauge("hitter_freq_percentile", nil).WithLabel("name", mw.name).WithLabel("p", "25").Set(hits[int(0.75*float64(n))].Freq())
	vmchain.Gauge("hitter_freq_percentile", nil).WithLabel("name", mw.name).WithLabel("p", "50").Set(hits[int(0.50*float64(n))].Freq())
	vmchain.Gauge("hitter_freq_percentile", nil).WithLabel("name", mw.name).WithLabel("p", "75").Set(hits[int(0.25*float64(n))].Freq())
	vmchain.Gauge("hitter_freq_percentile", nil).WithLabel("name", mw.name).WithLabel("p", "90").Set(hits[int(0.10*float64(n))].Freq())
	vmchain.Gauge("hitter_freq_percentile", nil).WithLabel("name", mw.name).WithLabel("p", "99").Set(hits[int(0.01*float64(n))].Freq())
}

func (mw *mwHeavy) Reset() {
	// ...
}

var _ = NewHeavy
