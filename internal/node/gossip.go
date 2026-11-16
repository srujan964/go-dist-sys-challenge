package node

import (
	"encoding/json"
	"fmt"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type GossipMsg struct {
	Type string  `json:"type"`
	Has  []int64 `json:"has"`
}

type Gossip struct {
	msg  GossipMsg
	dest string
}

func (n *Node) ReplyToGossip(msg maelstrom.Message) error {
	var body GossipMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("Unable to decode Gossip message - %w", err)
	}

	n.mu.Lock()
	hasSet := make(map[int64]struct{}, len(body.Has))
	for _, it := range body.Has {
		hasSet[it] = struct{}{}
		if _, exists := n.messages.elements[it]; !exists {
			n.messages.Add(it)
		}
	}
	delta := make([]int64, 0)
	for it := range n.messages.elements {
		if _, exists := hasSet[it]; !exists {
			delta = append(delta, it)
		}
	}
	n.mu.Unlock()

	outgoing := GossipMsg{
		Type: "gossip_ok",
		Has:  delta,
	}
	return n.maelstromNode.Send(msg.Src, outgoing)
}

func (n *Node) HandleGossipAck(msg maelstrom.Message) error {
	var body GossipMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("Unable to decode Gossip message - %w", err)
	}

	n.mu.Lock()
	for _, it := range body.Has {
		if _, exists := n.messages.elements[it]; !exists {
			n.messages.elements[it] = struct{}{}
		}
	}
	n.mu.Unlock()
	return nil
}

// Gossip messages with neighbours over `interval` milliseconds.
func (n *Node) Gossip(interval time.Duration) error {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-n.done:
				return
			case <-ticker.C:
				n.gossip()
			}
		}
	}()
	return nil
}

func (n *Node) gossip() error {
	n.mu.RLock()
	msg := GossipMsg{
		Type: "gossip",
		Has:  n.messages.Elements(),
	}
	n.mu.RUnlock()

	for _, neighbour := range n.neighbours {
		n.maelstromNode.Send(neighbour, msg)
	}
	return nil
}
