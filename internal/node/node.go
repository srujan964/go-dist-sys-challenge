package node

import (
	"log"
	"math/rand"
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
	kv            *maelstrom.KV
	seq           uint64

	mu     sync.RWMutex
	latest int
	done   chan bool
}

func CreateNode(n *maelstrom.Node, kv *maelstrom.KV) *Node {
	return &Node{
		id:            n.ID(),
		maelstromNode: n,
		kv:            kv,
		messages:      NewSet(),
		latest:        0,
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

func RandomNeighbours(neighbours []string) []string {
	rand.Shuffle(len(neighbours), func(i, j int) {
		neighbours[i], neighbours[j] = neighbours[j], neighbours[i]
	})
	m := len(neighbours)/2 + 1
	return neighbours[:m]
}
