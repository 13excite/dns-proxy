// Package config contains settings and basic logger configuration
package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// Config is the main config of the service
type Config struct {
	Addr        string         `yaml:"address"`
	Port        int            `yaml:"port"`
	MetricsPort int            `yaml:"metrics_port"`
	LogLevel    string         `yaml:"log_level"`
	LogEncoding string         `yaml:"log_encoding"`
	Upstream    UpstreamConfig `yaml:"upstream"`
}

// UpstreamConfig contains info about upstream DNS over TLS resolver
type UpstreamConfig struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Hostname string `yaml:"hostname"`
	Timeout  int    `yaml:"timeout"` // in seconds
}

// Defaults sets default values for the configuration
func (conf *Config) Defaults() {
	conf.Addr = "0.0.0.0"
	conf.Port = 1153
	conf.MetricsPort = 8090
	conf.LogLevel = "info"
	conf.LogEncoding = "console"
	conf.Upstream = UpstreamConfig{
		Address:  "1.1.1.1",
		Port:     853,
		Hostname: "one.one.one.one",
		Timeout:  2,
	}
}

// ReadConfigFile reading and parsing configuration yaml file
func (conf *Config) ReadConfigFile(configPath string) {
	yamlConfig, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlConfig, &conf)
	if err != nil {
		log.Fatal(fmt.Errorf("could not unmarshal config %v", conf), err)
	}
}
