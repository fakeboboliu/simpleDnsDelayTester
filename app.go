package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	dnsAddr string
	dnsPort uint
	query   string

	resolv resolver
)

var tpl = "%s %s %d\n" // host ip time

func main() {
	flag.StringVar(&dnsAddr, "addr", "1.1.1.1", "dns server to query")
	flag.UintVar(&dnsPort, "port", 53, "port of dns server")
	flag.StringVar(&query, "domain", "github.com", "domain to query")
	interval := flag.Duration("t", time.Second, "Interval (eg. 1s)")

	udp := flag.Bool("udp", true, "use udp mode (default)")
	tcp := flag.Bool("tcp", false, "use tcp mode")
	doh := flag.Bool("doh", false, "use dns-over-https mode")
	dot := flag.Bool("dot", false, "use dns-over-tls mode")

	help := flag.Bool("h", false, "display help message")

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	addr := fmt.Sprintf("%s:%d", dnsAddr, dnsPort)
	switch {
	case *tcp:
		resolv = (&dnsResolver{}).New("tcp", addr)
	case *doh:
		resolv = (&dohResolver{}).New(addr)
	case *dot:
		resolv = (&dnsResolver{}).New("tcp-tls", addr)
	case *udp:
		resolv = (&dnsResolver{}).New("udp", addr)
	}

	ctx := context.Background()
	ticker := time.NewTicker(*interval)
	for {
		<-ticker.C
		prev := time.Now()
		ip, err := resolv.GetIP(ctx, query)
		after := time.Now()
		cost := after.Sub(prev).Milliseconds()
		if err == ErrNoValidIP {
			fmt.Printf(tpl, query, "no_ipv4", cost)
			continue
		}
		if err != nil {
			fmt.Printf(tpl, query, "unknown", cost)
			continue
		}
		fmt.Printf(tpl, query, ip.String(), cost)
	}
}
