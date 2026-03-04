package main

import (
	"fmt"
	"net"
	"time"

	dnsv1 "github.com/miekg/dns"
)

// queryDNS performs a DNS query with specified parameters and prints the results.
func queryDNS(cfg *Config) {
	printHeader(fmt.Sprintf("DNS Query (UDPSize=%d)", cfg.UDPSize))

	// Parse query type
	var dnsType uint16
	if t, ok := dnsv1.StringToType[cfg.Type]; ok {
		dnsType = t
	} else {
		fmt.Printf("Invalid query type: %s\n", cfg.Type)
		return
	}

	msg := buildDNSMessage(cfg.Domain, uint16(cfg.UDPSize), dnsType)

	printRequest(cfg.Domain, cfg.Server, uint16(cfg.UDPSize))
	if cfg.Debug {
		printJSON("Request", msg)
	}

	// Execute DNS query with custom timeout
	client := &dnsv1.Client{Timeout: time.Duration(cfg.Timeout) * time.Second}
	response, rtt, err := client.Exchange(msg, cfg.Server)

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
func buildDNSMessage(domain string, udpSize uint16, qtype uint16) *dnsv1.Msg {
	msg := new(dnsv1.Msg)
	msg.SetQuestion(domain, qtype)
	// Add EDNS0 OPT record with client subnet information
	msg.Extra = append(msg.Extra, &dnsv1.OPT{
		Hdr: dnsv1.RR_Header{
			Name:   ".",
			Rrtype: dnsv1.TypeOPT,
			Class:  udpSize,  // UDP payload size
		},
		Option: []dnsv1.EDNS0{
			&dnsv1.EDNS0_SUBNET{
				Code:          dnsv1.EDNS0SUBNET,
				Family:        1,  // IPv4
				SourceNetmask: 0,  // No client subnet information
				Address:       net.ParseIP("0.0.0.0"),
			},
		},
	})
	return msg
}
