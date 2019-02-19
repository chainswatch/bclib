package net

import (
	"testing"
)

func TestQueue(t *testing.T) {
	q := NewQueue(3)

	s := "07df89107e2f4a03c3e055d745169ef1f8a35053202e26d40513a71fcff81372"
	var hash [32]byte
	copy(hash[:], s)
	q.Push(hash, nil)
	if len(q.invs) != 1 {
		t.Fatal("Length should be 1, is ", len(q.invs))
	}
	q.Push(hash, nil)
	if len(q.invs) != 1 {
		t.Fatal("Length should be 1, is ", len(q.invs))
	}
	hash[0] = 5
	q.Push(hash, nil)
	if len(q.invs) != 2 {
		t.Fatal("Length should still be 2")
	}
	hash[0] = 6
	q.Push(hash, nil)
	if len(q.invs) != 3 {
		t.Fatal("Length should still be 3")
	}

	hash[0] = 7
	q.Push(hash, nil)
	if len(q.invs) != 3 {
		t.Fatal("Length should still be 3")
	}
	if err := q.Update(hash, nil); err != nil {
		t.Fatal(err)
	}
}
