package stack

// Stack is a generic stack.
type Stack[T any] interface {
	Pop() (T, bool)
	Push(T)
	Peek() (T, bool)
	Len() int
	Cap() int
}

// stack is the default, private implementation of a stack.
type stack[T any] []T

// Len returns the number of elements in the stack.
func (s *stack[T]) Len() int {
	return len(*s)
}

// Cap returns the capacity of the stack.
func (s *stack[T]) Cap() int {
	return cap(*s)
}

// Pop removes and returns the top element of the stack.
// If the stack is empty, it returns the zero value of type T and false.
// Otherwise, it returns the top element and true.
func (s *stack[T]) Pop() (T, bool) {
	if len(*s) == 0 {
		return *new(T), false
	}
	t := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return t, true
}

// Push adds an element to the top of the stack.
func (s *stack[T]) Push(t T) {
	*s = append(*s, t)
}

// Peek returns the top element of the stack without removing it.
func (s *stack[T]) Peek() (T, bool) {
	if len(*s) == 0 {
		return *new(T), false
	}
	return (*s)[len(*s)-1], true
}

// New creates a new stack.
// If one argument is provided, it creates a stack with the specified capacity.
// If no arguments are provided, it creates a stack with the default capacity.
// If more than one argument is provided, it panics.
// If the provided capacity is negative, it creates a stack with the default capacity.
func New[T any](capacity ...int) Stack[T] {
	switch {
	case len(capacity) == 1:
		if capacity[0] < 0 {
			capacity[0] = 0
		}
		s := make(stack[T], 0, capacity[0])
		return &s
	case len(capacity) > 1:
		panic("stack.New() expects at most one argument")
	default:
		s := make(stack[T], 0)
		return &s
	}
}
