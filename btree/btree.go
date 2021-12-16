package btree

import (
	g "github.com/zyedidia/generic"
	"github.com/zyedidia/generic/iter"
)

const maxChildren = 64 // must be even and greater than 2

type KV[K g.Lesser[K], V any] struct {
	Key K
	Val V
}

type Tree[K g.Lesser[K], V any] struct {
	root   *node[K, V]
	height int
	n      int
}

type node[K g.Lesser[K], V any] struct {
	m        int
	children [maxChildren]entry[K, V]
}

type entry[K g.Lesser[K], V any] struct {
	key  K
	val  V
	next *node[K, V]
}

func NewTree[K g.Lesser[K], V any]() *Tree[K, V] {
	return &Tree[K, V]{
		root: &node[K, V]{},
	}
}

func (t *Tree[K, V]) Size() int {
	return t.n
}

func (t *Tree[K, V]) Get(key K) (V, bool) {
	return t.search(t.root, key, t.height)
}

func (t *Tree[K, V]) GetZ(key K) V {
	v, _ := t.search(t.root, key, t.height)
	return v
}

func (t *Tree[K, V]) search(x *node[K, V], key K, height int) (V, bool) {
	children := x.children

	if height == 0 {
		// leaf node
		for j := 0; j < x.m; j++ {
			if g.Compare(key, children[j].key) == 0 {
				return children[j].val, true
			}
		}
	} else {
		// internal node
		for j := 0; j < x.m; j++ {
			if x.m == j+1 || g.Compare(key, children[j+1].key) < 0 {
				return t.search(children[j].next, key, height-1)
			}
		}
	}
	var v V
	return v, false
}

func (t *Tree[K, V]) Put(key K, val V) {
	u := t.insert(t.root, key, val, t.height)
	t.n++
	if u == nil {
		return
	}

	n := &node[K, V]{
		m: 2,
	}
	n.children[0] = entry[K, V]{
		key:  t.root.children[0].key,
		next: t.root,
	}
	n.children[1] = entry[K, V]{
		key:  u.children[0].key,
		next: u,
	}
	t.root = n
	t.height++
}

func (t *Tree[K, V]) insert(h *node[K, V], key K, val V, height int) *node[K, V] {
	ent := entry[K, V]{
		key: key,
		val: val,
	}

	var j int
	if height == 0 {
		// leaf node
		for j = 0; j < h.m; j++ {
			if g.Compare(key, h.children[j].key) < 0 {
				break
			}
		}
	} else {
		// internal node
		for j = 0; j < h.m; j++ {
			if (j+1 == h.m) || g.Compare(key, h.children[j+1].key) < 0 {
				u := t.insert(h.children[j].next, key, val, height-1)
				if u == nil {
					return nil
				}
				j++
				ent.key = u.children[0].key
				ent.next = u
				break
			}
		}
	}

	for i := h.m; i > j; i-- {
		h.children[i] = h.children[i-1]
	}
	h.children[j] = ent
	h.m++
	if h.m < maxChildren {
		return nil
	}
	return t.split(h)
}

func (t *Tree[K, V]) split(h *node[K, V]) *node[K, V] {
	n := &node[K, V]{
		m: maxChildren / 2,
	}
	h.m = maxChildren / 2
	for j := 0; j < maxChildren/2; j++ {
		n.children[j] = h.children[maxChildren/2+j]
	}
	return n
}

func (t *Tree[K, V]) Iter() iter.Iter[KV[K, V]] {
	return t.iter(t.root)
}

func (t *Tree[K, V]) iter(n *node[K, V]) iter.Iter[KV[K, V]] {
	if n == nil {
		return func() (v KV[K, V], ok bool) {
			return v, false
		}
	}
}
