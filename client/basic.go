package main

import (
	"fmt"
	"os"
	"time"

	nats "github.com/nats-io/nats.go"
)

func main() {
	// Connect to a server
	nc, status := nats.Connect(nats.DefaultURL)

	fmt.Println(status)
	if status != nil {
		fmt.Println("Unable to connect to nats server")
		os.Exit(0)
	}

	// Simple Publisher
	nc.Publish("foo", []byte("Hello World"))

	// Simple Async Subscriber
	nc.Subscribe("foo", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	// Responding to a request message
	nc.Subscribe("request", func(m *nats.Msg) {
		m.Respond([]byte("answer is 42"))
	})

	// Simple Sync Subscriber
	sub, err := nc.SubscribeSync("foo")
	m, err := sub.NextMsg(10000 * time.Millisecond)

	// Channel Subscriber
	ch := make(chan *nats.Msg, 64)
	sub, err = nc.ChanSubscribe("foo", ch)
	msg := <-ch

	// Unsubscribe
	sub.Unsubscribe()

	// Drain
	sub.Drain()

	// Requests
	msg, err = nc.Request("help", []byte("help me"), 10*time.Millisecond)

	fmt.Println(msg)
	fmt.Println(m)
	fmt.Println(err)

	// Replies
	nc.Subscribe("help", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("I can help!"))
	})

	// Drain connection (Preferred for responders)
	// Close() not needed if this is called.
	nc.Drain()

	// Close connection
	nc.Close()
}
