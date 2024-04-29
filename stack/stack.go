package stack

import "errors"

// Stack represents a generic stack data structure.
type Stack[T any] []T

// New creates a new empty stack.
func New[T any]() *Stack[T] {
	return &Stack[T]{}
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(value T) {
	*s = append(*s, value)
}

// Pop removes and returns the top element from the stack.
// It returns an error if the stack is empty.
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

// Peek returns the top element of the stack without removing it.
// It returns an error if the stack is empty.
func (s *Stack[T]) Peek() (T, error) {
	if len(*s) == 0 {
		var empty T
		return empty, errors.New("empty stack")
	}
	return (*s)[len(*s)-1], nil
}

// Values returns the values of the stack.
func (s *Stack[T]) Values() []T {
	return *s
}

// Clear removes all elements from the stack.
func (s *Stack[T]) Clear() {
	*s = []T{}
}
