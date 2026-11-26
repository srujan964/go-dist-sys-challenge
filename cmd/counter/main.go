package main

import (
	"log"

	"github.com/srujan964/go-dist-sys-challenge/internal/node"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)

	node := node.CreateNode(n, kv)

	n.Handle("add", node.AddCounter)
	n.Handle("read", node.ReadCounter)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}

	node.Shutdown()
}
