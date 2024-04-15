package local_test

import (
	"context"
	"testing"

	log "github.com/sirupsen/logrus"
	testHook "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	inframon "github.com/yannickalex07/inframon/pkg"
	"github.com/yannickalex07/inframon/pkg/local"
)

func TestLogDestinationWithInfoLevel(t *testing.T) {
	// Arrange
	ctx := context.Background()

	hook := testHook.NewGlobal()
	destination, err := local.NewLogDestination("info")
	if err != nil {
		assert.FailNow(t, "Failed to create log destination", err)
	}

	n := inframon.Notification{
		Key:         "test-notification",
		Title:       "Test Notification",
		Description: "This is a test notification",
	}

	// Act
	err = destination.Handle(ctx, n)
	if err != nil {
		assert.FailNow(t, "Failed to handle notification", err)
	}

	// Assert
	assert.Equal(t, 1, len(hook.AllEntries()))
	assert.Equal(t, log.InfoLevel, hook.LastEntry().Level)

	msg := hook.LastEntry().Message
	assert.Equal(t, "Received Notification: {Key:test-notification Title:Test Notification Description:This is a test notification Logs:[] Links:map[]}", msg)
}

func TestLogDestinationWithNonExistingLevel(t *testing.T) {
	// Act
	_, err := local.NewLogDestination("test")

	// Assert
	assert.Error(t, err)
	assert.ErrorContains(t, err, "not a valid logrus Level")
}

func TestNewLogDestinationsWithAllLevels(t *testing.T) {
	// Arrange
	lvlStr := []string{
		"trace",
		"debug",
		"info",
		"warn",
		"error",
		"fatal",
		"panic",
	}

	// Act
	for _, l := range lvlStr {
		_, err := local.NewLogDestination(l)

		// Assert
		assert.NoError(t, err)
	}
}
