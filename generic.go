package generic

import (
	"fmt"
	"reflect"
	"unsafe"

	"golang.org/x/exp/constraints"

	"github.com/segmentio/fasthash/fnv1a"
)

var (
	tab64 = []int{
		63, 0, 58, 1, 59, 47, 53, 2,
		60, 39, 48, 27, 54, 33, 42, 3,
		61, 51, 37, 40, 49, 18, 28, 20,
		55, 30, 34, 11, 43, 14, 22, 4,
		62, 57, 46, 52, 38, 26, 32, 41,
		50, 36, 17, 19, 29, 10, 13, 21,
		56, 45, 25, 31, 35, 16, 9, 12,
		44, 24, 15, 8, 23, 7, 6, 5,
	}
)

// EqualsFn is a function that returns whether 'a' and 'b' are equal.
type EqualsFn[T any] func(a, b T) bool

// LessFn is a function that returns whether 'a' is less than 'b'.
type LessFn[T any] func(a, b T) bool

// HashFn is a function that returns the hash of 't'.
type HashFn[T any] func(t T) uint64

// Equals wraps the '==' operator for comparable types.
func Equals[T comparable](a, b T) bool {
	return a == b
}

// Less wraps the '<' operator for ordered types.
func Less[T constraints.Ordered](a, b T) bool {
	return a < b
}

// Compare uses a less function to determine the ordering of 'a' and 'b'. It returns:
//
// * -1 if a < b
//
// * 1 if a > b
//
// * 0 if a == b
func Compare[T any](a, b T, less LessFn[T]) int {
	if less(a, b) {
		return -1
	} else if less(b, a) {
		return 1
	}
	return 0
}

// NextPowerOf2 is a fast calculation implementation of 2^x.
// see: https://stackoverflow.com/questions/466204/rounding-up-to-next-power-of-2
// go:inline
func NextPowerOf2(i uint64) uint64 {
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32
	i++
	return i
}

// Log2 is fast calculation implementation of log2(x)
// see: https://stackoverflow.com/questions/11376288/fast-computing-of-log2-for-64-bit-integers
// go:inline
func Log2(value uint64) uint64 {
	value |= value >> 1
	value |= value >> 2
	value |= value >> 4
	value |= value >> 8
	value |= value >> 16
	value |= value >> 32

	index := ((value - (value >> 1)) * 0x07EDD5E59A4E28C2) >> 58
	return uint64(tab64[index])
}

// Swap exchange the values of the two given pointers.
// go:inline
func Swap[T any](a, b *T) {
	tmp := *a
	*a = *b
	*b = tmp
}

// Max returns the max of a and b.
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min returns the min of a and b.
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Clamp returns x constrained within [lo:hi] range.
// If x compares less than lo, returns lo; otherwise if hi compares less than x, returns hi; otherwise returns v.
func Clamp[T constraints.Ordered](x, lo, hi T) T {
	return Max(lo, Min(hi, x))
}

// MaxFunc returns the max of a and b using the less func.
func MaxFunc[T any](a, b T, less LessFn[T]) T {
	if less(b, a) {
		return a
	}
	return b
}

// MinFunc returns the min of a and b using the less func.
func MinFunc[T any](a, b T, less LessFn[T]) T {
	if less(a, b) {
		return a
	}
	return b
}

// ClampFunc returns x constrained within [lo:hi] range using the less func.
// If x compares less than lo, returns lo; otherwise if hi compares less than x, returns hi; otherwise returns v.
func ClampFunc[T any](x, lo, hi T, less LessFn[T]) T {
	return MaxFunc(lo, MinFunc(hi, x, less), less)
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
func HashInt16(i int16) uint64 {
	return hash(uint64(i))
}
func HashFloat32(i float32) uint64 {
	return hash(uint64(i))
}
func HashFloat64(i float64) uint64 {
	return hash(uint64(i))
}
func HashInt8(i int8) uint64 {
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

func hash(u uint64) uint64 {
	u ^= u >> 33
	u *= 0xff51afd7ed558ccd
	u ^= u >> 33
	u *= 0xc4ceb9fe1a85ec53
	u ^= u >> 33
	return u
}

// GetHasher returns a default hasher function for different Key types.
func GetHasher[Key any]() HashFn[Key] {
	var key Key
	kind := reflect.ValueOf(&key).Elem().Type().Kind()

	var (
		hashByte  = HashInt8
		hashWord  = HashUint16
		hashDword = HashUint32
		hashQword = HashUint64
		hashF32   = HashFloat32
		hashF64   = HashFloat64
		hashStr   = HashString
	)

	switch kind {
	case reflect.Int, reflect.Uint, reflect.Uintptr:
		switch unsafe.Sizeof(key) {
		case 2:
			return *(*func(Key) uint64)(unsafe.Pointer(&hashWord))
		case 4:
			return *(*func(Key) uint64)(unsafe.Pointer(&hashDword))
		case 8:
			return *(*func(Key) uint64)(unsafe.Pointer(&hashQword))

		default:
			panic(fmt.Errorf("unsupported integer byte size"))
		}

	case reflect.Int8, reflect.Uint8:
		return *(*func(Key) uint64)(unsafe.Pointer(&hashByte))
	case reflect.Int16, reflect.Uint16:
		return *(*func(Key) uint64)(unsafe.Pointer(&hashWord))
	case reflect.Int32, reflect.Uint32:
		return *(*func(Key) uint64)(unsafe.Pointer(&hashDword))
	case reflect.Int64, reflect.Uint64:
		return *(*func(Key) uint64)(unsafe.Pointer(&hashQword))
	case reflect.Float32:
		return *(*func(Key) uint64)(unsafe.Pointer(&hashF32))
	case reflect.Float64:
		return *(*func(Key) uint64)(unsafe.Pointer(&hashF64))
	case reflect.String:
		return *(*func(Key) uint64)(unsafe.Pointer(&hashStr))

	default:
		panic(fmt.Errorf("unsupported key type %T of kind %v", key, kind))
	}
}
