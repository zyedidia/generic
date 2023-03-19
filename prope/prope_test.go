package prope_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"testing"

	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/prope"
	"github.com/zyedidia/generic/rope"
)

func check(p *prope.Node[byte], r *rope.Node[byte], t *testing.T) {
	if !bytes.Equal(r.Value(), p.Value()) {
		t.Errorf("incorrect bytes: %s %s", string(r.Value()), string(p.Value()))
	}
	if r.Len() != p.Len() {
		t.Errorf("incorrect length: %d %d", r.Len(), p.Len())
	}
}

const datasz = 5000

func data() (*prope.Node[byte], *rope.Node[byte]) {
	data := randbytes(datasz)
	p := prope.New(data)
	r := rope.New(data)
	return p, r
}

func randrange(high int) (int, int) {
	i1 := rand.Intn(high)
	i2 := rand.Intn(high)
	return g.Min(i1, i2), g.Max(i1, i2)
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randbytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}

func TestMain(m *testing.M) {
	rope.SplitLength = 4
	rope.JoinLength = 2

	os.Exit(m.Run())
}

func TestContruction(t *testing.T) {
	p, r := data()
	check(p, r, t)
}

func TestSimpleInsert(t *testing.T) {
	p := prope.New([]byte("0123456789"))
	p2 := p.Insert(5, []byte("abcd"))
	if !bytes.Equal(p.Value(), []byte("0123456789")) {
		t.Errorf("history not kept: %s 0123456789", string(p.Value()))
	}
	if !bytes.Equal(p2.Value(), []byte("01234abcd56789")) {
		t.Errorf("insertion failed: %s 01234abcd56789", string(p.Value()))
	}
}

func TestSimpleRemove(t *testing.T) {
	p := prope.New([]byte("0123456789"))
	p2 := p.Remove(5, 9)
	if !bytes.Equal(p.Value(), []byte("0123456789")) {
		t.Errorf("history not kept: %s 0123456789", string(p.Value()))
	}
	if !bytes.Equal(p2.Value(), []byte("012349")) {
		t.Errorf("removal failed: %s 012349", string(p.Value()))
	}
}

func TestSimpleSplit(t *testing.T) {
	p := prope.New([]byte("0123456789"))
	pl, pr := p.SplitAt(5)
	if !bytes.Equal(p.Value(), []byte("0123456789")) {
		t.Errorf("history not kept: %s 0123456789", string(p.Value()))
	}
	if !bytes.Equal(pl.Value(), []byte("01234")) {
		t.Errorf("split failed: %s 01234", string(p.Value()))
	}
	if !bytes.Equal(pr.Value(), []byte("56789")) {
		t.Errorf("split failed: %s 56789", string(p.Value()))
	}
}

func TestInsertRemove(t *testing.T) {
	const nedits = 100
	const strlen = 20
	propes := make([]*prope.Node[byte], nedits)
	ropes := make([]*rope.Node[byte], nedits)

	propes[0], ropes[0] = data()

	for i := 0; i < nedits-1; i++ {
		nextPrope := propes[i]
		nextRope := rope.New(ropes[i].Value())
		low, high := randrange(nextRope.Len())
		nextPrope = nextPrope.Remove(low, high)
		nextRope.Remove(low, high)
		bstr := randbytes(strlen)
		nextPrope = nextPrope.Insert(low, bstr)
		nextRope.Insert(low, bstr)
		propes[i+1] = nextPrope
		ropes[i+1] = nextRope
	}

	for i := 0; i < nedits; i++ {
		check(propes[i], ropes[i], t)
	}
}

func TestSlice(t *testing.T) {
	const nslice = 100
	p, r := data()

	for i := 0; i < nslice; i++ {
		low, high := randrange(p.Len())

		pb := p.Slice(low, high)
		rb := r.Slice(low, high)

		if !bytes.Equal(pb, rb) {
			t.Errorf("slice not equal: %s %s", string(pb), string(rb))
		}
	}
}

func TestSplit(t *testing.T) {
	const nsplits = 100
	const strlen = 20
	propes := make([]*prope.Node[byte], nsplits)
	ropes := make([]*rope.Node[byte], nsplits)

	propes[0], ropes[0] = data()

	for i := 0; i < nsplits-1; i++ {
		nextPrope := propes[i]
		nextRope := rope.New(ropes[i].Value())
		splitidx := rand.Intn(nextRope.Len())

		nextPrope, _ = nextPrope.SplitAt(splitidx)
		nextRope, _ = nextRope.SplitAt(splitidx)

		data := randbytes(strlen)
		nextPrope = nextPrope.Insert(0, data)
		nextRope.Insert(0, data)

		propes[i+1] = nextPrope
		ropes[i+1] = nextRope
	}

	for i := 0; i < nsplits; i++ {
		check(propes[i], ropes[i], t)
	}
}

func Example() {
	r := prope.New([]byte("hello world"))

	fmt.Println(string(r.At(0)))

	r2 := r.Remove(5, r.Len())
	r3 := r2.Insert(5, []byte(" rope"))

	fmt.Println(string(r.Value()))
	fmt.Println(string(r2.Value()))
	fmt.Println(string(r3.Value()))
	// Output:
	// h
	// hello world
	// hello
	// hello rope
}
