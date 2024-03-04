// Package config contains settings and basic logger configuration
package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Host        string         `yaml:"host"`
	TCPPort     int            `yaml:"tcp_port"`
	UDPPort     int            `yaml:"udp_port"`
	MetricsPort int            `yaml:"metrics_port"`
	LogLevel    string         `yaml:"log_level"`
	LogEncoding string         `yaml:"log_encoding"`
	Upstream    UpstreamConfig `yaml:"upstream"`
}

type UpstreamConfig struct {
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	Hostname string `yaml:"hostname"`
	Timeout  int    `yaml:"timeout"` // in seconds
}

// Defaults sets default values for the configuration
func (conf *Config) Defaults() {
	conf.Host = "127.0.0.1"
	conf.TCPPort = 1153
	conf.UDPPort = 1153
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
