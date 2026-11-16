package main

import (
	"log"
	"time"

	"github.com/srujan964/go-dist-sys-challenge/internal/node"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// A Maelstrom node receives messages in the form:
// {
//   "src": "c1",
//   "dest": "n1",
//   "body": {
//     "type": "echo",
//     "msg_id": 1,
//     "echo": "Please echo 35"
//   }
// }

func main() {
	n := maelstrom.NewNode()
	node := node.CreateNode(n)

	n.Handle("echo", node.Echo)
	n.Handle("generate", node.GenerateId)
	n.Handle("broadcast", node.Broadcast)
	n.Handle("read", node.Read)
	n.Handle("topology", node.ApplyTopology)
	n.Handle("gossip", node.ReplyToGossip)
	n.Handle("gossip_ok", node.HandleGossipAck)
	n.Handle("broadcast_ok", node.BroadcastAck)

	node.Gossip(time.Duration(250 * time.Millisecond))

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

	node.Shutdown()
}
