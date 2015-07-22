package main

import (
	"testing"
)

func TestNewEventMatcher(t *testing.T) {
	matcher := NewEventMatcher("web", "start", "stop", "restart")
	if matcher.ServiceType != "web" {
		t.Error("expected ServiceType", "web", "got", matcher.ServiceType)
	}

	events := []string{"start", "stop", "restart"}
	for _, e := range events {
		if !matcher.Events[e] {
			t.Errorf("expected to find %s event", e)
		}
	}
}

func TestEventMatcherMatches(t *testing.T) {
	matcher := NewEventMatcher("web", "start", "stop", "restart")

	expectMatch := []map[string]string{
		map[string]string{"service": "web", "status": "start"},
		map[string]string{"service": "web", "status": "stop"},
		map[string]string{"service": "web", "status": "restart"},
	}

	for _, event := range expectMatch {
		if !matcher.Matches(event["service"], event["status"]) {
			t.Errorf("expected match for service: %s, status: %s", event["service"], event["status"])
		}
	}

	expectNoMatch := []map[string]string{
		map[string]string{"service": "web", "status": "die"},
		map[string]string{"service": "db", "status": "start"},
	}

	for _, event := range expectNoMatch {
		if matcher.Matches(event["service"], event["status"]) {
			t.Errorf("not expecting match for service: %s, status: %s", event["service"], event["status"])
		}
	}
}
