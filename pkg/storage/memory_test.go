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
	storage.SetLatestExecutionTime(lastExecutionTime)
	fetchedTime, err := storage.GetLatestExecutionTime()

	// - Assert
	assert.Nil(t, err)
	assert.True(t, lastExecutionTime.Equal(fetchedTime))
}

func TestMemoryStoreWasTimeoutHandledWithNonHandledTimeout(t *testing.T) {
	// - Arrange
	storage := storage.NewMemoryStore(1 * time.Hour)

	// - Act
	handled := storage.WasTimeoutHandled("job-id")

	// - Assert
	assert.False(t, handled)
}

func TestMemoryStoreWasTimeoutHandledWithHandledTimeout(t *testing.T) {
	// - Arrange
	storage := storage.NewMemoryStore(1 * time.Hour)

	// - Act
	storage.HandleTimeout("job-id", time.Now())
	handled := storage.WasTimeoutHandled("job-id")

	// - Assert
	assert.True(t, handled)
}

func TestMemoryStoreWasTimeoutHandledWithExpiredTimeout(t *testing.T) {
	// - Arrange
	// every timeout should expire after 1 second
	storage := storage.NewMemoryStore(1 * time.Second)

	// - Act
	storage.HandleTimeout("job-id", time.Now())

	time.Sleep(2 * time.Second)

	handled := storage.WasTimeoutHandled("job-id")

	// - Assert
	assert.False(t, handled)
}
