// Package generic provides types and constraints useful for implementing
// generic data structures. In particular, wrappers of primitive types are
// provided so that they implement Lesser, Comparable, and Hashable interfaces.
// This allows generic data structures that can use primitive types or custom
// user types. This package uses a custom murmur hash function for integer
// types and FNV1a for strings.
package generic

import (
	"constraints"

	"github.com/segmentio/fasthash/fnv1a"
)

func Equals[T comparable](a, b T) bool {
	return a == b
}

func Less[T constraints.Ordered](a, b T) bool {
	return a < b
}

func HashUint64(u uint64) uint64 {
	return hash(u)
}
func HashUint32(u uint32) uint64 {
	return hash(uint64(u))
}
func HashUint16(u uint16) uint64 {
	return hash(uint64(u))
}
func HashUint8(u uint8) uint64 {
	return hash(uint64(u))
}
func HashInt64(i int64) uint64 {
	return hash(uint64(i))
}
func HashInt32(i int32) uint64 {
	return hash(uint64(i))
}
func HashInt16(i int32) uint64 {
	return hash(uint64(i))
}
func HashInt8(i int32) uint64 {
	return hash(uint64(i))
}
func HashInt(i int) uint64 {
	return hash(uint64(i))
}
func HashUint(i uint) uint64 {
	return hash(uint64(i))
}
func HashString(s string) uint64 {
	return fnv1a.HashString64(s)
}
func HashBytes(b []byte) uint64 {
	return fnv1a.HashBytes64(b)
}

type Equaler[T any] func(a, b T) bool
type Lesser[T any] func(a, b T) bool
type Hasher[T any] func(t T) uint64

func hash(u uint64) uint64 {
	u ^= u >> 33
	u *= 0xff51afd7ed558ccd
	u ^= u >> 33
	u *= 0xc4ceb9fe1a85ec53
	u ^= u >> 33
	return u
}
