package rate

import (
	"testing"
	"time"
)

func TestLimiter_Allow(t *testing.T) {
	lim := NewLimiter(1*time.Second, 3)

	// Test allowing 3 requests within 1 second
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}

	// Test denying the 4th request within the same second
	if lim.Allow() {
		t.Errorf("expected false, got true")
	}
}

func TestLimiter_Allow_WithTimeAdvance(t *testing.T) {
	lim := NewLimiter(1*time.Second, 3)

	// Allow 3 requests in the first second
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}

	time.Sleep(2 * time.Second)

	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}

	if lim.Allow() {
		t.Errorf("expected false, got true")
	}
}

func TestLimiter_Allow_WithPartialWindow(t *testing.T) {
	lim := NewLimiter(1*time.Second, 3)

	// Allow 2 requests in the first second
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}

	time.Sleep(600 * time.Millisecond)

	// Allow 1 request, which should be allowed since we're only halfway through the window
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}

	// Deny the 2nd request in this half-window
	if lim.Allow() {
		t.Errorf("expected false, got true")
	}
}

func TestLimiter_Allow_ExactLimit(t *testing.T) {
	lim := NewLimiter(1*time.Second, 3) // 3 events per second

	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}
	if !lim.Allow() {
		t.Errorf("expected true, got false")
	}

	// Deny the 4th request
	if lim.Allow() {
		t.Errorf("expected false, got true")
	}
}

func TestLimiter_Allow_NegativeCases(t *testing.T) {
	lim := NewLimiter(1*time.Second, 3) // 3 events per second

	// Try to allow a large number of requests at once, which should be denied
	if lim.allowN(time.Now(), 5) {
		t.Errorf("expected false, got true")
	}

	// Allow a smaller number of requests that exceeds the limit after aggregation
	lim.allowN(time.Now(), 2) // 2 events
	if lim.allowN(time.Now(), 2) {
		t.Errorf("expected false, got true")
	}
}
