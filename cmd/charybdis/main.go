package main

import (
	"github.com/srujan964/go-dist-sys-challenge/internal/node"
	"log"

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

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
