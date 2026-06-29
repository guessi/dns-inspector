package main

import (
	"net/netip"
	"testing"

	dnsv2 "codeberg.org/miekg/dns"
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
			qtype:    dnsv2.TypeA,
			wantType: dnsv2.TypeA,
		},
		{
			name:     "TXT record query",
			domain:   "google.com.",
			udpSize:  512,
			qtype:    dnsv2.TypeTXT,
			wantType: dnsv2.TypeTXT,
		},
		{
			name:     "Zero UDP size",
			domain:   "test.com.",
			udpSize:  0,
			qtype:    dnsv2.TypeA,
			wantType: dnsv2.TypeA,
		},
		{
			name:     "MX record query",
			domain:   "example.com.",
			udpSize:  1232,
			qtype:    dnsv2.TypeMX,
			wantType: dnsv2.TypeMX,
		},
		{
			name:     "AAAA record query",
			domain:   "example.com.",
			udpSize:  4096,
			qtype:    dnsv2.TypeAAAA,
			wantType: dnsv2.TypeAAAA,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := buildDNSMessage(tt.domain, tt.udpSize, tt.qtype)

			// Check question section
			if len(msg.Question) != 1 {
				t.Errorf("expected 1 question, got %d", len(msg.Question))
			}

			if msg.Question[0].Header().Name != tt.domain {
				t.Errorf("expected domain %s, got %s", tt.domain, msg.Question[0].Header().Name)
			}

			if dnsv2.RRToType(msg.Question[0]) != tt.wantType {
				t.Errorf("expected qtype %d, got %d", tt.wantType, dnsv2.RRToType(msg.Question[0]))
			}

			// Check EDNS0 UDP size
			if msg.UDPSize != tt.udpSize {
				t.Errorf("expected UDP size %d, got %d", tt.udpSize, msg.UDPSize)
			}

			// Check SUBNET option is present regardless of UDP size
			hasSubnet := false
			for _, opt := range msg.Pseudo {
				if _, ok := opt.(*dnsv2.SUBNET); ok {
					hasSubnet = true
					break
				}
			}
			if !hasSubnet {
				t.Error("expected SUBNET option in pseudo section")
			}
		})
	}
}

func TestBuildDNSMessageEDNS0(t *testing.T) {
	msg := buildDNSMessage("example.com.", 4096, dnsv2.TypeA)

	// Verify EDNS0 UDP size is set
	if msg.UDPSize != 4096 {
		t.Errorf("expected UDPSize 4096, got %d", msg.UDPSize)
	}

	// Verify ECS option is present in pseudo section
	if len(msg.Pseudo) == 0 {
		t.Fatal("expected pseudo section options, got none")
	}

	// Check for SUBNET option
	hasECS := false
	for _, option := range msg.Pseudo {
		if subnet, ok := option.(*dnsv2.SUBNET); ok {
			hasECS = true
			if subnet.Family != 1 {
				t.Errorf("expected Family 1 (IPv4), got %d", subnet.Family)
			}
			if subnet.Netmask != 0 {
				t.Errorf("expected Netmask 0, got %d", subnet.Netmask)
			}
			if subnet.Address != netip.MustParseAddr("0.0.0.0") {
				t.Errorf("expected Address 0.0.0.0, got %s", subnet.Address)
			}
			break
		}
	}

	if !hasECS {
		t.Error("expected EDNS0 Client Subnet option")
	}
}

func TestBuildDNSMessagePackZeroUDPSize(t *testing.T) {
	msg := buildDNSMessage("example.com.", 0, dnsv2.TypeA)

	if err := msg.Pack(); err != nil {
		t.Fatalf("Pack() with UDPSize=0 failed: %v", err)
	}
	if len(msg.Data) == 0 {
		t.Fatal("Pack() produced empty buffer")
	}

	// Verify SUBNET option survives the round-trip
	var msg2 dnsv2.Msg
	msg2.Data = msg.Data
	if err := msg2.Unpack(); err != nil {
		t.Fatalf("Unpack() failed: %v", err)
	}

	hasSubnet := false
	for _, opt := range msg2.Pseudo {
		if _, ok := opt.(*dnsv2.SUBNET); ok {
			hasSubnet = true
			break
		}
	}
	if !hasSubnet {
		t.Error("SUBNET option lost after Pack/Unpack with UDPSize=0")
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
			wantValid: false, // dnsv2.StringToType is case-sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := dnsv2.StringToType[tt.qtype]
			if ok != tt.wantValid {
				t.Errorf("query type %q: got valid=%v, want valid=%v", tt.qtype, ok, tt.wantValid)
			}
		})
	}
}
