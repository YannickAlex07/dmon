package inframon_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	inframon "github.com/yannickalex07/inframon/pkg"
)

func TestMonitor(t *testing.T) {
	// Arrange
	ctx := context.Background()

	notifications := []inframon.Notification{
		{
			Key:         "1",
			Title:       "Hello, World!",
			Description: "This is a test notification",
			Logs:        []string{"log1", "log2"},
			Links:       map[string]*url.URL{},
		},
		// should not be in the destination as the key already exists in the storage
		{
			Key:         "2",
			Title:       "Hello, World!",
			Description: "This is a test notification",
			Logs:        []string{"log1", "log2"},
			Links:       map[string]*url.URL{},
		},
	}

	state := FakeState{Storage: map[string]interface{}{
		"2": "exsits",
	}}

	source := FakeSource{Notifications: notifications}
	destination := FakeDestination{}

	mon := inframon.Monitor{
		Sources:      []inframon.Source{&source},
		Destinations: []inframon.Destination{&destination},
		State:        &state,
	}

	// Act
	err := mon.Start(ctx)
	if err != nil {
		assert.Fail(t, err.Error())
	}

	// Assert
	assert.True(t, len(destination.Notifications) == 1)
	assert.ElementsMatch(t, notifications[0:1], destination.Notifications)
}
