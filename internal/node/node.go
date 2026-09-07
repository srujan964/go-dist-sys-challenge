package node

import (
	"fmt"
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
	wg   sync.WaitGroup
	done chan bool
}

func Key(msgType string, msgId int64, dest string) string {
	return fmt.Sprintf("%s-%d-%s", msgType, msgId, dest)
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

func (n *Node) setNeighbours(nodes []string) {
	n.neighbours = nodes
}

func (n *Node) Shutdown() {
	n.done <- true
}
