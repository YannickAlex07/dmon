package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yannickalex07/dmon/config"
)

func TestReadConfigFromInvalidPath(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	path := "/this/is/not/valid"

	// Act
	cfg, err := config.Read(path)

	// Assert
	assert.Nil(cfg)
	assert.NotNil(err)
}

func TestReadInvalidConfig(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	path := "./test/invalid.yml"

	// Act
	cfg, err := config.Read(path)

	// Assert
	assert.Nil(cfg)
	assert.NotNil(err)
}

func TestReadValidConfig(t *testing.T) {
	assert := assert.New(t)

	// Arrange
	path := "./test/valid.yml"

	// Act
	cfg, err := config.Read(path)

	// Assert
	assert.Nil(cfg)
	assert.Nil(err)
}
