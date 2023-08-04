package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/pkg/util"
)

func TestParseTimestampWithValidTimestamp(t *testing.T) {
	// - Arrange
	now := time.Now().UTC().Round(time.Second)
	nowStr := now.Format(time.RFC3339)

	// - Act
	parsed, err := util.ParseTimestamp(nowStr)

	if err != nil {
		t.Error(err)
	}

	// - Assert
	assert.True(t, now.Equal(parsed))
}

func TestParseTimestampWithWrongFormat(t *testing.T) {
	// - Arrange
	now := time.Now().UTC().Round(time.Second)
	nowStr := now.Format(time.RFC850) // different timestamp format

	// - Act
	result, err := util.ParseTimestamp(nowStr)

	// - Assert
	assert.True(t, result.IsZero())
	assert.Error(t, err)
}
