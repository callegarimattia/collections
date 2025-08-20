package skiplist_test

import (
	"math/rand/v2"
	"testing"

	"example.com/skiplist"
)

func BenchmarkInsert(b *testing.B) {
	randSize := 1_000_000
	keys := rand.Perm(randSize)
	s := skiplist.CreateSkipList[*int]()

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		s.Insert(keys[i%randSize], nil)
	}
}

func BenchmarkInsertPrefilled100k(b *testing.B) {
	randSize := 1_000_000
	keys := rand.Perm(randSize)
	s := skiplist.CreateSkipList[*int]()

	for i := range randSize {
		s.Insert(keys[i], nil)
	}

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		s.Insert(keys[i%randSize], nil)
	}
}

func BenchmarkDelete(b *testing.B) {
	const prefill = 1_000
	const maxKey = 200_000

	keys := rand.Perm(maxKey)
	s := skiplist.CreateSkipList[*int]()

	for i := range prefill {
		s.Insert(keys[i], nil)
	}

	b.ResetTimer()

	for b.Loop() {
		delKey := keys[rand.IntN(maxKey)]
		s.Delete(delKey)
	}
}

func BenchmarkDynamicInsertDelete(b *testing.B) {
	const prefill = 100_000
	const randSize = 10_000_000
	keys := rand.Perm(randSize)

	s := skiplist.CreateSkipList[*int]()
	for i := range prefill {
		s.Insert(keys[i], nil)
	}

	for i := 0; b.Loop(); i++ {
		s.Delete(keys[i%randSize])
		s.Insert(keys[i%randSize], nil)
	}
}

func BenchmarkGet(b *testing.B) {
	const prefill = 100_000
	keys := rand.Perm(prefill)
	s := skiplist.CreateSkipList[*int]()
	
	for i := range prefill {
		s.Insert(keys[i], nil)
	}
	
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		s.Get(keys[i%prefill])
	}
}

func BenchmarkMixedWorkload(b *testing.B) {
	const prefill = 50_000
	const randSize = 1_000_000
	keys := rand.Perm(randSize)
	s := skiplist.CreateSkipList[*int]()
	
	for i := range prefill {
		s.Insert(keys[i], nil)
	}
	
	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		keyIndex := i % randSize
		switch i % 10 {
		case 0, 1, 2, 3, 4, 5: // 60% reads
			s.Get(keys[keyIndex])
		case 6, 7, 8: // 30% inserts
			s.Insert(keys[keyIndex], nil)
		case 9: // 10% deletes
			s.Delete(keys[keyIndex])
		}
	}
}
