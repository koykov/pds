package prometheus

import (
	"math"

	"github.com/koykov/pbtk/heavy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	heavyItems.WithLabelValues(mw.name, result).Inc()
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
		heavyFreqDist.WithLabelValues(mw.name).Observe(f)
	}
	mean := sum / float64(n)
	heavyFreqMin.WithLabelValues(mw.name).Set(min_)
	heavyFreqMax.WithLabelValues(mw.name).Set(max_)
	heavyFreqMean.WithLabelValues(mw.name).Set(mean)
	heavyFreqSum.WithLabelValues(mw.name).Set(sum)

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
	heavyFreqStddev.WithLabelValues(mw.name).Set(stddev)
	heavyFreqVariance.WithLabelValues(mw.name).Set(variance)
	heavyFreqSkewness.WithLabelValues(mw.name).Set(skewness)
	heavyFreqCVar.WithLabelValues(mw.name).Set(coefVariation)

	// percentiles
	heavyPercentiles.WithLabelValues(mw.name, "25").Set(hits[int(0.75*float64(n))].Freq())
	heavyPercentiles.WithLabelValues(mw.name, "50").Set(hits[int(0.50*float64(n))].Freq())
	heavyPercentiles.WithLabelValues(mw.name, "75").Set(hits[int(0.25*float64(n))].Freq())
	heavyPercentiles.WithLabelValues(mw.name, "90").Set(hits[int(0.10*float64(n))].Freq())
	heavyPercentiles.WithLabelValues(mw.name, "99").Set(hits[int(0.01*float64(n))].Freq())
}

func (mw *mwHeavy) Reset() {
	// ...
}

func init() {
	heavyItems = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "heavy_items",
		Help: "Total items in heavy hitters",
	}, []string{"name", "result"})

	heavyFreqMin = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_freq_min",
		Help: "Minimum frequency in top list",
	}, []string{"name"})

	heavyFreqMax = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_freq_max",
		Help: "Maximum frequency in top list",
	}, []string{"name"})

	heavyFreqMean = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_freq_mean",
		Help: "Mean frequency in top list",
	}, []string{"name"})

	heavyFreqSum = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_freq_sum",
		Help: "Sum of all frequencies in top list",
	}, []string{"name"})

	heavyFreqStddev = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_freq_stddev",
		Help: "Standard deviation of frequencies",
	}, []string{"name"})

	heavyFreqVariance = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_freq_variance",
		Help: "Variance of frequencies",
	}, []string{"name"})

	heavyFreqSkewness = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_freq_skewness",
		Help: "Skewness of frequency distribution",
	}, []string{"name"})

	heavyFreqCVar = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_freq_cvar",
		Help: "Coefficient of variation",
	}, []string{"name"})

	heavyErrorBound = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_error_bound",
		Help: "Theoretical error bound",
	}, []string{"name"})

	heavyObservedMinError = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_observed_min_error",
		Help: "Observed min frequency / total elements",
	}, []string{"name"})

	heavyPercentiles = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "hitter_freq_percentile",
		Help: "Frequency percentiles in top list",
	}, []string{"name", "p"},
	)

	heavyFreqDist = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "heavy_freq_distribution",
		Help:    "Distribution of frequencies",
		Buckets: prometheus.ExponentialBuckets(1, 2, 10), // 1, 2, 4, 8,..., 512
	}, []string{"name"})

	prometheus.MustRegister(heavyItems, heavyFreqMin, heavyFreqMax, heavyFreqMean, heavyFreqSum, heavyFreqStddev,
		heavyFreqVariance, heavyFreqSkewness, heavyFreqCVar, heavyErrorBound, heavyObservedMinError, heavyPercentiles,
		heavyFreqDist)
}

var (
	heavyItems            *prometheus.CounterVec // total items
	heavyFreqMin          *prometheus.GaugeVec
	heavyFreqMax          *prometheus.GaugeVec
	heavyFreqMean         *prometheus.GaugeVec
	heavyFreqSum          *prometheus.GaugeVec
	heavyFreqStddev       *prometheus.GaugeVec
	heavyFreqVariance     *prometheus.GaugeVec
	heavyFreqSkewness     *prometheus.GaugeVec
	heavyFreqCVar         *prometheus.GaugeVec
	heavyErrorBound       *prometheus.GaugeVec
	heavyObservedMinError *prometheus.GaugeVec
	heavyPercentiles      *prometheus.GaugeVec
	heavyFreqDist         *prometheus.HistogramVec

	_ = NewHeavy
)
