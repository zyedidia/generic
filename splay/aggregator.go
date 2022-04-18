package splay

import (
	g "github.com/zyedidia/generic"
)

// Aggregator is an interface for aggregating values from a range.
// V is the value type
// A is the structure stored at each node which is used for aggregation.
type Aggregator[V, A any] interface {
	// PopUp aggregates values from the children of the parent node.
	PopUp(self, lchd, rchd *A) *A

	// PushDown populates the updates attached at a node to its children.
	PushDown(self, lchd, rchd *A)

	// FromValue generates an A from a value.
	FromValue(value V) *A

	// Value returns the value stored in a node.
	Value(self *A) V
}

type valueAggregator[V any] struct{}

func (_ valueAggregator[V]) PopUp(self, lchd, rchd *V) *V {
	return self
}
func (_ valueAggregator[V]) PushDown(self, lchd, rchd *V) {}

func (_ valueAggregator[V]) FromValue(value V) *V {
	return &value
}

func (_ valueAggregator[V]) Value(self *V) V {
	if self != nil {
		return *self
	} else {
		var v V
		return v
	}
}

// NewValueAggregator creates a ValueAggregator, which only stores the value.
func NewValueAggregator[V any]() Aggregator[V, V] {
	return &valueAggregator[V]{}
}

type minMaxAgg[V any] struct {
	min   V
	max   V
	value V
}

type minMaxAggregator[V any] struct {
	less g.LessFn[V]
}

func (_ minMaxAggregator[V]) FromValue(value V) *minMaxAgg[V] {
	return &minMaxAgg[V]{
		min:   value,
		max:   value,
		value: value,
	}
}

func (_ minMaxAggregator[V]) PushDown(self, lchd, rchd *minMaxAgg[V]) {}

func (_ minMaxAggregator[V]) Value(self *minMaxAgg[V]) V {
	if self != nil {
		return self.value
	} else {
		var v V
		return v
	}
}

func (a minMaxAggregator[V]) PopUp(self, lchd, rchd *minMaxAgg[V]) *minMaxAgg[V] {
	if self == nil {
		return nil
	}
	self.min = self.value
	self.max = self.value
	if lchd != nil {
		if a.less(lchd.min, self.min) {
			self.min = lchd.min
		}
		if a.less(self.max, lchd.max) {
			self.max = lchd.max
		}
	}
	if rchd != nil {
		if a.less(rchd.min, self.min) {
			self.min = rchd.min
		}
		if a.less(self.max, rchd.max) {
			self.max = rchd.max
		}
	}
	return self
}

func (a *minMaxAgg[V]) Min() V {
	if a != nil {
		return a.min
	} else {
		var v V
		return v
	}
}

func (a *minMaxAgg[V]) Max() V {
	if a != nil {
		return a.max
	} else {
		var v V
		return v
	}
}

// NewMinMaxAggregator creates a MinMaxAggregator,
// which collects minimal and maximal values in a range.
func NewMinMaxAggregator[V any](less g.LessFn[V]) Aggregator[V, minMaxAgg[V]] {
	return &minMaxAggregator[V]{
		less: less,
	}
}

type rangeAssignAgg[V any] struct {
	value    V
	assgined *V
}

type rangeAssignAggregator[V any] struct{}

func (a *rangeAssignAgg[V]) Assign(value V) {
	if a != nil {
		a.value = value
		a.assgined = &a.value
	}
}

func (_ rangeAssignAggregator[V]) FromValue(value V) *rangeAssignAgg[V] {
	return &rangeAssignAgg[V]{
		value:    value,
		assgined: nil,
	}
}

func (_ rangeAssignAggregator[V]) PushDown(self, lchd, rchd *rangeAssignAgg[V]) {
	if self != nil && self.assgined != nil {
		if lchd != nil {
			lchd.value = *self.assgined
			lchd.assgined = &lchd.value
		}
		if rchd != nil {
			rchd.value = *self.assgined
			rchd.assgined = &rchd.value
		}
		self.assgined = nil
	}
}

func (_ rangeAssignAggregator[V]) PopUp(self, lchd, rchd *rangeAssignAgg[V]) *rangeAssignAgg[V] {
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

// NewRangeSetAggregator creates a RangeAssignAggregator,
// which can update the values associated to a range of keys.
func NewRangeAssignAggregator[V any]() Aggregator[V, rangeAssignAgg[V]] {
	return &rangeAssignAggregator[V]{}
}
