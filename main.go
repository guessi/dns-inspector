// Package main provides a DNS debugging tool for testing DNS queries with different EDNS0 UDP buffer sizes.
package main

import "fmt"

func main() {
	// Parse command-line flags
	cfg := parseFlags()

	// Show version and exit
	if cfg.ShowVersion {
		fmt.Printf("dns-inspector version %s\n", version)
		return
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Println(err)
		return
	}

	queryDNS(cfg)
}
