package diet

import (
	"testing"
)

func assert(b bool, s string, t *testing.T) {
	if !b {
		t.Fatalf("assertion failed: %s", s)
	}
}

func TestBasic(t *testing.T) {
	dt := NewTree[uint64]()

	dt.Put(10, 20)
	dt.Put(30, 40)

	assert(dt.Contains(35, 37), "35, 37", t)
	assert(dt.Contains(11, 12), "11, 12", t)
	assert(!dt.Contains(21, 29), "21, 29", t)
	assert(!dt.Contains(15, 35), "15 35", t)
	dt.Put(21, 25)
	assert(dt.Contains(21, 23), "21 23", t)
	assert(!dt.Contains(24, 28), "24, 28", t)
	assert(!dt.Contains(0, 5), "0, 5", t)
	dt.Remove(22, 23)
	assert(!dt.Contains(21, 23), "21, 23", t)
	assert(dt.Contains(18, 21), "18, 21", t)
	dt.Put(22, 23)
	dt.Put(26, 29)
	assert(dt.Contains(10, 29), "10, 29", t)
	assert(dt.Contains(10, 40), "10, 40", t)
}
