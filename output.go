package main

import (
	"encoding/json"
	"fmt"
	"time"

	dnsv1 "github.com/miekg/dns"
)

// printHeader prints a formatted section header.
func printHeader(title string) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println(title)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}

// printJSON marshals and prints the given value as JSON.
func printJSON(label string, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("%s JSON: Error marshaling: %v\n\n", label, err)
		return
	}
	fmt.Printf("%s JSON:\n%s\n\n", label, data)
}

// printRequest prints DNS request details.
func printRequest(domain, server string, udpSize uint16) {
	fmt.Printf("📤 Request:\n")
	fmt.Printf("   Domain: %s\n", domain)
	fmt.Printf("   Server: %s\n", server)
	fmt.Printf("   EDNS0 UDP Size: %d bytes\n\n", udpSize)
}

// printError prints DNS query error information.
func printError(err error) {
	fmt.Printf("📥 Response: ❌ ERROR\n")
	fmt.Printf("   %v\n\n", err)
}

// printResponse prints DNS response details including answers and round-trip time.
func printResponse(r *dnsv1.Msg, rtt time.Duration) {
	truncated := ""
	if r.Truncated {
		truncated = " ⚠️  TRUNCATED"
	}
	fmt.Printf("📥 Response: ✅ SUCCESS (RTT: %v)%s\n", rtt, truncated)
	fmt.Printf("   Response Code: %s\n", dnsv1.RcodeToString[r.Rcode])
	fmt.Printf("   Message Size: %d bytes\n", r.Len())
	fmt.Printf("   Answers: %d, Authority: %d, Additional: %d\n", len(r.Answer), len(r.Ns), len(r.Extra))
	
	// Check for EDNS0 in response
	if opt := r.IsEdns0(); opt != nil {
		fmt.Printf("   Server EDNS0 UDP Size: %d bytes\n", opt.UDPSize())
	}
	
	if len(r.Answer) > 0 {
		fmt.Println("   Answer Records:")
		for _, ans := range r.Answer {
			fmt.Printf("     • %s\n", ans.String())
		}
	}
	fmt.Println()
}
