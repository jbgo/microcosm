package main

import (
	"fmt"
	docker "github.com/jbgo/microcosm/dockerclient"
)

type EventHandler struct {
	Matcher *EventMatcher
	Action  Action
}

type Payload struct {
	Container   *docker.Container
	DockerEvent *docker.Event
	Handler     *EventHandler
}

type Action interface {
	Invoke(event *Payload) error
}

type Agent struct {
	Docker   docker.DockerClient
	Handlers []*EventHandler
	Channel  chan *docker.Event
	Errors   chan error
}

func New(client docker.DockerClient) *Agent {
	return &Agent{
		Docker:   client,
		Handlers: []*EventHandler{},
		Errors:   make(chan error),
	}
}

func (agent *Agent) Register(handler *EventHandler) {
	agent.Handlers = append(agent.Handlers, handler)
}

func (agent *Agent) Cancel(handler *EventHandler) *EventHandler {
	index := -1

	for i, h := range agent.Handlers {
		if h == handler {
			index = i
			break
		}
	}

	if index >= 0 {
		cancelled := agent.Handlers[index]
		agent.Handlers = append(agent.Handlers[:index], agent.Handlers[index+1:]...)
		return cancelled
	}

	return nil
}

func (agent *Agent) handleEvent(event *docker.Event) {
	matchingHandlers := []*EventHandler{}

	container, err := agent.Docker.InspectContainer(event.ContainerID)
	if err != nil {
		agent.Errors <- err
		return
	}

	for _, handler := range agent.Handlers {
		if handler.Matcher.Matches(container.Labels["microcosm.type"], event.Status) {
			matchingHandlers = append(matchingHandlers, handler)
		}
	}

	for _, handler := range matchingHandlers {
		payload := Payload{container, event, handler}
		handler.Action.Invoke(&payload)
	}
}

func (agent *Agent) Listen() {
	if agent.Channel != nil {
		agent.Errors <- fmt.Errorf("agent already listening for events")
	}

	agent.Channel = make(chan *docker.Event, 10)
	err := agent.Docker.AddEventListener(agent.Channel)
	agent.Errors <- err

	if err != nil {
		return
	}

	for {
		event, ok := <-agent.Channel
		if !ok {
			break
		}
		go agent.handleEvent(event)
	}
}
