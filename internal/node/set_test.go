package node_test

import (
	"testing"

	"github.com/srujan964/go-dist-sys-challenge/internal/node"
)

func TestSet_Intersect(t *testing.T) {
	s := node.NewSet()
	for _, it := range []int64{0, 1, 2, 3, 4, 5} {
		s.Add(it)
	}

	other := node.NewSet()
	for _, it := range []int64{3, 4, 5, 6, 7} {
		other.Add(it)
	}

	result := s.Intersect(other)

	if !result.Contains(3) && !result.Contains(4) && !result.Contains(5) {
		t.Errorf("%v does not contain the common elements of %v and %v", result, s, other)
	}

	if result.Contains(0) || result.Contains(1) ||
		result.Contains(2) || result.Contains(6) || result.Contains(7) {
		t.Errorf("%v contains elements found only in one set", result)
	}
}

func TestSet_Union(t *testing.T) {
	s := node.NewSet()
	for _, it := range []int64{0, 1, 2, 3, 4, 5} {
		s.Add(it)
	}

	other := node.NewSet()
	for _, it := range []int64{3, 4, 5, 6, 7} {
		other.Add(it)
	}

	result := s.Intersect(other)

	if !result.Contains(0) && !result.Contains(1) && !result.Contains(2) &&
		!result.Contains(3) && !result.Contains(4) && !result.Contains(5) &&
		!result.Contains(6) && !result.Contains(7) {
		t.Errorf("%v does not contain the combined elements of %v and %v", result, s, other)
	}
}

func TestSet_Difference(t *testing.T) {
	s := node.NewSet()
	for _, it := range []int64{0, 1, 2, 3, 4, 5} {
		s.Add(it)
	}

	other := node.NewSet()
	for _, it := range []int64{3, 4, 5, 6, 7} {
		other.Add(it)
	}

	result := s.Difference(other)

	if !result.Contains(0) || !result.Contains(1) || !result.Contains(2) ||
		result.Contains(3) || result.Contains(4) || result.Contains(5) {
		t.Errorf("%v does not contain the set difference of %v and %v", result, s, other)
	}
}
