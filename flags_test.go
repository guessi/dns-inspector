package main

import "testing"

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: Config{
				Server:  "8.8.8.8:53",
				UDPSize: 4096,
			},
			wantErr: false,
		},
		{
			name: "valid IPv6",
			cfg: Config{
				Server:  "[2001:4860:4860::8888]:53",
				UDPSize: 512,
			},
			wantErr: false,
		},
		{
			name: "valid hostname",
			cfg: Config{
				Server:  "dns.google:53",
				UDPSize: 1232,
			},
			wantErr: false,
		},
		{
			name: "missing port",
			cfg: Config{
				Server:  "8.8.8.8",
				UDPSize: 4096,
			},
			wantErr: true,
		},
		{
			name: "invalid server format",
			cfg: Config{
				Server:  "not-a-valid-address",
				UDPSize: 4096,
			},
			wantErr: true,
		},
		{
			name: "UDP size too large",
			cfg: Config{
				Server:  "8.8.8.8:53",
				UDPSize: 65536,
			},
			wantErr: true,
		},
		{
			name: "UDP size way too large",
			cfg: Config{
				Server:  "8.8.8.8:53",
				UDPSize: 100000,
			},
			wantErr: true,
		},
		{
			name: "UDP size at maximum",
			cfg: Config{
				Server:  "8.8.8.8:53",
				UDPSize: 65535,
			},
			wantErr: false,
		},
		{
			name: "UDP size zero",
			cfg: Config{
				Server:  "8.8.8.8:53",
				UDPSize: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
