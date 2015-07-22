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
	t.Skip("TODO")
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
