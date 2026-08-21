package node

import (
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Node struct {
	id            string
	messages      *Set
	neighbours    []string
	maelstromNode *maelstrom.Node
	seq           uint64
}

type Set struct {
	mu       *sync.RWMutex
	elements map[int64]struct{}
}

func (s *Set) Contains(element int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.elements[element]
	return exists
}

func (s *Set) Add(element int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.elements[element]
	if !exists {
		s.elements[element] = struct{}{}
	}
	return nil
}

func (s *Set) ReadAll() []int64 {
	ret := make([]int64, 0, len(s.elements))
	s.mu.RLock()
	defer s.mu.RUnlock()
	for e := range s.elements {
		ret = append(ret, e)
	}
	return ret
}

func NewSet() *Set {
	return &Set{
		mu:       &sync.RWMutex{},
		elements: map[int64]struct{}{},
	}
}

func CreateNode(n *maelstrom.Node) *Node {
	return &Node{
		id:            n.ID(),
		maelstromNode: n,
		messages:      NewSet(),
	}
}

func (n *Node) setID() {
	n.id = n.maelstromNode.ID()
}

func (n *Node) setNeighbours(nodes []string) {
	n.neighbours = nodes
}
