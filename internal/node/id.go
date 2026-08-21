package node

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type IdGenerateMsg struct {
	Type  string `json:"type"`
	MsgId int64  `json:"msg_id"`
}

type IdGenerateReply struct {
	Type      string `json:"type"`
	InReplyTo int64  `json:"in_reply_to"`
	Id        string `json:"id"`
}

const (
	nodeBits     = 6
	sequenceBits = 16
	nodeMask     = (1 << nodeBits) - 1
	sequenceMask = (1 << sequenceBits) - 1
)

// Generate a node based growing ID
func (n *Node) GenerateId(msg maelstrom.Message) error {
	var body IdGenerateMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		log.Fatalf("Unable to decode message body - %v", err)
		return err
	}

	id, err := n.nextId(msg.Dest)
	if err != nil {
		log.Fatalf("Error generating id - %v", err)
		return err
	}

	reply := IdGenerateReply{
		Type:      "generate_ok",
		InReplyTo: body.MsgId,
		Id:        id,
	}

	return n.maelstromNode.Reply(msg, reply)
}

// Generate a 64 bit snowflake-esque id based on the current timestamp, node id,
// and a local counter.
//
// 41 bits for timestamp -- milliseconds since epoch
// 6 bit node ID -- assumes no more than 64 nodes
// 16 bit node local sequence -- assumes no more than 65536 requests to a node per millisecond
func (n *Node) nextId(node string) (string, error) {
	n.seq++

	var id uint64
	parsedNodeId, err := strconv.Atoi(strings.TrimPrefix(node, "n"))
	if err != nil {
		log.Printf("Cannot convert node id to int - %v", err)
		return "", err
	}

	tmstmp := uint64(time.Now().UnixMilli())
	nodeId := uint64(parsedNodeId) & nodeMask

	id = (tmstmp << (nodeBits + sequenceBits)) |
		(nodeId << (sequenceBits)) |
		(n.seq & sequenceMask)

	guid := fmt.Sprintf("%X", id)
	return guid, nil
}
