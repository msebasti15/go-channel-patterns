package blockingchannel

import (
	"log"
	"time"
)

type BlockingChannel struct {
	capacity      int
	receiveCh     <-chan interface{}
	sendCh        chan<- interface{}
	timeout       time.Duration
	timeoutTicker *time.Ticker
}

// NewBlockingChannel creates a new blocking returns *BlockingChannel
func NewBlockingChannel(capacity int, timeout time.Duration) *BlockingChannel {
	ch := make(chan interface{}, capacity)

	return &BlockingChannel{
		receiveCh:     ch,
		sendCh:        ch, // using the same channel for demo simplicity
		timeout:       timeout,
		timeoutTicker: time.NewTicker(timeout),
	}
}

// Send sends a message and keeps retrying until it has success
// if the channel is full, it would keep retrying sending the message until it succeeds.
// this is a simple demonstration, in a real-world scenario, you might want to add some flow control or logging.
// also, in a real-world scenario, you might want to handle the case where the channel is closed by the other side.
// for example, you might want to close the channel when the sender is done sending messages.
// and you might want to add error handling and retry logic in case of failures.
// this would be a complex task and might require additional libraries or frameworks.
func (b *BlockingChannel) Send(msg interface{}) {
	// if send returns true Send returns
	// if send returns false, Send would keep retrying sending the message.
	for !b.send(msg) {
		//this section would be used for logging and metrics or add some flow control
	}
}

// send sends a message and returns true if successful, false otherwise
func (b *BlockingChannel) send(msg interface{}) bool {
	b.timeoutTicker.Reset(b.timeout)
	defer b.timeoutTicker.Stop()

	select {
	case b.sendCh <- msg:
		log.Printf("msg[%v] sent", msg)
		return true

	case <-b.timeoutTicker.C:
		log.Printf("msg[%v] sending timed out, retrying", msg)
		return false
	}
}

// Process reads messages from the channel and calls the provided function for each message
func (b *BlockingChannel) Process(f func(msg interface{})) {
	for msg := range b.receiveCh {
		f(msg)
	}
}
