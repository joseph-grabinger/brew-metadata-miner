package stack

import (
	"testing"
)

func TestPush(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	values := s.Values()
	expected := []int{1, 2, 3}

	if len(values) != len(expected) {
		t.Errorf("expected length %d, got %d", len(expected), len(values))
	}

	for i, v := range values {
		if v != expected[i] {
			t.Errorf("expected value %d at index %d, got %d", expected[i], i, v)
		}
	}
}

func TestPop(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	value, err := s.Pop()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if value != 3 {
		t.Errorf("expected value 3, got %d", value)
	}

	value, err = s.Pop()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if value != 2 {
		t.Errorf("expected value 2, got %d", value)
	}

	value, err = s.Pop()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if value != 1 {
		t.Errorf("expected value 1, got %d", value)
	}

	_, err = s.Pop()
	if err == nil {
		t.Error("expected an error, got nil")
	} else if err.Error() != "empty stack" {
		t.Errorf("expected error 'empty stack', got %v", err)
	}
}

func TestPeek(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	value, err := s.Peek()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if value != 3 {
		t.Errorf("expected value 3, got %d", value)
	}

	value, err = s.Peek()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if value != 3 {
		t.Errorf("expected value 3, got %d", value)
	}

	s.Pop()
	s.Pop()
	s.Pop()

	_, err = s.Peek()
	if err == nil {
		t.Error("expected an error, got nil")
	} else if err.Error() != "empty stack" {
		t.Errorf("expected error 'empty stack', got %v", err)
	}
}

func TestClear(t *testing.T) {
	s := New[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	s.Clear()

	values := s.Values()
	if len(values) != 0 {
		t.Errorf("expected length 0, got %d", len(values))
	}
}
