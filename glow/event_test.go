package glow

import "testing"

func TestOneEvent(t *testing.T) {
	called := false
	ee := make(EventEmitter)
	ee.On("ping", func(Msg) { called = true })

	if called {
		t.Error("event fired too soon")
	}

	ee.Emit("ping")

	if !called {
		t.Error("event did not fire")
	}
}

func TestMultipleEventHandlers(t *testing.T) {
	calls := 0
	ee := make(EventEmitter)
	ee.On("ping", func(Msg) { calls += 1 })
	ee.On("ping", func(Msg) { calls += 10 })

	ee.Emit("ping")

	if calls != 11 {
		t.Error("expected 11, got:", calls)
	}
}

func TestMultipleEvents(t *testing.T) {
	calls := 0
	ee := make(EventEmitter)
	ee.On("ping", func(Msg) { calls += 1 })
	ee.On("pong", func(Msg) { calls += 10 })
	ee.On("blah", func(Msg) { calls += 100 })

	ee.Emit("ping")

	if calls != 1 {
		t.Error("expected 1, got:", calls)
	}

	ee.Emit("pong")

	if calls != 11 {
		t.Error("expected 11, got:", calls)
	}
}

func TestMultipleEventsAndHandlers(t *testing.T) {
	calls := 0
	ee := make(EventEmitter)
	ee.On("ping", func(Msg) { calls += 1 })
	ee.On("pong", func(Msg) { calls += 10 })
	ee.On("pong", func(Msg) { calls += 100 })

	ee.Emit("ping")
	ee.Emit("pong")
	ee.Emit("ping")
	ee.Emit("pong")
	ee.Emit("ping")

	if calls != 223 {
		t.Error("expected 223, got:", calls)
	}
}

func TestEventWithArgs(t *testing.T) {
	var args Msg
	ee := make(EventEmitter)
	ee.On("ping", func(m Msg) { args = m })

	ee.Emit("ping", 1, "a", 2)

	if args.String() != "1 a 2" {
		t.Error("expected '1 a 2', got:", args)
	}
}
