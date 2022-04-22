package aggregator_test

import (
	"math/rand"
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/aggregator"
)

func TestValueAggregator(t *testing.T) {
	aggarr := aggregator.NewArray(100, aggregator.NewValueAggregator[int]())
	stdarr := make([]int, 100)

	const nops = 2000
	for i := 0; i < nops; i++ {
		key := rand.Intn(100)
		val := rand.Intn(0xffff)
		op := rand.Intn(2)

		switch op {
		case 0:
			stdarr[key] = val
			aggarr.Put(key, val)
		case 1:
			stdv := stdarr[key]
			aggv := aggarr.Get(key)
			if stdv != aggv {
				t.Fatalf("value mismatch: %v != %v", stdv, aggv)
			}
		}
	}
}

func TestMinMaxAggregator(t *testing.T) {
	const length = 10
	aggarr := aggregator.NewArray(length, aggregator.NewMinMaxAggregator(g.Less[int]))
	stdarr := make([]int, length)

	const nops = 3000
	for i := 0; i < nops; i++ {
		key := rand.Intn(length)
		val := rand.Intn(0xffff)
		op := rand.Intn(3)

		switch op {
		case 0:
			stdarr[key] = val
			aggarr.Put(key, val)
		case 1:
			stdv := stdarr[key]
			aggv := aggarr.Get(key)
			if stdv != aggv {
				t.Fatalf("value mismatch: %v != %v", stdv, aggv)
			}
		case 2:
			l := key
			r := rand.Intn(length)
			if l > r {
				l, r = r, l
			}
			r += 1
			agg := aggarr.Range(l, r)
			if agg == nil {
				continue
			}
			aggMin := agg.Min()
			aggMax := agg.Max()
			stdMin := 0x10000
			stdMax := 0
			for i := l; i < r; i++ {
				stdMin = g.Min(stdMin, stdarr[i])
				stdMax = g.Max(stdMax, stdarr[i])
			}
			if aggMin != stdMin || aggMax != stdMax {
				t.Fatalf("value mismatch: (%v, %v) != (%v, %v)", aggMin, aggMax, stdMin, stdMax)
			}
		}
	}
}

func TestRangeAssignAggregator(t *testing.T) {
	const length = 10
	aggarr := aggregator.NewArray(length, aggregator.NewRangeAssignAggregator[int]())
	stdarr := make([]int, length)

	const nops = 4000
	for i := 0; i < nops; i++ {
		key := rand.Intn(length)
		val := rand.Intn(0xffff)
		op := rand.Intn(3)

		switch op {
		case 0:
			stdarr[key] = val
			aggarr.Put(key, val)
		case 1:
			stdv := stdarr[key]
			aggv := aggarr.Get(key)
			if stdv != aggv {
				t.Fatalf("value mismatch: %v != %v", stdv, aggv)
			}
		case 2:
			l := key
			r := rand.Intn(length)
			if l > r {
				l, r = r, l
			}
			r += 1
			aggarr.Range(l, r).Assign(val)
			for i := l; i < r; i++ {
				stdarr[i] = val
			}
		}
	}

	for i := 0; i < length; i++ {
		if stdarr[i] != aggarr.Get(i) {
			t.Fatalf("value mismatch: %v != %v", stdarr[i], aggarr.Get(i))
		}
	}
}
