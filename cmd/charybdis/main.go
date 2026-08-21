package main

import (
	"github.com/srujan964/go-dist-sys-challenge/internal/handler"
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
	idGenerator := handler.NewGenerator()

	state := handler.NewState()

	n.Handle("echo", func(msg maelstrom.Message) error {
		resp, err := handler.HandleEcho(msg)
		if err != nil {
			log.Fatal(err)
		}

		return n.Reply(msg, resp)
	})

	n.Handle("generate", func(msg maelstrom.Message) error {
		resp, err := idGenerator.Generate(msg)
		if err != nil {
			log.Fatal(err)
		}

		return n.Reply(msg, resp)
	})

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		resp, err := state.Broadcast(msg, n)
		if err != nil {
			log.Fatal(err)
		}

		return n.Reply(msg, resp)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		resp, err := state.CurrentMessages()
		if err != nil {
			log.Fatal(err)
		}

		return n.Reply(msg, resp)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		resp, err := state.ApplyTopology(msg)
		if err != nil {
			log.Fatal(err)
		}

		return n.Reply(msg, resp)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
