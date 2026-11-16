package node

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Node struct {
	id            string
	messages      *Set
	neighbours    []string
	maelstromNode *maelstrom.Node
	seq           uint64

	mu   sync.RWMutex
	done chan bool
}

func CreateNode(n *maelstrom.Node) *Node {
	return &Node{
		id:            n.ID(),
		maelstromNode: n,
		messages:      NewSet(),
		done:          make(chan bool),
	}
}

func (n *Node) setID() {
	n.id = n.maelstromNode.ID()
}

func (n *Node) InitTopology(k int) {
	sorted := append([]string(nil), n.maelstromNode.NodeIDs()...)
	sort.Strings(sorted)

	idx := -1
	for i, id := range sorted {
		if id == n.id {
			idx = i
			break
		}
	}

	if idx == -1 || k < 1 {
		return
	}
	var neighbours []string
	if idx > 0 {
		parent := (idx - 1) / k
		neighbours = append(neighbours, sorted[parent])
	}
	first := k*idx + 1
	last := k*idx + k
	for c := first; c <= last && c < len(sorted); c++ {
		neighbours = append(neighbours, sorted[c])
	}

	n.neighbours = neighbours
}

func (n *Node) Shutdown() {
	n.done <- true
}

func NumericalID(nodeId string) (int, error) {
	parsedNodeId, err := strconv.Atoi(strings.TrimPrefix(nodeId, "n"))
	if err != nil {
		log.Printf("Cannot convert node id to int - %v", err)
		return 0, err
	}
	return int(parsedNodeId), nil
}
