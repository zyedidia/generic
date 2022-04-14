package generic_test

import (
	"fmt"
	"math"
	"time"

	"github.com/zyedidia/generic"
)

func ExampleMax() {
	fmt.Println(generic.Max(7, 3))
	fmt.Println(generic.Max(2*time.Second, 3*time.Second).Milliseconds())
	// Output:
	// 7
	// 3000
}

func ExampleMin() {
	fmt.Println(generic.Min(7, 3))
	fmt.Println(generic.Min(2*time.Second, 3*time.Second).Milliseconds())
	// Output:
	// 3
	// 2000
}

func ExampleClamp() {
	fmt.Println(generic.Clamp(500, 400, 600))
	fmt.Println(generic.Clamp(200, 400, 600))
	fmt.Println(generic.Clamp(800, 400, 600))

	fmt.Println(generic.Clamp(5*time.Second, 4*time.Second, 6*time.Second).Milliseconds())
	fmt.Println(generic.Clamp(2*time.Second, 4*time.Second, 6*time.Second).Milliseconds())
	fmt.Println(generic.Clamp(8*time.Second, 4*time.Second, 6*time.Second).Milliseconds())

	fmt.Println(generic.Clamp(1.5, 1.4, 1.8))
	fmt.Println(generic.Clamp(1.5, 1.8, 1.8))
	fmt.Println(generic.Clamp(1.5, 2.1, 1.9))

	// Output:
	// 500
	// 400
	// 600
	// 5000
	// 4000
	// 6000
	// 1.5
	// 1.8
	// 2.1
}

func lessMagnitude(a, b float64) bool {
	return math.Abs(a) < math.Abs(b)
}

func ExampleMaxFunc() {
	fmt.Println(generic.MaxFunc(2.5, -3.1, lessMagnitude))
	// Output:
	// -3.1
}

func ExampleMinFunc() {
	fmt.Println(generic.MinFunc(2.5, -3.1, lessMagnitude))
	// Output:
	// 2.5
}

func ExampleClampFunc() {
	fmt.Println(generic.ClampFunc(1.5, 1.4, 1.8, lessMagnitude))
	fmt.Println(generic.ClampFunc(1.5, 1.8, 1.8, lessMagnitude))
	fmt.Println(generic.ClampFunc(1.5, 2.1, 1.9, lessMagnitude))
	fmt.Println(generic.ClampFunc(-1.5, -1.4, -1.8, lessMagnitude))
	fmt.Println(generic.ClampFunc(-1.5, -1.8, -1.8, lessMagnitude))
	fmt.Println(generic.ClampFunc(-1.5, -2.1, -1.9, lessMagnitude))
	fmt.Println(generic.ClampFunc(1.5, -1.5, -1.5, lessMagnitude))

	// Output:
	// 1.5
	// 1.8
	// 2.1
	// -1.5
	// -1.8
	// -2.1
	// 1.5
}
