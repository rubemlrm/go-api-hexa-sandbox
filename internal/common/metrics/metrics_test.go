package metrics

import (
	"os"
	"testing"
	"time"
)

func TestInitMeter(t *testing.T) {
	// Run InitMeter in a goroutine since it blocks on context.Done()
	go func() {
		InitMeter()
	}()

	time.Sleep(500 * time.Millisecond)

	// Send interrupt signal to unblock context.Done()
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("failed to find process: %v", err)
	}
	if err := p.Signal(os.Interrupt); err != nil {
		t.Fatalf("failed to send interrupt: %v", err)
	}

	// Wait a bit to allow goroutine to exit
	time.Sleep(500 * time.Millisecond)
}
