package node

import (
	"context"
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type GcounterRPC struct {
	Type  string `json:"type"`
	Delta int    `json:"delta"`
}

type GcounterRPCReply struct {
	Type  string `json:"type"`
	Value *int   `json:"value,omitempty"`
}

func (n *Node) AddCounter(msg maelstrom.Message) error {
	var body GcounterRPC
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("Unable to decode GcounterRPC message body - %w", err)
	}
	ctx := context.Background()

	updated := n.latest + body.Delta
	n.kv.CompareAndSwap(ctx, n.maelstromNode.ID(), n.latest, updated, true)
	n.latest = updated
	return n.maelstromNode.Reply(msg, GcounterRPCReply{
		Type: "add_ok",
	})
}

func (n *Node) ReadCounter(msg maelstrom.Message) error {
	var body GcounterRPC
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("Unable to decode GcounterRPC message body - %w", err)
	}
	ctx := context.Background()

	var total int = 0
	for _, x := range n.maelstromNode.NodeIDs() {
		val, err := n.kv.ReadInt(ctx, x)
		if err != nil {
			val = 0
		}
		total += val
	}
	return n.maelstromNode.Reply(msg, GcounterRPCReply{
		Type:  "read_ok",
		Value: &total,
	})
}
