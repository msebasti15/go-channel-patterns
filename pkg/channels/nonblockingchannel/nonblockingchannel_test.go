package nonblockingchannel

import (
	"testing"
	"time"
)

const timeout = 5 * time.Microsecond

func TestNonBlockingChannels_send_no_blocks(t *testing.T) {
	// Arrange
	bc := NewNonBlockingChannel(1, timeout)

	go bc.Process(
		noDelay,
	)

	msg1 := "1st message"
	msg2 := "2nd message"

	// Act
	sent1 := bc.Send(msg1)

	time.Sleep(10 * time.Millisecond)

	sent2 := bc.Send(msg2)

	time.Sleep(20 * time.Microsecond)
	// Assert
	if !sent1 {
		t.Error("Expected send to be successful msg1")
	}

	if !sent2 {
		t.Error("Expected send to be successful msg2")
	}

}

func TestNonBlockingChannels_send_blocks(t *testing.T) {
	// Arrange
	bc := NewNonBlockingChannel(1, timeout)

	msg1 := "1st message"
	msg2 := "2nd message"

	// Act
	sent1 := bc.Send(msg1)
	sent2 := bc.Send(msg2)

	time.Sleep(5 * time.Microsecond)
	// Assert
	if !sent1 {
		t.Error("Expected send to be successful msg1")
	}

	if sent2 {
		t.Error("Expected send to be blocked msg2")
	}
}

func TestNonBlockingChannelsAllMessageSent(t *testing.T) {
	// Arrange
	bc := NewNonBlockingChannel(1, timeout)

	go bc.Process(
		noDelay,
	)

	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	// Act
	for idx := 0; idx < 10; idx++ {
		bc.Send(idx)
		<-ticker.C
	}
}

func TestNonBlockingChannelsWithDroppedMessages(t *testing.T) {
	// Arrange
	bc := NewNonBlockingChannel(1, timeout)

	go bc.Process(
		withDelay,
	)

	ticker := time.NewTicker(timeout)
	defer ticker.Stop()

	// Act
	for idx := 0; idx < 10; idx++ {
		bc.Send(idx)
		<-ticker.C
	}
}

func noDelay(msg interface{}) {
	// do nothing
}

func withDelay(msg interface{}) {
	time.Sleep(2 * timeout)
}
