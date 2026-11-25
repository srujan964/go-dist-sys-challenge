package node

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

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

func RandomNeighbours(neighbours []string) []string {
	rand.Shuffle(len(neighbours), func(i, j int) {
		neighbours[i], neighbours[j] = neighbours[j], neighbours[i]
	})
	m := len(neighbours)/2 + 1
	return neighbours[:m]
}

func (n *Node) Broadcast(msg maelstrom.Message) error {
	value, err := extractMessage(msg)
	if err != nil {
		return err
	}
	n.mu.Lock()
	if n.messages.Contains(value) {
		n.mu.Unlock()
		return n.maelstromNode.Reply(msg, RPCReply{
			Type: "broadcast_ok",
		})
	}
	n.messages.Add(value)
	n.mu.Unlock()

	err = n.maelstromNode.Reply(msg, RPCReply{
		Type: "broadcast_ok",
	})
	if err != nil {
		return err
	}

	for _, d := range RandomNeighbours(n.neighbours) {
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
	return nil
}

func (n *Node) BroadcastAck(msg maelstrom.Message) error {
	return nil
}
