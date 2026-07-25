package hyperloglog

import (
	"bytes"
	_ "embed"
	"encoding/binary"
	"math"
)

//go:embed bias.bin
var biasRaw []byte

// empirical bias correction pairs
// loads from local binary due to huge size
var bias [15][][2]float64

func init() {
	var err error
	r := bytes.NewBuffer(biasRaw)
	for i := 0; i < 15; i++ {
		var buf [8]byte
		if _, err = r.Read(buf[:]); err != nil {
			return
		}
		n := binary.LittleEndian.Uint64(buf[:])
		for j := uint64(0); j < n; j++ {
			if _, err = r.Read(buf[:]); err != nil {
				return
			}
			dist := binary.LittleEndian.Uint64(buf[:])
			if _, err = r.Read(buf[:]); err != nil {
				return
			}
			bias_ := binary.LittleEndian.Uint64(buf[:])
			bias[i] = append(bias[i], [2]float64{math.Float64frombits(dist), math.Float64frombits(bias_)})
		}
	}
}

type biasp struct{ d, e float64 }

func biasEstimation(p uint64, e float64) float64 {
	_ = bias[14]
	const k = 6 // K-nn
	var keys [6]biasp
	ssize := len(bias[p])
	var eidx int
	{
		// std::lower_bound
		for i := 0; i < ssize; i++ {
			if bias[p][i][0] < e {
				eidx = i
				continue
			}
			break
		}
	}
	{
		// std::partial_sort_copy
		lo, hi := 0, k
		if eidx > k {
			lo, hi = eidx-k, eidx
		} else if eidx+k <= ssize {
			lo, hi = eidx, eidx+k
		}
		for i := lo; i < hi; i++ {
			keys[i-lo].d, keys[i-lo].e = bias[p][i][0], bias[p][i][1]
		}
		biasQsort(keys[:], 0, k-1, e)
	}
	var s, ws float64
	{
		// std::accumulate
		for i := 0; i < k; i++ {
			s += keys[i].e * 1 / (math.Abs(keys[i].d-e) + 1e-5)
			ws += 1 / (math.Abs(keys[i].d-e) + 1e-5)
		}
	}
	return s / ws
}

func biasPivot(p []biasp, lo, hi int, e float64) int {
	if len(p) == 0 {
		return 0
	}
	pi := &p[hi]
	i := lo - 1
	_ = p[len(p)-1]
	for j := lo; j <= hi-1; j++ {
		if math.Abs(p[j].d-e) < math.Abs(pi.d-e) {
			i++
			p[i], p[j] = p[j], p[i]
		}
	}
	if i < hi {
		p[i+1], p[hi] = p[hi], p[i+1]
	}
	return i + 1
}

func biasQsort(p []biasp, lo, hi int, e float64) {
	if lo < hi {
		pi := biasPivot(p, lo, hi, e)
		biasQsort(p, lo, pi-1, e)
		biasQsort(p, pi+1, hi, e)
	}
}
