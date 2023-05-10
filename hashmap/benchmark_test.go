// package bench is inspired by https://tessil.github.io/2016/08/29/benchmark-hopscotch-map.html
package hashmap_test

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/hashmap"
)

func getRanges() []int {
	r := os.Getenv("RANGES")
	if r == "" {
		r = "200000 400000 600000 800000 1000000 1200000 1400000 1600000 1800000 2000000 2200000 2400000 2600000 2800000 3000000"
	}
	rangesStr := strings.Split(r, " ")
	rangesInt := make([]int, len(rangesStr))

	for i := range rangesStr {
		var err error
		rangesInt[i], err = strconv.Atoi(rangesStr[i])
		if err != nil {
			panic(err)
		}
	}
	return rangesInt
}

//go:noinline
func handleElem(key uint64, val uint64) {}

func getMapNames() []string {
	m := os.Getenv("MAPS")
	if m == "" {
		m = "std linear robin"
	}
	return strings.Split(m, " ")
}

func createMapUint64(n int, mapName string) hashmap.HashMap[uint64, uint64] {
	switch mapName {
	case "std":
		m := make(map[uint64]uint64, n)
		return hashmap.HashMap[uint64, uint64]{
			Put: func(k uint64, v uint64) {
				m[k] = v
			},
			Get: func(k uint64) (uint64, bool) {
				v, ok := m[k]
				return v, ok
			},
			Remove: func(k uint64) {
				delete(m, k)
			},
			Each: func(callback func(key uint64, val uint64)) {
				for k, v := range m {
					callback(k, v)
				}
			},
			Load: func() float64 {
				return -1.0 //unknown
			},
		}
	case "robin":
		m := hashmap.NewRobinMapWithHasher[uint64, uint64](g.HashUint64)
		m.Reserve(uintptr(n))
		return hashmap.HashMap[uint64, uint64]{
			Get:     m.Get,
			Reserve: m.Reserve,
			Put:     m.Put,
			Remove:  m.Remove,
			Clear:   m.Clear,
			Size:    m.Size,
			Each:    m.Each,
			Load:    m.Load,
		}
	case "linear":
		m := hashmap.New[uint64, uint64](uint64(n), g.Equals[uint64], g.HashUint64)
		//m.Reserve(uintptr(n))
		return hashmap.HashMap[uint64, uint64]{
			Get:     m.Get,
			Reserve: m.Reserve,
			Put:     m.Put,
			Remove:  m.Remove,
			Clear:   m.Clear,
			Size:    m.Size,
			Each:    m.Each,
			Load:    m.Load,
		}
	default:
		panic(fmt.Sprintln("unknown map:", mapName))
	}
}

func genRandUInt64Array(n int) []uint64 {
	arr := make([]uint64, n)
	for i := range arr {
		arr[i] = rand.Uint64()
	}
	return arr
}

func genShuffled64Array(n int) []uint64 {
	arr := make([]uint64, n)
	for i := range arr {
		arr[i] = uint64(i)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func genDifferentRandUInt64Array(in []uint64) []uint64 {
	out := make([]uint64, len(in))
	sort.Slice(in, func(i int, j int) bool { return in[i] < in[j] })

	for j := 0; j < len(out); {
		x := rand.Uint64()
		_, found := sort.Find(len(in), func(i int) int {
			return int(x - in[i])
		})
		if !found {
			out[j] = x
			j++
		}
	}

	return out
}

func report(b *testing.B, n int, load float64) {
	b.ReportAllocs()
	b.ReportMetric(float64(n), "N-runs")
	b.ReportMetric(load, "Load")
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	b.ReportMetric(float64(mem.Alloc), "Bytes")
}

// Before the test, a vector with the values [0, N) is generated and shuffled.
// Then for each value in the vector, the key-value pair (k, 1) is inserted into the hash map.
func BenchmarkRandomShuffleInsertsU64(b *testing.B) {
	for _, mapName := range getMapNames() {
		for _, r := range getRanges() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := -1.0
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMapUint64(0, mapName)
					arr := genShuffled64Array(r)

					b.StartTimer()
					for j := range arr {
						m.Put(arr[j], 1)
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, a vector with random values in range [0, 2^64-1) is generated.
// Then for each value in the vector, the key-value pair (k, 1) is inserted into the hash map.
func BenchmarkRandomFullInsertsInsertsU64(b *testing.B) {
	for _, mapName := range getMapNames() {
		for _, r := range getRanges() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := -1.0
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMapUint64(0, mapName)
					arr := genRandUInt64Array(r)

					b.StartTimer()
					for j := range arr {
						m.Put(arr[j], 1)
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Same as the random full inserts test but the reserve method of the hash map is called beforehand
// to avoid any rehash during the insertion. It provides a fair comparison even if the growth factor
// of each hash map is different.
func BenchmarkRandomFullInsertsWithReserveInsertsU64(b *testing.B) {
	for _, mapName := range getMapNames() {
		for _, r := range getRanges() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := -1.0
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMapUint64(r, mapName)
					arr := genRandUInt64Array(r)

					b.StartTimer()
					for j := range arr {
						m.Put(arr[j], 1)
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n elements in the same way as in the random full insert test are added.
// Each key is deleted one by one in a different and random order than the one they were inserted.
func BenchmarkRandomFullDeletesU64(b *testing.B) {
	for _, mapName := range getMapNames() {
		for _, r := range getRanges() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := -1.0
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMapUint64(r, mapName)
					arr := genRandUInt64Array(r)
					for j := range arr {
						m.Put(arr[j], 1)
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })

					b.StartTimer()
					for j := range arr {
						m.Remove(arr[j])
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n elements are inserted in the same way as in the random shuffle inserts test.
// Read each key-value pair is look up in a different and random order than the one they were inserted.
func BenchmarkRandomShuffleReadsU64(b *testing.B) {
	for _, mapName := range getMapNames() {
		for _, r := range getRanges() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := -1.0
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMapUint64(0, mapName)
					arr := genShuffled64Array(r)
					for j := range arr {
						m.Put(arr[j], 1)
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })

					b.StartTimer()
					for j := range arr {
						_, found := m.Get(arr[j])
						if !found {
							b.Fatal("inserted key not found")
						}
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n elements are inserted in the same way as in the random full inserts test.
// Read each key-value pair is look up in a different and random order than the one they were inserted.
func BenchmarkFullReadsU64(b *testing.B) {
	for _, mapName := range getMapNames() {
		for _, r := range getRanges() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := -1.0
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMapUint64(0, mapName)
					arr := genRandUInt64Array(r)
					for j := range arr {
						m.Put(arr[j], 1)
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })

					b.StartTimer()
					for j := range arr {
						_, found := m.Get(arr[j])
						if !found {
							b.Fatal("inserted key not found")
						}
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n elements are inserted in the same way as in the random full inserts test.
// Then a another vector of n random elements different from the inserted elements is generated
// which is tried to search in the hash map.
func BenchmarkFullReadsMissesU64(b *testing.B) {
	for _, mapName := range getMapNames() {
		for _, r := range getRanges() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := -1.0
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMapUint64(0, mapName)
					arr := genRandUInt64Array(r)
					for j := range arr {
						m.Put(arr[j], 1)
					}

					other := genDifferentRandUInt64Array(arr)
					rand.Shuffle(len(other), func(i, j int) { other[i], other[j] = other[j], other[i] })

					b.StartTimer()
					for j := range arr {
						_, found := m.Get(other[j])
						if found {
							b.Fatal("missed key was found")
						}
					}
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, we insert n elements in the same way as in the random full inserts test
// before deleting half of these values randomly. We then try to read all the original values
// in a different order which will lead to 50% hits and 50% misses.
func BenchmarkRandomFullReadsAfterDeletingHalf(b *testing.B) {
	for _, mapName := range getMapNames() {
		for _, r := range getRanges() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := -1.0
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMapUint64(0, mapName)
					arr := genRandUInt64Array(r)
					for j := range arr {
						m.Put(arr[j], 1)
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
					numRemoved := len(arr) / 2
					for j := 0; j < numRemoved; j++ {
						m.Remove(arr[j])
					}
					rand.Shuffle(len(arr), func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })

					b.StartTimer()
					ac := 0
					for j := range arr {
						_, found := m.Get(arr[j])
						x := 0
						if !found {
							x = 1
						}
						ac = ac + x
					}
					b.StopTimer()
					if ac != numRemoved {
						b.Fatal("unexpected lookup accumulation:", ac)
					}

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}

// Before the test, n elements are inserted in the same way as in the random full inserts test.
// Then use the hash map iterators to read all the key-value pairs.
func BenchmarkRandomFullIteration(b *testing.B) {
	for _, mapName := range getMapNames() {
		for _, r := range getRanges() {
			b.Run(fmt.Sprintf("%s-%d", mapName, r), func(b *testing.B) {
				load := -1.0
				for i := 0; i < b.N; i++ {
					b.StopTimer()

					m := createMapUint64(0, mapName)
					arr := genRandUInt64Array(r)
					for j := range arr {
						m.Put(arr[j], 1)
					}

					b.StartTimer()
					m.Each(handleElem)
					b.StopTimer()

					load = m.Load()
				}
				report(b, r, load)
			})
		}
	}
}
