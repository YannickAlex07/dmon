package util_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/util"
)

func TestParseInvalidTimestamp(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	invalidTimestamp := "this-is-invalid"

	// Act
	res, err := util.ParseTimestamp(invalidTimestamp)

	// Assert
	assert.Equal(time.Time{}, res)
	assert.NotNil(err)
}

func TestParseValidTimestamp(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	ts := "2000-01-01T10:00:00Z"
	expectedTime := time.Date(2000, 01, 01, 10, 00, 00, 00, time.UTC)

	// Act
	res, err := util.ParseTimestamp(ts)

	// Assert
	assert.Equal(expectedTime, res)
	assert.Nil(err)
}
