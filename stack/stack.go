package stack

import "errors"

type Stack[T any] []T

func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Push(i T) {
	*s = append(*s, i)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(*s) == 0 {
		var empty T
		return empty, errors.New("empty stack")
	}
	i := len(*s) - 1
	elem := (*s)[i]
	*s = (*s)[:i]
	return elem, nil
}

func (s *Stack[T]) Values() []T {
	return *s
}

func (s *Stack[T]) Clear() {
	*s = []T{}
}
