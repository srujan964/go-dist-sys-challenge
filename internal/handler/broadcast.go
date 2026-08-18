package handler

import (
	"encoding/json"
	"fmt"

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
	messages []int64
	topology map[string][]string
}

func NewState() *State {
	return &State{
		messages: make([]int64, 0),
	}
}

func (s *State) Store(msg maelstrom.Message, n *maelstrom.Node) (RPCReply, error) {
	var body BroadcastMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return RPCReply{}, fmt.Errorf("Unable to decode Broadcast message body - %w", err)
	}
	s.messages = append(s.messages, body.Message)
	return RPCReply{
		Type: "broadcast_ok",
	}, nil
}

func (s *State) CurrentMessages() (RPCReply, error) {
	return RPCReply{
		Type:     "read_ok",
		Messages: s.messages,
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
