package storage_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/pkg/storage"
)

func TestMemoryStoreGetAndSetExecutionTime(t *testing.T) {
	// - Arrange
	storage := storage.NewMemoryStore(1 * time.Hour)

	// we append some time to make sure no time.Now() can be accidentally correct
	lastExecutionTime := time.Now().Add(1 * time.Hour)

	// - Act
	storeError := storage.SetLatestExecutionTime(lastExecutionTime)
	fetchedTime, fetchError := storage.GetLatestExecutionTime()

	// - Assert
	assert.Nil(t, storeError)
	assert.Nil(t, fetchError)
	assert.True(t, lastExecutionTime.Equal(fetchedTime))
}

func TestMemoryStoreIsTimeoutStoredWithNonStoredTimeout(t *testing.T) {
	// - Arrange
	storage := storage.NewMemoryStore(1 * time.Hour)

	// - Act
	handled, err := storage.IsTimeoutStored("job-id")

	// - Assert
	assert.Nil(t, err)
	assert.False(t, handled)
}

func TestMemoryStoreIsTimeoutStoredWithStoredTimeout(t *testing.T) {
	// - Arrange
	storage := storage.NewMemoryStore(1 * time.Hour)

	// - Act
	storeError := storage.StoreTimeout("job-id", time.Now())
	handled, fetchError := storage.IsTimeoutStored("job-id")

	// - Assert
	assert.Nil(t, storeError)
	assert.Nil(t, fetchError)
	assert.True(t, handled)
}

func TestMemoryStoreIsTimeoutStoredWithExpiredTimeout(t *testing.T) {
	// - Arrange
	// every timeout should expire after 1 second
	storage := storage.NewMemoryStore(1 * time.Second)

	// - Act
	storeError := storage.StoreTimeout("job-id", time.Now())

	time.Sleep(2 * time.Second)

	handled, fetchError := storage.IsTimeoutStored("job-id")

	// - Assert
	assert.Nil(t, storeError)
	assert.Nil(t, fetchError)
	assert.False(t, handled)
}
