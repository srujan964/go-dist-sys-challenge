package handler

import (
	"encoding/json"
	"fmt"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type BroadcastMsg struct {
	Type    string `json:"type"`
	Message int64  `json:"message"`
}

type RPCReadMsg struct {
	Type     string  `json:"type"`
	Messages []int64 `json:"messages"`
}

type RPCToplogyMsg struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type RPCReply struct {
	Type     string  `json:"type"`
	Messages []int64 `json:"messages,omitempty"`
}

type State struct {
	mu       *sync.Mutex
	topology map[string][]string
	messages map[int64]int64
}

func NewState() *State {
	return &State{
		mu:       &sync.Mutex{},
		messages: make(map[int64]int64),
	}
}

func (s *State) Broadcast(msg maelstrom.Message, n *maelstrom.Node) (RPCReply, error) {
	var body BroadcastMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return RPCReply{}, fmt.Errorf("Unable to decode Broadcast message body - %w", err)
	}

	if seenMessage(s, body.Message) {
		return RPCReply{
			Type: "broadcast_ok",
		}, nil
	}

	if dest, ok := s.topology[msg.Dest]; ok {
		for _, d := range dest {
			if d == msg.Src {
				continue
			}
			broadcastMsg := BroadcastMsg{
				Type:    "broadcast",
				Message: body.Message,
			}
			n.RPC(d, broadcastMsg, func(msg maelstrom.Message) error { return nil })
		}
	}
	return RPCReply{
		Type: "broadcast_ok",
	}, nil
}

func seenMessage(s *State, value int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.messages[value]
	if !exists {
		s.messages[value] = 0
	}
	return exists
}

func (s *State) CurrentMessages() (RPCReply, error) {
	msgs := make([]int64, 0)
	for k := range s.messages {
		msgs = append(msgs, k)
	}

	return RPCReply{
		Type:     "read_ok",
		Messages: msgs,
	}, nil
}

func (s *State) ApplyTopology(msg maelstrom.Message) (RPCReply, error) {
	var body RPCToplogyMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return RPCReply{}, fmt.Errorf("Unable to decode Topology message body - %v", err)
	}
	s.topology = body.Topology
	return RPCReply{
		Type: "topology_ok",
	}, nil
}
