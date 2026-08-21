package node

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type RPCToplogyMsg struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

func (n *Node) ApplyTopology(msg maelstrom.Message) error {
	n.setID()
	var body RPCToplogyMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("Unable to decode Topology message body - %w", err)
	}
	neighbours, exists := body.Topology[n.id]
	if !exists {
		return fmt.Errorf("Node %s does not exist in topology %v", n.id, body.Topology)
	}
	n.setNeighbours(neighbours)
	return n.maelstromNode.Reply(msg, RPCReply{
		Type: "topology_ok",
	})
}
