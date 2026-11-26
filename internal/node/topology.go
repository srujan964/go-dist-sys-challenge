package node

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type RPCToplogyMsg struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

func (n *Node) ApplyTopology(msg maelstrom.Message) error {
	n.setID()

	n.InitTopology(5)

	return n.maelstromNode.Reply(msg, BroadcastRPCReply{
		Type: "topology_ok",
	})
}
