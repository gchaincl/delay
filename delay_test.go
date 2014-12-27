package delay

import (
	"testing"
	"time"
)

func TestNewDelayerStartsWithNoPendingTimers(t *testing.T) {
	d := NewDelayer(func(key, payload string) {}, 0)

	got := d.Pending()
	if got != 0 {
		t.Errorf("Expected 0 pending, got %d\n", got)
	}
}

func TestRegisterTriggersCallbackAfterSpecifiedTime(t *testing.T) {
	sem := make(chan bool)
	cb := func(key, payload string) {
		sem <- true
	}

	d := NewDelayer(cb, 10*time.Millisecond)

	d.Register("a", "message")
	start := time.Now()
	<-sem

	elapsed := time.Since(start)
	if elapsed <= 10*time.Millisecond {
		t.Errorf("Callback should run after 10ms, %fms elasped\n",
			float64(elapsed)/float64(time.Millisecond),
		)
	}
}

func TestRegisteredCallbacksAreCanceled(t *testing.T) {
	d := NewDelayer(func(key, payload string) {}, 1*time.Second)
	d.Register("a", "message")

	d.Cancel("a")

	got := d.Pending()
	if got != 0 {
		t.Errorf("It should not be pending callbacks, got %d\n", got)
	}
}

func TestCallbacksAreUpdated(t *testing.T) {
	c := make(chan string)

	d := NewDelayer(func(key, payload string) {
		c <- payload
	}, 10*time.Millisecond)

	d.Register("a", "1")
	d.Register("a", "2")
	d.Register("a", "3")

	got := <-c
	if got != "3" {
		t.Errorf("Message expected '%s', got '%s'", "3", got)
	}
}
