package local_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/pkg/local"
)

func TestMemoryStoreStoringWithoutExpire(t *testing.T) {
	// Arrange
	ctx := context.Background()
	storage := local.NewMemoryStorage(10 * time.Second)

	// Act
	err := storage.Store(ctx, "key", "value", false)
	if err != nil {
		t.Errorf("Error storing value: %v", err)
	}

	// Assert
	value, err := storage.Get(ctx, "key")
	if err != nil {
		t.Errorf("Error getting value: %v", err)
	}

	valueStr, ok := value.(string)
	if !ok {
		t.Errorf("Failed to cast value to string")
	}

	assert.Equal(t, "value", valueStr)
}

func TestMemoryStoreStoringWithExpire(t *testing.T) {
	// Arrange
	ctx := context.Background()
	storage := local.NewMemoryStorage(1 * time.Second)

	// Act
	err := storage.Store(ctx, "key", "value", true)
	if err != nil {
		t.Errorf("Error storing value: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Assert
	exists, err := storage.Exists(ctx, "key")
	if err != nil {
		t.Errorf("Error checking for value: %v", err)
	}

	assert.False(t, exists)
}

func TestMemoryStoreExistsWithExistingKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	storage := local.NewMemoryStorage(1 * time.Second)

	err := storage.Store(ctx, "key", "value", false)
	if err != nil {
		t.Errorf("Error storing value: %v", err)
	}

	// Act
	exists, err := storage.Exists(ctx, "key")
	if err != nil {
		t.Errorf("Error checking for value: %v", err)
	}

	// Assert
	assert.True(t, exists)
}

func TestMemoryStoreExistsWithNonExistingKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	storage := local.NewMemoryStorage(1 * time.Second)

	// Act
	exists, err := storage.Exists(ctx, "key")
	if err != nil {
		t.Errorf("Error checking for value: %v", err)
	}

	// Assert
	assert.False(t, exists)
}

func TestMemoryStoreGetWithExistingKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	storage := local.NewMemoryStorage(1 * time.Second)

	err := storage.Store(ctx, "key", "value", false)
	if err != nil {
		t.Errorf("Error storing value: %v", err)
	}

	// Act
	value, err := storage.Get(ctx, "key")
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

func TestMemoryStoreGetWithNonExistingKey(t *testing.T) {
	// Arrange
	ctx := context.Background()
	storage := local.NewMemoryStorage(1 * time.Second)

	// Act
	value, err := storage.Get(ctx, "key")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, value)
}
