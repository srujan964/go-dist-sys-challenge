# go-dist-sys-challenge
Solving Fly.io's distributed systems challenges

The first three challenges are run under the `charybdis` cmd which wraps a Maelstrom node with handlers for Echo, ID Generation, and Broadcasting, along with a Gossip helper.

Challenge four (Grow-only counter) is built as its own command `counter`.

## Usage

Build the project as you usually would a Go project.
```
go build -o build ./cmd/...
```

Run the first three challenges with:
```
   ./maelstrom test -w echo --bin build/charybdis --node-count 1 --time-limit 10
```

Run challenge four with:
```
  ./maelstrom test -w g-counter --bin build/counter --node-count 3 --rate 100 --time-limit 20 --nemesis partition
```

### Notes

1) Globally unique ID
  Each node generates a Snowflake-like ID that uses the node's ID (numerical part), a local, node-only, increasing sequence, and a microsecond-precision timestamp.

2) Broadcast
  The broadcast solution started out forwarding every incoming message to all other neighbours as defined by the topology provided by Maelstrom.    

    To tolerate network partitions, each node also gossips its current state of messages to its neighbours in an interval.
    Node 'A' gossips its entire message state to node 'B', and node 'B' replies with all the messages it has that 'A' hasn't seen yet.

    The final solution for this challenge, requiring efficient broadcasting, moved to a custom cluster topology and randomized broadcasts to a majority of neighbours.
    
