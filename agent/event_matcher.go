package main

type EventMatcher struct {
	Events      map[string]bool
	ServiceType string
}

func NewEventMatcher(serviceType string, eventNames ...string) *EventMatcher {
	matcher := &EventMatcher{Events: map[string]bool{}}
	matcher.ServiceType = serviceType
	for _, e := range eventNames {
		matcher.Events[e] = true
	}
	return matcher
}

func (matcher *EventMatcher) Matches(serviceType, eventStatus string) bool {
	return serviceType == matcher.ServiceType && matcher.Events[eventStatus]
}
