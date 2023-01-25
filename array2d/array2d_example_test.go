package array2d_test

import (
	"fmt"
	"strings"

	"github.com/zyedidia/generic/array2d"
)

type Sudoku struct {
	arr array2d.Array2D[byte]
}

func (s Sudoku) PrintBoard() {
	var sb strings.Builder
	for y := 0; y < s.arr.Height(); y++ {
		if y%3 == 0 {
			sb.WriteString("+-------+-------+-------+\n")
		}
		for x := 0; x < s.arr.Width(); x++ {
			if x%3 == 0 {
				sb.WriteString("| ")
			}
			val := s.arr.Get(x, y)
			if val == 0 {
				sb.WriteByte(' ')
			} else {
				fmt.Fprint(&sb, val)
			}
			sb.WriteByte(' ')
		}
		sb.WriteString("|\n")
	}
	sb.WriteString("+-------+-------+-------+\n")
	fmt.Print(sb.String())
}

func ExampleArray2D() {
	s := Sudoku{
		arr: array2d.OfJagged(9, 9, [][]byte{
			{5, 3, 0, 0, 7, 0, 0, 0, 0},
			{6, 0, 0, 1, 9, 5, 0, 0, 0},
			{0, 9, 8, 0, 0, 0, 0, 6, 0},
			{8, 0, 0, 0, 6, 0, 0, 0, 3},
			{4, 0, 0, 8, 0, 3, 0, 0, 1},
			{7, 0, 0, 0, 2, 0, 0, 0, 6},
			{0, 6, 0, 0, 0, 0, 2, 8, 0},
			{0, 0, 0, 4, 1, 9, 0, 0, 5},
			{0, 0, 0, 0, 8, 0, 0, 7, 9},
		}),
	}

	s.arr.Set(2, 5, 3)

	s.PrintBoard()

	// Output:
	// +-------+-------+-------+
	// | 5 3   |   7   |       |
	// | 6     | 1 9 5 |       |
	// |   9 8 |       |   6   |
	// +-------+-------+-------+
	// | 8     |   6   |     3 |
	// | 4     | 8   3 |     1 |
	// | 7   3 |   2   |     6 |
	// +-------+-------+-------+
	// |   6   |       | 2 8   |
	// |       | 4 1 9 |     5 |
	// |       |   8   |   7 9 |
	// +-------+-------+-------+
}
