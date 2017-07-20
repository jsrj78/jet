package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/jeelabs/jet/glow"
)

func TestOneNotification(t *testing.T) {
	called := false
	nf := glow.MakeNotifier()
	nf.On("ping", func(glow.Message) { called = true })

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
	nf := glow.MakeNotifier()
	nf.On("ping", func(glow.Message) { calls += 1 })
	nf.On("ping", func(glow.Message) { calls += 10 })

	nf.Notify("ping")

	if calls != 11 {
		t.Error("expected 11, got:", calls)
	}
}

func TestDifferentNotifications(t *testing.T) {
	calls := 0
	nf := glow.MakeNotifier()
	nf.On("ping", func(glow.Message) { calls += 1 })
	nf.On("pong", func(glow.Message) { calls += 10 })
	nf.On("blah", func(glow.Message) { calls += 100 })

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
	nf := glow.MakeNotifier()
	nf.On("ping", func(glow.Message) { calls += 1 })
	nf.On("pong", func(glow.Message) { calls += 10 })
	nf.On("pong", func(glow.Message) { calls += 100 })

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
	var args glow.Message
	nf := glow.MakeNotifier()
	nf.On("ping", func(m glow.Message) { args = m })

	nf.Notify("ping", 1, "a", nil)

	if args.String() != "1 a []" {
		t.Error("expected '1 a []', got:", args)
	}
}

func TestNotificationOff(t *testing.T) {
	called := false
	nf := glow.MakeNotifier()
	l := nf.On("ping", func(glow.Message) { called = true })

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

func TestRunning(t *testing.T) {
	glow.Now = 0
	now := time.Now()
	glow.Run(10)
	elapsed := time.Since(now)

	if glow.Now != 10 {
		t.Error("expected 10, got:", glow.Now)
	}
	if elapsed > time.Millisecond {
		t.Error("simulated time should be instant, was:", elapsed)
	}
}

func TestNoNextTimer(t *testing.T) {
	if glow.NextTimer >= 0 {
		t.Error("there should be no timeouts pending")
	}
}

func TestMultipleTimers(t *testing.T) {
	glow.Now = 0
	t1, t2, t3 := -1, -1, -1
	glow.SetTimer(123, func() { t1 = glow.Now })
	glow.SetTimer(789, func() { t2 = glow.Now })
	glow.SetTimer(456, func() { t3 = glow.Now })

	if glow.NextTimer != 123 {
		t.Error("expected", 123, "got:", glow.NextTimer)
	}

	glow.Run(1000)

	if glow.Now != 1000 {
		t.Error("expected 1000, got:", glow.Now)
	}

	if t1 != 123 {
		t.Error("expected", 123, "got:", t1)
	}
	if t2 != 789 {
		t.Error("expected", 789, "got:", t2)
	}
	if t3 != 456 {
		t.Error("expected", 456, "got:", t3)
	}

	if glow.NextTimer >= 0 {
		t.Error("there should be no timeouts pending")
	}
}

func TestCancelledNextTimer(t *testing.T) {
	glow.Now = 0
	l := glow.SetTimer(123, func() {})

	if glow.NextTimer != 123 {
		t.Error("expected", 123, "got:", glow.NextTimer)
	}

	glow.CancelTimer(l)

	if glow.NextTimer >= 0 {
		t.Error("there should be no timeouts pending")
	}
}

func TestPeriodicimer(t *testing.T) {
	glow.Now = 0
	v := []int{}
	glow.SetPeriodic(123, func() { v = append(v, glow.Now) })

	if glow.NextTimer != 123 {
		t.Error("expected", 123, "got:", glow.NextTimer)
	}

	defer glow.Stop() // TODO manually stop all timers
	glow.Run(500)

	if fmt.Sprint(v) != "[123 246 369 492]" {
		t.Error("expected '[123 246 369 492]', got:", fmt.Sprint(v))
	}
	if glow.NextTimer != 615 {
		t.Error("expected 615, got:", glow.NextTimer)
	}
}
