package stack

type Stack[T any] struct {
	entries []T
}

func New[T any]() *Stack[T] {
	return &Stack[T]{
		entries: nil,
	}
}

func (s *Stack[T]) Push(value T) {
	s.entries = append(s.entries, value)
}

func (s *Stack[T]) Pop() T {
	v := s.entries[len(s.entries)-1]
	s.entries = s.entries[:len(s.entries)-1]
	return v
}

func (s *Stack[T]) Peek() T {
	return s.entries[len(s.entries)-1]
}

func (s *Stack[T]) Size() int {
	return len(s.entries)
}
