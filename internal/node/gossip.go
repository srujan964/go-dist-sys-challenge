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
	received := SetFrom(body.Has)
	n.mu.RLock()
	unseen := received.Difference(n.messages)
	n.mu.RUnlock()

	if unseen.Size() > 0 {
		n.mu.Lock()
		for _, msg := range unseen.Elements() {
			n.messages.Add(msg)
		}
		n.mu.Unlock()
	}

	n.mu.RLock()
	delta := n.messages.Difference(received)
	n.mu.RUnlock()

	outgoing := GossipMsg{
		Type: "gossip_ok",
		Has:  delta.Elements(),
	}

	n.maelstromNode.Send(msg.Src, outgoing)

	return nil
}

func (n *Node) HandleGossipAck(msg maelstrom.Message) error {
	var body GossipMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("Unable to decode Gossip message - %w", err)
	}
	received := SetFrom(body.Has)
	n.mu.RLock()
	unseen := received.Difference(n.messages)
	n.mu.RUnlock()

	if unseen.Size() > 0 {
		n.mu.Lock()
		for _, msg := range unseen.Elements() {
			n.messages.Add(msg)
		}
		n.mu.Unlock()
	}
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
