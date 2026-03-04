package main

import (
	"testing"

	dnsv1 "github.com/miekg/dns"
)

func TestBuildDNSMessage(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		udpSize  uint16
		qtype    uint16
		wantType uint16
	}{
		{
			name:     "A record query",
			domain:   "example.com.",
			udpSize:  4096,
			qtype:    dnsv1.TypeA,
			wantType: dnsv1.TypeA,
		},
		{
			name:     "TXT record query",
			domain:   "google.com.",
			udpSize:  512,
			qtype:    dnsv1.TypeTXT,
			wantType: dnsv1.TypeTXT,
		},
		{
			name:     "Zero UDP size",
			domain:   "test.com.",
			udpSize:  0,
			qtype:    dnsv1.TypeA,
			wantType: dnsv1.TypeA,
		},
		{
			name:     "MX record query",
			domain:   "example.com.",
			udpSize:  1232,
			qtype:    dnsv1.TypeMX,
			wantType: dnsv1.TypeMX,
		},
		{
			name:     "AAAA record query",
			domain:   "example.com.",
			udpSize:  4096,
			qtype:    dnsv1.TypeAAAA,
			wantType: dnsv1.TypeAAAA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := buildDNSMessage(tt.domain, tt.udpSize, tt.qtype)

			// Check question section
			if len(msg.Question) != 1 {
				t.Errorf("expected 1 question, got %d", len(msg.Question))
			}

			if msg.Question[0].Name != tt.domain {
				t.Errorf("expected domain %s, got %s", tt.domain, msg.Question[0].Name)
			}

			if msg.Question[0].Qtype != tt.wantType {
				t.Errorf("expected qtype %d, got %d", tt.wantType, msg.Question[0].Qtype)
			}

			// Check EDNS0 OPT record
			opt := msg.IsEdns0()
			if opt == nil {
				t.Fatal("expected EDNS0 OPT record, got nil")
			}

			if opt.UDPSize() != tt.udpSize {
				t.Errorf("expected UDP size %d, got %d", tt.udpSize, opt.UDPSize())
			}
		})
	}
}

func TestBuildDNSMessageEDNS0(t *testing.T) {
	msg := buildDNSMessage("example.com.", 4096, dnsv1.TypeA)

	// Verify EDNS0 is present
	opt := msg.IsEdns0()
	if opt == nil {
		t.Fatal("EDNS0 OPT record should be present")
	}

	// Verify ECS option is present
	if len(opt.Option) == 0 {
		t.Fatal("expected EDNS0 options, got none")
	}

	// Check for ECS option
	hasECS := false
	for _, option := range opt.Option {
		if option.Option() == dnsv1.EDNS0SUBNET {
			hasECS = true
			break
		}
	}

	if !hasECS {
		t.Error("expected EDNS0 Client Subnet option")
	}
}

func TestQueryTypeValidation(t *testing.T) {
	tests := []struct {
		name      string
		qtype     string
		wantValid bool
	}{
		{
			name:      "valid A record",
			qtype:     "A",
			wantValid: true,
		},
		{
			name:      "valid AAAA record",
			qtype:     "AAAA",
			wantValid: true,
		},
		{
			name:      "valid TXT record",
			qtype:     "TXT",
			wantValid: true,
		},
		{
			name:      "valid MX record",
			qtype:     "MX",
			wantValid: true,
		},
		{
			name:      "valid NS record",
			qtype:     "NS",
			wantValid: true,
		},
		{
			name:      "valid CNAME record",
			qtype:     "CNAME",
			wantValid: true,
		},
		{
			name:      "invalid type",
			qtype:     "INVALID",
			wantValid: false,
		},
		{
			name:      "empty string",
			qtype:     "",
			wantValid: false,
		},
		{
			name:      "lowercase valid",
			qtype:     "a",
			wantValid: false, // dnsv1.StringToType is case-sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := dnsv1.StringToType[tt.qtype]
			if ok != tt.wantValid {
				t.Errorf("query type %q: got valid=%v, want valid=%v", tt.qtype, ok, tt.wantValid)
			}
		})
	}
}
