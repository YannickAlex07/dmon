package local_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/inframon/pkg/local"
)

func TestMemoryStateStoringWithoutExpire(t *testing.T) {
	// Arrange
	ctx := context.Background()
	state := local.NewMemoryState(10 * time.Second)

	// Act
	err := state.Store(ctx, "key", "value", false)
	if err != nil {
		t.Errorf("Error storing value: %v", err)
	}

	// Assert
	value, err := state.Get(ctx, "key")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}

	valueStr, ok := value.(string)
	if !ok {
		t.Errorf("Failed to cast value to string")
	}

	assert.Equal(t, "value", valueStr)
}

func TestMemoryStateStoringWithExpire(t *testing.T) {
	// Arrange
	ctx := context.Background()
	state := local.NewMemoryState(1 * time.Second)

	// Act
	err := state.Store(ctx, "key", "value", true)
	if err != nil {
		t.Errorf("Error storing value: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Assert
	exists, err := state.Exists(ctx, "key")
	if err != nil {
		t.Errorf("Error checking for value: %v", err)
	}

	assert.False(t, exists)
}

func TestMemoryStateExistsWithExistingKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	state := local.NewMemoryState(1 * time.Second)

	err := state.Store(ctx, "key", "value", false)
	if err != nil {
		t.Errorf("Error storing value: %v", err)
	}

	// Act
	exists, err := state.Exists(ctx, "key")
	if err != nil {
		t.Errorf("Error checking for value: %v", err)
	}

	// Assert
	assert.True(t, exists)
}

func TestMemoryStateExistsWithNonExistingKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	state := local.NewMemoryState(1 * time.Second)

	// Act
	exists, err := state.Exists(ctx, "key")
	if err != nil {
		t.Errorf("Error checking for value: %v", err)
	}

	// Assert
	assert.False(t, exists)
}

func TestMemoryStateGetWithExistingKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	state := local.NewMemoryState(1 * time.Second)

	err := state.Store(ctx, "key", "value", false)
	if err != nil {
		t.Errorf("Error storing value: %v", err)
	}

	// Act
	value, err := state.Get(ctx, "key")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}

	// Assert
	valueStr, ok := value.(string)
	if !ok {
		t.Errorf("Failed to cast value to string")
	}

	assert.Equal(t, "value", valueStr)
}

func TestMemoryStateGetWithNonExistingKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	state := local.NewMemoryState(1 * time.Second)

	// Act
	value, err := state.Get(ctx, "key")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, value)
}
