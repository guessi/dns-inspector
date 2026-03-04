package main

import (
	"flag"
	"fmt"
	"net"
)

const version = "0.1.0"

// Config holds all command-line flag values
type Config struct {
	Domain      string
	Server      string
	UDPSize     uint
	Type        string
	Timeout     uint
	Debug       bool
	ShowVersion bool
}

// parseFlags parses command-line flags and returns a Config
func parseFlags() *Config {
	cfg := &Config{}
	
	flag.StringVar(&cfg.Domain, "domain", "api.github.com.", "Domain name to query (must end with .)")
	flag.StringVar(&cfg.Server, "server", "1.1.1.1:53", "DNS server address (host:port)")
	flag.UintVar(&cfg.UDPSize, "udpsize", 4096, "EDNS0 UDP payload size (0-65535, use 0 to test truncation)")
	flag.StringVar(&cfg.Type, "type", "A", "Query type (A, AAAA, TXT, MX, etc.)")
	flag.UintVar(&cfg.Timeout, "timeout", 3, "Query timeout in seconds (default: 3)")
	flag.BoolVar(&cfg.Debug, "debug", false, "Print full request/response JSON")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "Show version information")
	flag.Parse()
	
	return cfg
}

// Validate checks if the configuration is valid and returns an error if not
func (c *Config) Validate() error {
	// Validate server address format
	if _, _, err := net.SplitHostPort(c.Server); err != nil {
		return fmt.Errorf("invalid server address format: %v", err)
	}

	// Validate UDP size range
	if c.UDPSize > 65535 {
		return fmt.Errorf("invalid UDP size: %d (must be 0-65535)", c.UDPSize)
	}

	return nil
}
