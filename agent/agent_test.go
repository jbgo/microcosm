package main

import (
	"fmt"
	docker "github.com/jbgo/microcosm/dockerclient"
	"testing"
	"time"
)

type logEventAction struct {
	Events chan string
}

func (action *logEventAction) Invoke(payload *Payload) error {
	msg := fmt.Sprintf("ID: %s, Event: %s", payload.Container.ID, payload.DockerEvent.Status)
	action.Events <- msg
	return nil
}

func TestRegister(t *testing.T) {
	agent := New(&FakeDockerClient{})

	handler := &EventHandler{
		Matcher: NewEventMatcher("proxy", "start", "stop"),
		Action:  &logEventAction{},
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
	agent := New(&FakeDockerClient{})

	handler := &EventHandler{
		Matcher: NewEventMatcher("web", "start", "restart"),
		Action:  &logEventAction{},
	}

	otherHandler := &EventHandler{
		Matcher: NewEventMatcher("db", "stop"),
		Action:  &logEventAction{},
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

type FakeDockerClient struct {
	Containers map[string]*docker.Container
	Listeners  []chan<- *docker.Event
}

func newFakeDockerClient() *FakeDockerClient {
	d := &FakeDockerClient{}

	web := &docker.Container{
		ID:     "aweb",
		Labels: map[string]string{"microcosm.type": "web"},
	}

	proxy := &docker.Container{
		ID:     "aproxy",
		Labels: map[string]string{"microcosm.type": "proxy"},
	}

	d.Containers = map[string]*docker.Container{
		"aproxy": proxy,
		"aweb":   web,
	}

	return d
}

func (d *FakeDockerClient) InspectContainer(id string) (*docker.Container, error) {
	return d.Containers[id], nil
}

func (d *FakeDockerClient) DispatchEvent(containerID, eventStatus string) {
	event := &docker.Event{ContainerID: containerID, Status: eventStatus}
	for _, listener := range d.Listeners {
		listener <- event
	}
}

func (d *FakeDockerClient) AddEventListener(listener chan<- *docker.Event) error {
	d.Listeners = append(d.Listeners, listener)
	return nil
}

func failOnAnyError(t *testing.T, agent *Agent) {
	for {
		err := <-agent.Errors
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestListen(t *testing.T) {
	client := newFakeDockerClient()
	agent := New(client)

	action := &logEventAction{Events: make(chan string, 10)}

	handler1 := &EventHandler{
		Matcher: NewEventMatcher("proxy", "start", "stop"),
		Action:  action,
	}

	handler2 := &EventHandler{
		Matcher: NewEventMatcher("web", "start", "restart"),
		Action:  action,
	}

	agent.Register(handler1)
	agent.Register(handler2)

	go agent.Listen()

	// The agent will send a nil error message to indicate it is listening
	// for events, otherwise this will be an actual error.
	err := <-agent.Errors
	if err != nil {
		t.Fatal("agent not listening", err)
	}
	go failOnAnyError(t, agent)

	client.DispatchEvent("aproxy", "start")   // match
	client.DispatchEvent("aweb", "stop")      // no match
	client.DispatchEvent("aproxy", "restart") // no match
	client.DispatchEvent("aweb", "start")     // match
	client.DispatchEvent("aproxy", "stop")    // match
	client.DispatchEvent("aweb", "restart")   // match
	client.DispatchEvent("aproxy", "start")   // match
	client.DispatchEvent("aweb", "stop")      // no match

	expected := []string{
		"ID: aproxy, Event: start",
		"ID: aweb, Event: start",
		"ID: aproxy, Event: stop",
		"ID: aweb, Event: restart",
		"ID: aproxy, Event: start",
	}

	// wait for events to process, then close channels and continue
	go func() {
		time.Sleep(1 * time.Millisecond)
		close(agent.Channel)
		close(action.Events)
	}()

	actual := []string{}
	for {
		msg, ok := <-action.Events
		if !ok {
			break
		}
		actual = append(actual, msg)
	}

	if len(actual) != len(expected) {
		t.Fatalf("expected %d events, got %d", len(expected), len(actual))
	}

	for i := range actual {
		if actual[i] != expected[i] {
			t.Errorf("expected event with %s, got %s", expected[i], actual[i])
		}
	}
}
