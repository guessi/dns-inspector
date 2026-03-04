# DNS Inspector

[![Go Version](https://img.shields.io/github/go-mod/go-version/guessi/dns-inspector)](https://github.com/guessi/dns-inspector/blob/main/go.mod "Go module version")
[![License](https://img.shields.io/github/license/guessi/dns-inspector)](https://github.com/guessi/dns-inspector/blob/main/LICENSE "MIT License")
[![Go Report Card](https://goreportcard.com/badge/github.com/guessi/dns-inspector)](https://goreportcard.com/report/github.com/guessi/dns-inspector "Go Report Card for dns-inspector")

A Go CLI tool for debugging DNS queries with different EDNS0 UDP buffer sizes. Diagnoses truncation issues and helps catch silent failures where responses break under smaller buffer configurations.

## Why This Tool?

Read the full story: [Why DNS Breaks in Production but Works Locally](https://guessi.github.io/posts/2026/why-dns-breaks-in-production-but-works-locally/)

## Quick Start

```bash
go install github.com/guessi/dns-inspector@latest
dns-inspector -domain google.com.
```

## Usage

![dns-inspector usage and available flags](assets/usage.svg)

## Examples

With a small UDP buffer size (200 bytes), the response is truncated — zero answers returned:

![dns-inspector demo showing truncated DNS response](assets/demo-truncated.svg)

Increasing the buffer to 4096 bytes returns the full response with all 13 answers:

![dns-inspector demo showing full DNS response](assets/demo-full.svg)

### When UDP Is Not Enough

Some domains have records too large for UDP, even at the maximum buffer size. The TC (truncation) bit signals that clients should retry over TCP:

![dns-inspector showing truncation even at max UDP buffer size](assets/udp-limit.svg)

## Background

DNS responses that exceed the UDP buffer size get truncated — the server sets the TC bit and returns no answers, signaling the client to retry over TCP. Common buffer sizes:

- **512 bytes** — original limit ([RFC 1035](https://tools.ietf.org/html/rfc1035))
- **1232 bytes** — recommended minimum ([DNS Flag Day 2020](https://www.dnsflagday.net/2020/))
- **4096 bytes** — EDNS0 default ([RFC 6891](https://tools.ietf.org/html/rfc6891))

## Dependencies

- [miekg/dns](https://github.com/miekg/dns) — DNS library for Go

## License

[MIT](LICENSE)
