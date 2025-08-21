# collections

Collections in Go

1. SkipList
Skiplist is a probabilistic data structure that allows for fast search, insertion, and deletion operations. It is built on top of a linked list and uses multiple levels to achieve logarithmic time complexity for these operations.
The implementation leverages the `math/rand/v2` package to generate random levels for the skip list nodes, ensuring that the average time complexity remains logarithmic.
It's concurrent safe, meaning it can be used in multi-threaded applications without additional synchronization.
It's also designed to be memory efficient, using a compact representation of nodes and levels.
Additionally, it leverages a local node pool to avoid frequent memory allocations, which can be a performance bottleneck in high-throughput applications.

2. Stack
A simple, generic, stack implementation in Go. It provides basic stack operations such as push, pop, and peek, allowing for easy management of a collection of elements in a last-in-first-out (LIFO) manner.

[![Go Reference](https://pkg.go.dev/badge/github.com/callegarimattia/collections.svg)](https://pkg.go.dev/github.com/callegarimattia/collections)
