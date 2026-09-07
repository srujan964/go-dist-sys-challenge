package node

import (
	"encoding/json"
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastMsg struct {
	Type    string `json:"type"`
	Message int64  `json:"message"`
}

type RPCReply struct {
	Type      string  `json:"type"`
	InReplyTo int64   `json:"in_reply_to,omitempty"`
	Messages  []int64 `json:"messages,omitempty"`
}

func extractMessage(msg maelstrom.Message) (int64, error) {
	var body BroadcastMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return 0, fmt.Errorf("Unable to decode Broadcast message body - %w", err)
	}
	return body.Message, nil
}

// Partition tolerance: assume a message does not reach a node.
//
// We can attach a handler for every message sent to a neighbouring node,
// that deletes the configured retry for this message to the destination node.
func (n *Node) Broadcast(msg maelstrom.Message) error {
	value, err := extractMessage(msg)
	if err != nil {
		return err
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.messages.Contains(value) {
		return n.maelstromNode.Reply(msg, RPCReply{
			Type: "broadcast_ok",
		})
	}

	n.messages.Add(value)

	for _, d := range n.neighbours {
		if d == msg.Src {
			continue
		}
		broadcastMsg := BroadcastMsg{
			Type:    "broadcast",
			Message: value,
		}
		err := n.maelstromNode.Send(d, broadcastMsg)
		if err != nil {
			log.Fatalf("Failed to broadcast to neighbour - %v", err)
		}
	}
	return n.maelstromNode.Reply(msg, RPCReply{
		Type: "broadcast_ok",
	})
}

func (n *Node) BroadcastAck(msg maelstrom.Message) error {
	return nil
}
