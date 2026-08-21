package node

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type RPCReadMsg struct {
	Type     string  `json:"type"`
	Messages []int64 `json:"messages"`
}

func (n *Node) Read(msg maelstrom.Message) error {
	var body RPCReadMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return fmt.Errorf("Unable to decode Read message body - %w", err)
	}

	return n.maelstromNode.Reply(msg, RPCReply{
		Type:     "read_ok",
		Messages: n.messages.ReadAll(),
	})
}
