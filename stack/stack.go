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

func (s *Stack[T]) Pop() (t T) {
	if len(s.entries) == 0 {
		return t
	}
	v := s.entries[len(s.entries)-1]
	s.entries = s.entries[:len(s.entries)-1]
	return v
}

func (s *Stack[T]) Peek() (t T) {
	if len(s.entries) == 0 {
		return t
	}
	return s.entries[len(s.entries)-1]
}

func (s *Stack[T]) Size() int {
	return len(s.entries)
}
