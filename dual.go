package main

import (
	"context"
	"net"

	"github.com/miekg/dns"
)

type dnsResolver struct {
	server string
	stack  string

	r *dns.Client

	co *dns.Conn
}

func (u *dnsResolver) New(stack string, server string) resolver {
	r := new(dns.Client)
	r.SingleInflight = false
	r.Net = stack
	if stack != "udp" {
		var err error
		u.co, err = r.Dial(server)
		if err != nil {
			u.co = nil
		}
	}

	return &dnsResolver{
		server: server,
		stack:  stack,
		r:      r,
	}
}

func (u *dnsResolver) GetIP(ctx context.Context, domain string) (ip net.IP, err error) {
	m := &dns.Msg{
		MsgHdr: dns.MsgHdr{
			Authoritative:     false,
			AuthenticatedData: false,
			CheckingDisabled:  false,
			RecursionDesired:  true,
			Opcode:            dns.OpcodeQuery,
			Rcode:             dns.RcodeSuccess,
		},
		Question: make([]dns.Question, 1),
	}
	qt := dns.TypeA
	qc := uint16(dns.ClassINET)
	m.Question[0] = dns.Question{Name: dns.Fqdn(domain), Qtype: qt, Qclass: qc}
	m.Id = dns.Id()

	var r *dns.Msg
	if u.stack != "udp" {
		if u.co == nil {
			u.co, err = u.r.Dial(u.server)
			if err != nil {
				u.co = nil
			}
		}
		r, _, err = u.r.ExchangeWithConn(m, u.co)
		if err != nil {
			u.co = nil
		}
	} else {
		r, _, err = u.r.Exchange(m, u.server)
	}

	if r == nil || len(r.Answer) <= 0 {
		err = ErrNoValidIP
		return
	}
	return r.Answer[0].(*dns.A).A, err
}
