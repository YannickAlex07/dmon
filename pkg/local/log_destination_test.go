package local_test

import (
	"context"
	"testing"

	testHook "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	inframon "github.com/yannickalex07/inframon/pkg"
	"github.com/yannickalex07/inframon/pkg/local"
)

func TestLogDestination(t *testing.T) {
	// Arrange
	ctx := context.Background()

	hook := testHook.NewGlobal()
	destination := local.LogDestination{}

	n := inframon.Notification{
		Key:         "test-notification",
		Title:       "Test Notification",
		Description: "This is a test notification",
	}

	// Act
	err := destination.Handle(ctx, n)
	if err != nil {
		assert.FailNow(t, "Failed to handle notification", err)
	}

	// Assert
	assert.Equal(t, 1, len(hook.AllEntries()))

	msg := hook.LastEntry().Message
	assert.Equal(t, "Received Notification: {Key:test-notification Title:Test Notification Description:This is a test notification Logs:[] Links:map[]}", msg)
}
