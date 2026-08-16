package handler

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"log"
)

type EchoMsg struct {
	Type  string `json:"type"`
	MsgId int64  `json:"msg_id"`
	Echo  string `json:"echo"`
}

type EchoReply struct {
	Type      string `json:"type"`
	MsgId     int64  `json:"msg_id"`
	InReplyTo int64  `json:"in_reply_to"`
	Echo      string `json:"echo"`
}

func HandleEcho(msg maelstrom.Message) (EchoReply, error) {
	var body EchoMsg
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		log.Fatalf("Unable to decode message body - %v", err)
		return EchoReply{}, err
	}

	reply := EchoReply{
		Type:      "echo_ok",
		MsgId:     1,
		InReplyTo: body.MsgId,
		Echo:      body.Echo,
	}

	return reply, nil
}
