package stack

import "testing"

func TestPopReleasesUnusedSlot(t *testing.T) {
	for _, try := range []bool{false, true} {
		name := "Pop"
		if try {
			name = "TryPop"
		}
		t.Run(name, func(t *testing.T) {
			s := New[*int]()
			first, second := 1, 2
			s.Push(&first)
			s.Push(&second)
			storage := s.entries[:cap(s.entries)]
			var got *int
			if try {
				var ok bool
				got, ok = s.TryPop()
				if !ok {
					t.Fatal("missing value")
				}
			} else {
				got = s.Pop()
			}
			if got != &second || s.Size() != 1 || s.Peek() != &first {
				t.Fatal("LIFO contract changed")
			}
			if storage[1] != nil {
				t.Fatal("popped pointer remains reachable from unused backing slot")
			}
			if try {
				s.TryPop()
			} else {
				s.Pop()
			}
			if storage[0] != nil {
				t.Fatal("last popped pointer remains reachable")
			}
			s.Push(&second)
			if s.Pop() != &second {
				t.Fatal("reuse")
			}
		})
	}
}
