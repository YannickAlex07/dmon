package gmon

import "net/url"

// A notification that needs to be handled by a handler
type Notification struct {
	// The title of the notification
	Title string

	// An overview over what the notification is about
	Overview string

	// Logs associated with the notification
	Logs []string

	// A map of string -> urls links
	Links map[string]*url.URL
}
