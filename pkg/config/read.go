package config

import (
	"os"

	"github.com/yannickalex07/dmon/pkg/models"
	"gopkg.in/yaml.v3"
)

func Read(path string) (*models.Config, error) {
	// open file
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// unmarshal config
	var c *models.Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
