package main

import (
	"testing"
)

func TestRegister(t *testing.T) {
	agent := New()

	handler := &EventHandler{
		Matcher: NewEventMatcher("proxy", "start", "stop"),
		Command: &Command{},
	}

	agent.Register(handler)

	if len(agent.Handlers) != 1 {
		t.Errorf("expecting 1 registered handler, got %d", len(agent.Handlers))
		t.FailNow()
	}

	registeredHandler := agent.Handlers[0]
	if registeredHandler != handler {
		t.Error("expected", handler, "got", registeredHandler)
	}
}

func TestCancel(t *testing.T) {
	agent := New()

	handler := &EventHandler{
		Matcher: NewEventMatcher("web", "start", "restart"),
		Command: &Command{},
	}

	otherHandler := &EventHandler{
		Matcher: NewEventMatcher("db", "stop"),
		Command: &Command{},
	}

	agent.Register(handler)
	agent.Register(otherHandler)
	cancelledHandler := agent.Cancel(handler)

	if len(agent.Handlers) != 1 {
		t.Errorf("expecting 1 registered handler, got %d", len(agent.Handlers))
	}

	if cancelledHandler != handler {
		t.Error("expected cancelled handler to be", handler, "got", cancelledHandler)
	}
}

func TestFindTasksForEvent(t *testing.T) {
	t.Skip("TODO")
}

func TestExecuteTask(t *testing.T) {
	t.Skip("TODO")
}

func TestListen(t *testing.T) {
	t.Skip("TODO")
}
