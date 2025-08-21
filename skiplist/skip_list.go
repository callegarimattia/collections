// Skip List Implementation in Go.
// This implementation is a generic skip list that supports concurrent access.
// It leverages a sync.Pool for node reuse to minimize memory allocation overhead.
// It uses a novel random level generator based on a geometric distribution, instead
// of the traditional iterational coin flip method.
package skiplist

import (
	"math"
	"math/rand/v2"
	"sync"
)

// The probability of promoting a node to the next level.
// It can be adjusted to control the average height of the skip list.
var LAYER_PROMOTION_PROB = 0.25

// The maximum number of layers in the skip list.
const MAX_LAYER = 32

// The SkipList struct represents a skip list.
// It contains a head node, the current highest level, the size of the list,
// a sync.Pool for node reuse, and a rwmutex.
type SkipList[T any] struct {
	head  *skipListNode[T] // Head node of the skip list
	level int              // Current highest level of the skip list
	size  int              // Number of elements in the skip list
	pool  sync.Pool        // Pool for reusing nodes to reduce memory allocation overhead
	mu    sync.RWMutex     // Mutex for thread-safe operations
}

type skipListNode[T any] struct {
	forward [MAX_LAYER + 1]*skipListNode[T] // Pointers to the next nodes at each level
	key     int                             //	Key of the node. Used for sorting.
	level   int                             // The number of levels this node has
	value   T                               // Value associated with the key
}

// CreateSkipList initializes a new skip list with the default parameters.
// The skip list is generic and can hold any type of value.
// The random level generator is set to a geometric distribution with a promotion probability
// decided by the package level variable `LAYER_PROMOTION_PROB`.
func CreateSkipList[T any]() *SkipList[T] {
	s := &SkipList[T]{pool: newPool[T]()}
	s.head = s.createNode(MAX_LAYER, 0, *new(T)) // Create a head node with maximum level
	return s
}

// Get retrieves the value associated with the given key from the skip list.
// Time complexity is O(log n) on average.
func (s *SkipList[T]) Get(key int) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	x := s.head
	for i := s.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
	}

	x = x.forward[0]

	if x != nil && x.key == key {
		return x.value, true
	} else {
		var zero T
		return zero, false
	}
}

// Insert adds a new node with the given key and value to the skip list.
// Time complexity is O(log n) on average.
func (s *SkipList[T]) Insert(key int, val T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	update := make([]*skipListNode[T], MAX_LAYER+1)
	x := s.head
	for i := s.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]
	if x != nil && x.key == key {
		x.value = val
		return
	}

	lvl := defaultRngLevelGen(LAYER_PROMOTION_PROB, MAX_LAYER)
	if lvl > s.level {
		for i := s.level + 1; i <= lvl; i++ {
			update[i] = s.head
		}
		s.level = lvl
	}

	x = s.createNode(lvl, key, val)
	for i := 0; i <= lvl; i++ {
		x.forward[i] = update[i].forward[i]
		update[i].forward[i] = x
	}

	s.size++
}

// Delete removes the node with the given key from the skip list.
// Time complexity is O(log n) on average.
func (s *SkipList[T]) Delete(key int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	update := [MAX_LAYER + 1]*skipListNode[T]{}
	x := s.head

	for i := s.level; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < key {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[0]

	if x == nil || x.key != key {
		return false
	}

	for i := 0; i <= s.level; i++ {
		if update[i].forward[i] != x {
			break
		}
		update[i].forward[i] = x.forward[i]
	}

	for s.level > 0 && s.head.forward[s.level] == nil {
		s.level--
	}

	s.size--
	s.freeNode(x)

	return true
}

// Size returns the number of elements in the skip list.
func (s *SkipList[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.size
}

func (s *SkipList[T]) createNode(level, key int, value T) *skipListNode[T] {
	n := s.pool.Get().(*skipListNode[T])
	n.level = level
	n.key = key
	n.value = value
	return n
}

func (s *SkipList[T]) freeNode(node *skipListNode[T]) {
	var zero T
	node.value = zero
	for i := range node.forward {
		node.forward[i] = nil
	}
	node.level = 0
	node.key = 0
	s.pool.Put(node)
}

func defaultRngLevelGen(p float64, m int) int {
	return min(m, geometric(p))
}

// geometric distribution sampler.
// panics if p is not in (0,1).
func geometric(p float64) int {
	if p <= 0 || p >= 1 {
		panic("nice try, p must be in (0,1)")
	}
	u := rand.Float64() // uniform(0,1)

	return int(math.Ceil(math.Log(1-u) / math.Log(p))) // math is hard.
}

// Generates a new sync.Pool for skip list nodes of type T.
func newPool[T any]() sync.Pool {
	return sync.Pool{
		New: func() any {
			return &skipListNode[T]{}
		},
	}
}
