package main

import ()

type EventHandler struct {
	Matcher *EventMatcher
	Command *Command
}

type Command struct {
	Container string
	Action    string
}

type Agent struct {
	Handlers []*EventHandler
}

func New() *Agent {
	return &Agent{
		Handlers: []*EventHandler{},
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
