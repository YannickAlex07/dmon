package dataflow_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/internal/gcp/dataflow"
)

func TestMessageLevelFromString(t *testing.T) {
	// Arrange
	values := map[string]dataflow.MessageLevel{
		"JOB_MESSAGE_IMPORTANCE_UNKNOWN": dataflow.LEVEL_UNKNOWN,
		"JOB_MESSAGE_DEBUG":              dataflow.LEVEL_DEBUG,
		"JOB_MESSAGE_DETAILED":           dataflow.LEVEL_DETAILED,
		"JOB_MESSAGE_BASIC":              dataflow.LEVEL_BASIC,
		"JOB_MESSAGE_WARNING":            dataflow.LEVEL_WARNING,
		"JOB_MESSAGE_ERROR":              dataflow.LEVEL_ERROR,
	}

	// Act
	for str, value := range values {
		actual := dataflow.MessageLevelFromString(str)

		// Assert
		assert.Equal(t, value, actual)
	}
}

func TestMessageLevelFromStringInvalidValue(t *testing.T) {
	// Act
	actual := dataflow.MessageLevelFromString("random")

	// Assert
	assert.Equal(t, dataflow.LEVEL_UNKNOWN, actual)
}
