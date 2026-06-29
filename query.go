package main

import (
	"context"
	"fmt"
	"net/netip"
	"time"

	dnsv2 "codeberg.org/miekg/dns"
)

// queryDNS performs a DNS query with specified parameters and prints the results.
func queryDNS(cfg *Config) {
	printHeader(fmt.Sprintf("DNS Query (UDPSize=%d)", cfg.UDPSize))

	// Parse query type
	dnsType, ok := dnsv2.StringToType[cfg.Type]
	if !ok {
		fmt.Printf("Invalid query type: %s\n", cfg.Type)
		return
	}

	msg := buildDNSMessage(cfg.Domain, uint16(cfg.UDPSize), dnsType)

	printRequest(cfg.Domain, cfg.Server, uint16(cfg.UDPSize))
	if cfg.Debug {
		printJSON("Request", msg)
	}

	// Execute DNS query with custom timeout
	client := dnsv2.NewClient()
	client.Transport.ReadTimeout = time.Duration(cfg.Timeout) * time.Second
	client.Transport.WriteTimeout = time.Duration(cfg.Timeout) * time.Second
	response, rtt, err := client.Exchange(context.Background(), msg, "udp", cfg.Server)

	if err != nil {
		printError(err)
		return
	}

	printResponse(response, rtt)
	if cfg.Debug {
		printJSON("Response", response)
	}
}

// buildDNSMessage creates a DNS query message with EDNS0 client subnet option.
func buildDNSMessage(domain string, udpSize uint16, qtype uint16) *dnsv2.Msg {
	msg := dnsv2.NewMsg(domain, qtype)
	msg.UDPSize = udpSize
	// Add EDNS0 client subnet option
	msg.Pseudo = append(msg.Pseudo, &dnsv2.SUBNET{
		Family:  1, // IPv4
		Netmask: 0, // No client subnet information
		Address: netip.MustParseAddr("0.0.0.0"),
	})
	return msg
}
