package keiho

import "net/url"

// A notification that needs to be handled by a handler
type Notification struct {
	// The title of the notification
	Title string

	// The description of what the notification is about
	Description string

	// Logs associated with the notification
	Logs []string

	// A map of string -> urls links
	Links map[string]*url.URL
}
