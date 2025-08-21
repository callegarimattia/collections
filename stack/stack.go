package stack

// Stack is the default, private implementation of a Stack.
type Stack[T any] []T

// Len returns the number of elements in the stack.
func (s *Stack[T]) Len() int {
	return len(*s)
}

// Cap returns the capacity of the stack.
func (s *Stack[T]) Cap() int {
	return cap(*s)
}

// Pop removes and returns the top element of the stack.
// If the stack is empty, it returns the zero value of type T and false.
// Otherwise, it returns the top element and true.
func (s *Stack[T]) Pop() (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	t := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return t, true
}

// Push adds an element to the top of the stack.
func (s *Stack[T]) Push(t T) {
	*s = append(*s, t)
}

// Peek returns the top element of the stack without removing it.
func (s *Stack[T]) Peek() (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	return (*s)[len(*s)-1], true
}

// New creates a new stack with an optional initial capacity.
func New[T any](cap ...int) *Stack[T] {
	if len(cap) > 1 {
		panic("stack.New: only one capacity argument is allowed")
	}
	if len(cap) == 1 && cap[0] < 0 {
		panic("stack.New: capacity cannot be negative")
	}
	if len(cap) == 0 {
		cap = append(cap, 10) // Default capacity
	}
	s := make(Stack[T], 0, cap[0])
	return &s
}
