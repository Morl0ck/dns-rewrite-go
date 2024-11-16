package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/miekg/dns"
)

// Config represents the structure of the configuration file.
type Config struct {
	RewriteEntries map[string]string `json:"rewrite_entries"`
	ListenAddress  string            `json:"listen_address"`
	UpstreamDNS    string            `json:"upstream_dns"`
}

var config Config

func main() {
	// Load configuration
	err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Set up DNS handler
	dns.HandleFunc(".", handleDNSRequest)

	// Start DNS server
	server := &dns.Server{Addr: config.ListenAddress, Net: "udp"}
	log.Printf("Starting DNS server on %s", config.ListenAddress)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// loadConfig reads and parses the configuration file.
func loadConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read config file: %w", err)
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return fmt.Errorf("could not parse config file: %w", err)
	}

	// Normalize domain names to lowercase
	normalized := make(map[string]string)
	for domain, ip := range config.RewriteEntries {
		normalized[strings.ToLower(domain)] = ip
	}
	config.RewriteEntries = normalized

	return nil
}

// handleDNSRequest processes incoming DNS queries.
func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	msg.Authoritative = true

	for _, question := range r.Question {
		domain := strings.ToLower(question.Name)
		qType := question.Qtype

		if newIP, exists := config.RewriteEntries[domain]; exists {
			log.Printf("Rewriting DNS for %s to %s", domain, newIP)
			rr, err := dns.NewRR(fmt.Sprintf("%s A %s", domain, newIP))
			if err != nil {
				log.Printf("Failed to create RR: %v", err)
				continue
			}
			msg.Answer = append(msg.Answer, rr)
		} else {
			// Forward the query to upstream DNS
			c := new(dns.Client)
			in, _, err := c.Exchange(r, config.UpstreamDNS)
			if err != nil {
				log.Printf("Failed to forward query: %v", err)
				msg.SetRcode(r, dns.RcodeServerFailure)
				continue
			}
			msg.Answer = append(msg.Answer, in.Answer...)
			msg.Authoritative = in.Authoritative
			msg.RecursionAvailable = in.RecursionAvailable
			msg.Rcode = in.Rcode
		}

		// Handle other record types (e.g., AAAA) if needed
		if qType == dns.TypeAAAA {
			msg.SetRcode(r, dns.RcodeNameError)
		}
	}

	w.WriteMsg(&msg)
}
