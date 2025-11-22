package events

import (
	"time"
)

// ClickEvent represents the data for a click event
type ClickEvent struct {
	ShortCode string
	IPAddress string
	ClickedAt time.Time
}

// ClickQueue is a buffered channel to hold click events
var ClickQueue = make(chan ClickEvent, 1000)

// PublishClickEvent sends a click event to the queue
func PublishClickEvent(event ClickEvent) {
	select {
	case ClickQueue <- event:
		// Event sent successfully
	default:
		// Queue is full, drop event or log warning
		// In a production system, we might want better handling here
	}
}
