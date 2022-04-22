// Package aggregator provides interfaces of aggregator,
// which is used by data structures to aggregate values of a range.
package aggregator

import (
	g "github.com/zyedidia/generic"
)

// Aggregator is an interface for aggregating values from a range.
// V is the value type
// A is the structure stored at each node which is used for aggregation.
// R is the range-aggregated view type
type Aggregator[V, A, R any] interface {
	// PopUp aggregates values from the children of the parent node.
	// PopUp should be called on 'root' after any node in its subtree is updated,
	// including itself.
	PopUp(root *A, children []*A) *A

	// PushDown populates the updates attached at a node to its children.
	// PushDown should be called on 'root' every time before one of its
	// children is accessed, or 'root' is modified.
	PushDown(root *A, children []*A)

	// FromValue generates an A from a value.
	FromValue(value V) A

	// Value returns the value stored in 'node'.
	Value(node *A) V

	// RangeView returns the aggregated view of a range.
	// 'subTrees' contains the aggregators of the roots of subtrees
	// which are included in the range as a whole.
	// 'values' contains the aggregators of the nodes that are
	// included in the range themselves, excluding their children.
	RangeView(subTrees, values []*A) R
}

type valueAggregator[V any] struct{}

func (_ valueAggregator[V]) PopUp(self *V, children []*V) *V {
	return self
}
func (_ valueAggregator[V]) PushDown(self *V, children []*V) {}

func (_ valueAggregator[V]) FromValue(value V) V {
	return value
}

func (_ valueAggregator[V]) Value(self *V) V {
	if self != nil {
		return *self
	} else {
		var v V
		return v
	}
}

func (_ valueAggregator[V]) RangeView(nodes, values []*V) *V {
	return nil
}

// NewValueAggregator creates a ValueAggregator, which only stores the value.
func NewValueAggregator[V any]() Aggregator[V, V, *V] {
	return &valueAggregator[V]{}
}

type minMaxAgg[V any] struct {
	min   V
	max   V
	value V
}

type minMaxRV[V any] struct {
	min V
	max V
}

type minMaxAggregator[V any] struct {
	less g.LessFn[V]
}

func (_ minMaxAggregator[V]) FromValue(value V) minMaxAgg[V] {
	return minMaxAgg[V]{
		min:   value,
		max:   value,
		value: value,
	}
}

func (_ minMaxAggregator[V]) PushDown(self *minMaxAgg[V], children []*minMaxAgg[V]) {}

func (_ minMaxAggregator[V]) Value(self *minMaxAgg[V]) V {
	if self != nil {
		return self.value
	} else {
		var v V
		return v
	}
}

func (a minMaxAggregator[V]) PopUp(self *minMaxAgg[V], children []*minMaxAgg[V]) *minMaxAgg[V] {
	if self == nil {
		return nil
	}
	self.min = self.value
	self.max = self.value
	for _, chd := range children {
		if chd == nil {
			continue
		}
		if a.less(chd.min, self.min) {
			self.min = chd.min
		}
		if a.less(self.max, chd.max) {
			self.max = chd.max
		}
	}
	return self
}

func (a minMaxAggregator[V]) RangeView(nodes, values []*minMaxAgg[V]) *minMaxRV[V] {
	if len(nodes) == 0 && len(values) == 0 {
		return nil
	}
	ret := &minMaxRV[V]{}
	if len(nodes) > 0 {
		ret.min = nodes[0].min
		ret.max = nodes[0].max
	} else {
		ret.min = values[0].value
		ret.max = values[0].value
	}
	for _, n := range nodes {
		if n == nil {
			continue
		}
		if a.less(n.min, ret.min) {
			ret.min = n.min
		}
		if a.less(ret.max, n.max) {
			ret.max = n.max
		}
	}
	for _, v := range values {
		if v == nil {
			continue
		}
		if a.less(v.value, ret.min) {
			ret.min = v.value
		}
		if a.less(ret.max, v.value) {
			ret.max = v.max
		}
	}
	return ret
}

func (a *minMaxRV[V]) Min() V {
	if a != nil {
		return a.min
	} else {
		var v V
		return v
	}
}

func (a *minMaxRV[V]) Max() V {
	if a != nil {
		return a.max
	} else {
		var v V
		return v
	}
}

// NewMinMaxAggregator creates a MinMaxAggregator,
// which collects minimal and maximal values in a range.
func NewMinMaxAggregator[V any](less g.LessFn[V]) Aggregator[V, minMaxAgg[V], *minMaxRV[V]] {
	return &minMaxAggregator[V]{
		less: less,
	}
}

type rangeAssignAgg[V any] struct {
	value    V
	assgined *V
}

type rangeAssignRV[V any] struct {
	nodes  []*rangeAssignAgg[V]
	values []*rangeAssignAgg[V]
}

type rangeAssignAggregator[V any] struct{}

func (_ rangeAssignAggregator[V]) FromValue(value V) rangeAssignAgg[V] {
	return rangeAssignAgg[V]{
		value:    value,
		assgined: nil,
	}
}

func (_ rangeAssignAggregator[V]) PushDown(self *rangeAssignAgg[V], children []*rangeAssignAgg[V]) {
	if self != nil && self.assgined != nil {
		for _, chd := range children {
			if chd != nil {
				chd.value = *self.assgined
				chd.assgined = &chd.value
			}
		}
		self.assgined = nil
	}
}

func (_ rangeAssignAggregator[V]) PopUp(self *rangeAssignAgg[V], children []*rangeAssignAgg[V]) *rangeAssignAgg[V] {
	return self
}

func (_ rangeAssignAggregator[V]) Value(self *rangeAssignAgg[V]) V {
	if self != nil {
		return self.value
	} else {
		var v V
		return v
	}
}

func (a rangeAssignAggregator[V]) RangeView(nodes, values []*rangeAssignAgg[V]) rangeAssignRV[V] {
	return rangeAssignRV[V]{
		nodes:  nodes,
		values: values,
	}
}

// Assign assign 'value' to all nodes in the range.
func (a rangeAssignRV[V]) Assign(value V) {
	if a.nodes != nil {
		for _, n := range a.nodes {
			n.value = value
			n.assgined = &n.value
		}
	}
	if a.values != nil {
		for _, n := range a.values {
			n.value = value
		}
	}
}

// NewRangeSetAggregator creates a RangeAssignAggregator,
// which can update the values associated to a range of keys.
func NewRangeAssignAggregator[V any]() Aggregator[V, rangeAssignAgg[V], rangeAssignRV[V]] {
	return &rangeAssignAggregator[V]{}
}
