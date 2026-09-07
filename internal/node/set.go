package node

type Set struct {
	elements map[int64]struct{}
}

func (s *Set) Size() int {
	return len(s.elements)
}

// Returns a Set containing the elements present in both s and other.
func (s *Set) Intersect(other *Set) *Set {
	new := &Set{
		elements: make(map[int64]struct{}, 0),
	}
	var (
		bigger  = s
		smaller = other
	)
	if bigger.Size() < smaller.Size() {
		bigger, smaller = other, s
	}

	for it := range smaller.elements {
		if bigger.Contains(it) {
			new.Add(it)
		}
	}
	return new
}

// Returns a Set containing the elements in the set but not in the other.
func (s *Set) Difference(other *Set) *Set {
	size := max(0, s.Size()-other.Size())
	new := &Set{
		elements: make(map[int64]struct{}, size),
	}

	for it := range s.elements {
		if !other.Contains(it) {
			new.Add(it)
		}
	}
	return new
}

// Returns a Set containing the elements of both s and other combined.
func (s *Set) Union(other *Set) *Set {
	size := max(s.Size(), other.Size())
	new := &Set{
		elements: make(map[int64]struct{}, size),
	}
	new.AddAll(s)
	new.AddAll(other)
	return new
}

func (s *Set) Contains(element int64) bool {
	_, exists := s.elements[element]
	return exists
}

func (s *Set) Add(element int64) error {
	if !s.Contains(element) {
		s.elements[element] = struct{}{}
	}
	return nil
}

func (s *Set) Elements() []int64 {
	ret := []int64{}
	if len(s.elements) == 0 {
		return ret
	}

	for e := range s.elements {
		ret = append(ret, e)
	}
	return ret
}

func (s *Set) AddAll(other *Set) bool {
	updated := false
	for it := range other.elements {
		if s.Contains(it) {
			continue
		}
		s.Add(it)
		updated = true
	}
	return updated
}

func SetFrom(elements []int64) *Set {
	new := &Set{
		elements: make(map[int64]struct{}, len(elements)),
	}
	for _, it := range elements {
		new.Add(it)
	}
	return new
}

func NewSet() *Set {
	return &Set{
		elements: map[int64]struct{}{},
	}
}
