package glow

import (
	"testing"
	"time"
)

func TestOneNotification(t *testing.T) {
	called := false
	nf := make(Notifier)
	nf.On("ping", func(Message) { called = true })

	if called {
		t.Error("event fired too soon")
	}

	nf.Notify("ping")

	if !called {
		t.Error("event did not fire")
	}
}

func TestMultipleNotificationHandlers(t *testing.T) {
	calls := 0
	nf := make(Notifier)
	nf.On("ping", func(Message) { calls += 1 })
	nf.On("ping", func(Message) { calls += 10 })

	nf.Notify("ping")

	if calls != 11 {
		t.Error("expected 11, got:", calls)
	}
}

func TestDifferentNotifications(t *testing.T) {
	calls := 0
	nf := make(Notifier)
	nf.On("ping", func(Message) { calls += 1 })
	nf.On("pong", func(Message) { calls += 10 })
	nf.On("blah", func(Message) { calls += 100 })

	nf.Notify("ping")

	if calls != 1 {
		t.Error("expected 1, got:", calls)
	}

	nf.Notify("pong")

	if calls != 11 {
		t.Error("expected 11, got:", calls)
	}
}

func TestDifferentAndMultipleNotifications(t *testing.T) {
	calls := 0
	nf := make(Notifier)
	nf.On("ping", func(Message) { calls += 1 })
	nf.On("pong", func(Message) { calls += 10 })
	nf.On("pong", func(Message) { calls += 100 })

	nf.Notify("ping")
	nf.Notify("pong")
	nf.Notify("ping")
	nf.Notify("pong")
	nf.Notify("ping")

	if calls != 223 {
		t.Error("expected 223, got:", calls)
	}
}

func TestNotificationWithArgs(t *testing.T) {
	var args Message
	nf := make(Notifier)
	nf.On("ping", func(m Message) { args = m })

	nf.Notify("ping", 1, "a", nil)

	if args.String() != "1 a []" {
		t.Error("expected '1 a []', got:", args)
	}
}

func TestNotificationOff(t *testing.T) {
	called := false
	nf := make(Notifier)
	l := nf.On("ping", func(Message) { called = true })

	nf.Notify("ping")

	if !called {
		t.Error("event did not fire")
	}

	nf.Off(l)
	called = false
	nf.Notify("ping")

	if called {
		t.Error("event fired again")
	}
}

func TestSleep(t *testing.T) {
	now := time.Now()
	t0 := Now
	Sleep(10)
	diff := Now - t0
	elapsed := time.Since(now)

	if diff != 10 {
		t.Error("expected 10, got:", diff)
	}
	if elapsed > time.Millisecond {
		t.Error("simulated time should be instant, was:", elapsed)
	}
}

func TestNoNextTimeout(t *testing.T) {
	if NextTimeout() >= 0 {
		t.Error("there should be no timeouts pending")
	}
}

func TestNextTimeout(t *testing.T) {
	t0 := Now
	SetTimeout(123, func(Message) { t.Error("should not happen") })

	if NextTimeout() != t0+123 {
		t.Error("expected", t0+123, "got:", NextTimeout())
	}
}
