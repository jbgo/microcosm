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
