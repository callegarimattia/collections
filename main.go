package main

import (
	"fmt"

	"example.com/skiplist"
)

func main() {
	s := skiplist.CreateSkipList[string]()
	s.Insert(1, "one")
	s.Insert(3, "three")
	s.Insert(25, "twenty-five")
	s.Insert(15, "fifteen")
	s.Insert(5, "five")
	s.Print()
	v, _ := s.Get(3)
	fmt.Printf("Searching for 15: %v\n", v)
	v, _ = s.Get(15)
	fmt.Printf("Searching for 100: %v\n", v)
	s.Delete(3)
	s.Delete(25)
	fmt.Println("After deletions:")
	s.Print()
}
