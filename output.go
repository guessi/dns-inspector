package main

import (
	"encoding/json"
	"fmt"
	"time"

	dnsv2 "codeberg.org/miekg/dns"
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
func printResponse(r *dnsv2.Msg, rtt time.Duration) {
	truncated := ""
	if r.Truncated {
		truncated = " ⚠️  TRUNCATED"
	}
	fmt.Printf("📥 Response: ✅ SUCCESS (RTT: %v)%s\n", rtt, truncated)
	fmt.Printf("   Response Code: %s\n", dnsv2.RcodeToString[r.Rcode])
	fmt.Printf("   Message Size: %d bytes\n", r.Len())
	// ARCOUNT is the DNS header field counting records in the Additional
	// section. When EDNS0 is used the server returns an OPT record there, so
	// ARCOUNT includes it. dnsv2 folds that OPT into r.UDPSize instead of
	// keeping it in r.Extra, so add it back to match ARCOUNT (what dig shows).
	additional := len(r.Extra)
	if r.UDPSize > 0 {
		additional++
	}
	fmt.Printf("   Answers: %d, Authority: %d, Additional: %d\n", len(r.Answer), len(r.Ns), additional)

	// Check for EDNS0 in response
	if r.UDPSize > 0 {
		fmt.Printf("   Server EDNS0 UDP Size: %d bytes\n", r.UDPSize)
	}

	if len(r.Answer) > 0 {
		fmt.Println("   Answer Records:")
		for _, ans := range r.Answer {
			fmt.Printf("     • %s\n", ans.String())
		}
	}
	fmt.Println()
}
