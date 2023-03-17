package diet

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Tree implements a discrete interval encoding tree (DIET). Intervals
// must be non-overlapping and will be merged together when possible.
// The DIET supports querying if a range is fully contained by the
// intervals that have been added to the tree.
type Tree[I constraints.Integer] struct {
	root *node[I]
}

// NewTree returns an empty DIET.
func NewTree[I constraints.Integer]() *Tree[I] {
	return &Tree[I]{}
}

// Contains checks if the range [start, end] is fully contained within
// the DIET.
func (t *Tree[I]) Contains(start, end I) bool {
	return t.root.contains(start, end)
}

// Put inserts a new interval [start, end].
func (t *Tree[I]) Put(start, end I) {
	t.root = insert(start, end, t.root)
}

// Remove the interval [start, end], which must be fully contained within the DIET.
func (t *Tree[I]) Remove(start, end I) {
	t.root = remove(start, end, t.root)
}

func (t *Tree[I]) dump() {
	dump(t.root, 0)
}

type node[I constraints.Integer] struct {
	start, end I

	left  *node[I]
	right *node[I]
}

func dump[I constraints.Integer](n *node[I], level int) {
	if n == nil {
		return
	}
	for i := 0; i < level; i++ {
		fmt.Print(" ")
	}
	fmt.Println("start:", n.start, "end:", n.end)
	dump(n.left, level+2)
	dump(n.right, level+2)
}

func (n *node[I]) icontains(val I) bool {
	return val >= n.start && val <= n.end
}

func (n *node[I]) contains(start, end I) bool {
	if n == nil {
		return false
	} else if n.icontains(start) && n.icontains(end) {
		return true
	} else if end < n.start {
		return n.left.contains(start, end)
	} else if start > n.end {
		return n.right.contains(start, end)
	}
	// TODO: is last branch needed?
	return false
}

func splitmax[I constraints.Integer](n *node[I]) (I, I, *node[I]) {
	if n.right == nil {
		return n.start, n.end, n.left
	}

	u, v, rp := splitmax(n.right)
	return u, v, &node[I]{
		start: n.start,
		end:   n.end,
		left:  n.left,
		right: rp,
	}
}

func splitmin[I constraints.Integer](n *node[I]) (I, I, *node[I]) {
	if n.left == nil {
		return n.start, n.end, n.right
	}

	u, v, lp := splitmin(n.left)
	return u, v, &node[I]{
		start: n.start,
		end:   n.end,
		left:  lp,
		right: n.right,
	}
}

func joinleft[I constraints.Integer](n node[I]) *node[I] {
	if n.left == nil {
		return &n
	}
	xp, yp, lp := splitmax(n.left)
	if yp+1 == n.start {
		return &node[I]{
			start: xp,
			end:   n.end,
			left:  lp,
			right: n.right,
		}
	}
	return &n
}

func joinright[I constraints.Integer](n node[I]) *node[I] {
	if n.right == nil {
		return &n
	}
	xp, yp, rp := splitmin(n.right)
	if n.end+1 == xp {
		return &node[I]{
			start: n.start,
			end:   yp,
			left:  n.left,
			right: rp,
		}
	}
	return &n
}

func insert[I constraints.Integer](zstart, zend I, n *node[I]) *node[I] {
	if n == nil {
		return &node[I]{
			start: zstart,
			end:   zend,
		}
	}
	if zend < n.start {
		if zend+1 == n.start {
			return joinleft(node[I]{
				start: zstart,
				end:   n.end,
				left:  n.left,
				right: n.right,
			})
		} else {
			return &node[I]{
				start: n.start,
				end:   n.end,
				left:  insert(zstart, zend, n.left),
				right: n.right,
			}
		}
	} else if zstart > n.end {
		if zstart == n.end+1 {
			return joinright(node[I]{
				start: n.start,
				end:   zend,
				left:  n.left,
				right: n.right,
			})
		} else {
			return &node[I]{
				start: n.start,
				end:   n.end,
				left:  n.left,
				right: insert(zstart, zend, n.right),
			}
		}
	} else {
		return n
	}
}

func merge[I constraints.Integer](l *node[I], r *node[I]) *node[I] {
	if r == nil {
		return l
	}
	if l == nil {
		return r
	}
	x, y, lp := splitmax(l)
	return &node[I]{
		start: x,
		end:   y,
		left:  lp,
		right: r,
	}
}

func remove[I constraints.Integer](zstart, zend I, n *node[I]) *node[I] {
	if n == nil {
		return nil
	}
	if zend < n.start {
		return &node[I]{
			start: n.start,
			end:   n.end,
			left:  remove(zstart, zend, n.left),
			right: n.right,
		}
	} else if zstart > n.end {
		return &node[I]{
			start: n.start,
			end:   n.end,
			left:  n.left,
			right: remove(zstart, zend, n.right),
		}
	} else if zstart == n.start {
		if zend == n.end {
			return merge(n.left, n.right)
		} else {
			return &node[I]{
				start: zend + 1,
				end:   n.end,
				left:  n.left,
				right: n.right,
			}
		}
	} else if zend == n.end {
		return &node[I]{
			start: n.start,
			end:   zstart - 1,
			left:  n.left,
			right: n.right,
		}
	} else {
		return &node[I]{
			start: n.start,
			end:   zstart - 1,
			left:  n.left,
			right: &node[I]{
				start: zend + 1,
				end:   n.end,
				left:  nil,
				right: n.right,
			},
		}
	}
}
