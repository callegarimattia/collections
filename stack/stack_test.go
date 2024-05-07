package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	t.Run("A new stack should be empty", func(t *testing.T) {
		s := New[int]()
		assert.Empty(t, s, "A new stack should be empty")
	})
	t.Run("A stack should have a length of 0 when empty", func(t *testing.T) {
		s := New[int]()
		assert.Zero(t, s.Len(), "An empty stack should have a length of 0")
	})

	t.Run(
		"A stack created using New() with a negative capacity should have a length of 0 and capacity of 0",
		func(t *testing.T) {
			s := New[int](-1)
			assert.Zero(t, s.Len(), "Should have a length of 0")
			assert.Zero(t, s.Cap(), "Should have capacity of 0")
		},
	)

	t.Run(
		"A stack created using New() with a positive capacity should have a length of 0 and the specified capacity",
		func(t *testing.T) {
			s := New[int](10)
			assert.Zero(t, s.Len(), "Should have a length of 0")
			assert.Equal(t, 10, s.Cap(), "Should have capacity of 10")
		},
	)

	t.Run("A stack created using New() with more the one argument should panic",
		func(t *testing.T) {
			assert.Panics(
				t,
				func() { New[int](1, 2) },
				"New() should panic when more than one argument is provided",
			)
		},
	)

	t.Run(
		"A stack should return the top element of the stack without removing it when using Peek()",
		func(t *testing.T) {
			s := New[int]()
			s.Push(1)
			s.Push(2)
			s.Push(3)
			v, ok := s.Peek()
			assert.True(t, ok, "The top element should exist")
			assert.Equal(t, 3, v, "The top element should be 3")
			assert.Equal(t, 3, s.Len(), "The stack should have 3 elements")
		},
	)

	t.Run(
		"A stack should remove and return the top element of the stack when using Pop()",
		func(t *testing.T) {
			s := New[int]()
			s.Push(1)
			s.Push(2)
			s.Push(3)
			v, ok := s.Pop()
			assert.True(t, ok, "The top element should exist")
			assert.Equal(t, 3, v, "The top element should be 3")
			assert.Equal(t, 2, s.Len(), "The stack should have 2 elements")
		},
	)

	t.Run(
		"A stack should return false when trying to Pop() an empty stack",
		func(t *testing.T) {
			s := New[int]()
			v, ok := s.Pop()
			assert.False(t, ok, "The top element should not exist")
			assert.Zero(t, v, "The top element should be zero")
			assert.Zero(t, s.Len(), "The stack should have 0 elements")
		},
	)

	t.Run(
		"A stack should return false when trying to Peek() an empty stack",
		func(t *testing.T) {
			s := New[int]()
			v, ok := s.Peek()
			assert.False(t, ok, "The top element should not exist")
			assert.Zero(t, v, "The top element should be zero")
			assert.Zero(t, s.Len(), "The stack should have 0 elements")
		},
	)

	t.Run(
		"A stack should store new elements in the order they were added",
		func(t *testing.T) {
			s := New[int]()
			s.Push(1)
			s.Push(2)
			s.Push(3)
			v, ok := s.Pop()
			assert.True(t, ok, "The top element should exist")
			assert.Equal(t, 3, v, "The top element should be 3")
			v, ok = s.Pop()
			assert.True(t, ok, "The top element should exist")
			assert.Equal(t, 2, v, "The top element should be 2")
			v, ok = s.Pop()
			assert.True(t, ok, "The top element should exist")
			assert.Equal(t, 1, v, "The top element should be 1")
			assert.Zero(t, s.Len(), "The stack should have 0 elements")
		},
	)
}
