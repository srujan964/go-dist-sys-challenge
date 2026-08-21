package node

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastMsg struct {
	Type    string `json:"type"`
	Message int64  `json:"message"`
}

type RPCReply struct {
	Type     string  `json:"type"`
	Messages []int64 `json:"messages,omitempty"`
}

func (n *Node) Broadcast(msg maelstrom.Message) error {
	var body BroadcastMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("Unable to decode Broadcast message body - %w", err)
	}

	if n.messages.Contains(body.Message) {
		return n.maelstromNode.Reply(msg, RPCReply{
			Type: "broadcast_ok",
		})
	}
	n.messages.Add(body.Message)
	for _, d := range n.neighbours {
		if d == msg.Src {
			continue
		}
		broadcastMsg := BroadcastMsg{
			Type:    "broadcast",
			Message: body.Message,
		}
		n.maelstromNode.RPC(d, broadcastMsg, func(msg maelstrom.Message) error { return nil })
	}

	return n.maelstromNode.Reply(msg, RPCReply{
		Type: "broadcast_ok",
	})
}
