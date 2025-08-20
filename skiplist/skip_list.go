// Skip List Implementation in Go.
// This implementation is a generic skip list that supports concurrent access.
package skiplist

import (
	"fmt"
	"math"
	"math/rand/v2"
	"sync"
)

const (
	LAYER_PROMOTION_PROB = 0.25
	MAX_LAYER            = 32
)

type SkipList[T any] interface {
	Get(key int) (T, bool)
	Insert(key int, val T)
	Delete(key int) bool
	Size() int
}

var _ SkipList[int] = (*skipList[int])(nil)

type skipList[T any] struct {
	head  *skipListNode[T] // Head node of the skip list
	level int              // Current highest level of the skip list
	size  int              // Number of elements in the skip list
	rng   func() int       // Function to generate random levels for new nodes
	pool  sync.Pool        // Pool for reusing nodes to reduce memory allocation overhead
	mu    sync.RWMutex     // Mutex for thread-safe operations
}

type skipListNode[T any] struct {
	forward [MAX_LAYER + 1]*skipListNode[T] // Pointers to the next nodes at each level
	key     int                             //	Key of the node. Used for sorting.
	level   int                             // The number of levels this node has
	value   T
}

func CreateSkipList[T any]() *skipList[T] {
	s := &skipList[T]{
		rng:  defaultRngLevelGen(LAYER_PROMOTION_PROB, MAX_LAYER),
		pool: newPool[T](),
	}
	s.head = s.createNode(MAX_LAYER, 0, *new(T)) // Create a head node with maximum level
	return s
}

// Get retrieves the value associated with the given key from the skip list.
func (s *skipList[T]) Get(key int) (T, bool) {
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
func (s *skipList[T]) Insert(key int, val T) {
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

	lvl := s.rng()
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

func (s *skipList[T]) Delete(key int) bool {
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

func (s *skipList[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.size
}

func (s *skipList[T]) Print() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := s.level; i >= 0; i-- {
		x := s.head.forward[i]
		for x != nil {
			print(x.key, " ")
			x = x.forward[i]
		}
		println()
	}
}

func (n *skipListNode[T]) String() string {
	return "{Key: " + fmt.Sprint(n.key) + ", Value: " + fmt.Sprint(n.value) + ", Level: " + fmt.Sprint(n.level) + "}"
}

func (s *skipList[T]) createNode(level, key int, value T) *skipListNode[T] {
	n := s.pool.Get().(*skipListNode[T])
	n.level = level
	n.key = key
	n.value = value
	return n
}

func (s *skipList[T]) freeNode(node *skipListNode[T]) {
	var zero T
	node.value = zero
	for i := range node.forward {
		node.forward[i] = nil
	}
	node.level = 0
	node.key = 0
	s.pool.Put(node)
}

func defaultRngLevelGen(p float64, m int) func() int {
	return func() int {
		return min(m, geometric(p))
	}
}

// geometric distribution sampler.
func geometric(p float64) int {
	if p <= 0 || p >= 1 {
		panic("nice try, p must be in (0,1)")
	}
	u := rand.Float64() // uniform(0,1)
	return int(math.Ceil(math.Log(1-u) / math.Log(p)))
}

func newPool[T any]() sync.Pool {
	return sync.Pool{
		New: func() any {
			return &skipListNode[T]{}
		},
	}
}
