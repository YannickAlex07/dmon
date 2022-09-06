package main

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Project struct {
		ID       string `yaml:"id"`
		Location string `yaml:"location"`
	} `yaml:"project"`
	Slack struct {
		Token   string `yaml:"token"`
		Channel string `yaml:"channel"`
	} `yaml:"slack"`
}

func parseConfig(path string) *Config {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		panic("Failed to read config file with err " + err.Error())
	}

	var c *Config
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		panic("Failed to read config file with err " + err.Error())
	}

	return c
}

func main() {
	// Parse CLI Arguments
	configPath := flag.String("c", "./config.yaml", "Path to the config file")
	flag.Parse()

	// Parse Config
	cfg := parseConfig(*configPath)

	// Start Monitor
	startMonitor(*cfg)
}
