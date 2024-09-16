package nonblockingchannel

import (
	"log"
	"time"
)

type NonBlockingChannel struct {
	capacity  int
	timeout   time.Duration
	receiveCh <-chan interface{}
	sendCh    chan<- interface{}
}

func NewNonBlockingChannel(capacity int, timeout time.Duration) *NonBlockingChannel {
	ch := make(chan interface{}, capacity)
	return &NonBlockingChannel{
		receiveCh: ch,
		sendCh:    ch,
		timeout:   timeout,
	}
}

// Send returns true if the message was sent ok
// If the channel is full, it drops the message and returns false.
// For a more sophisticated non-blocking behavior, a channel buffer size greater than 0 can be used.
// This would require additional logic to handle the case when the channel is full and the sender is faster than the receiver.
// In a real-world scenario, this would involve implementing a queue or a buffer mechanism.
// However, for simplicity and demonstration purposes, this example does not include such a feature.
// In a real-world application, you would want to handle this situation in a more robust way.
func (nc *NonBlockingChannel) Send(msg interface{}) bool {
	select {
	case nc.sendCh <- msg:
		log.Printf("msg[%v] send", msg)
		return true
	default:
		// message can be dropped, but other clever mechanism for failover can be done here
		log.Printf("channel full, msg[%v] drop", msg)
		return false
	}
}

func (nc *NonBlockingChannel) Process(f func(msg interface{})) {
	for msg := range nc.receiveCh {
		f(msg)
	}
}
