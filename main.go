package main

import (
	"flag"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Project struct {
		ID       string `yaml:"id"`
		Location string `yaml:"location"`
	} `yaml:"project"`
}

func parseConfig(path string) *Config {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic("Failed to read config file with err " + err.Error())
	}

	var c *Config
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		panic("Failed to read config file with err " + err.Error())
	}

	return c
}

func main() {
	// Parse CLI Arguments
	configPath := flag.String("config", "config.yaml", "Path to the config file")
	flag.Parse()

	// Parse Config
	cfg := parseConfig(*configPath)

	// Start Monitor
	startMonitor(*cfg)
}
