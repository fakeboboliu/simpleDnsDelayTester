package main

import (
	"context"
	"net"

	"github.com/babolivier/go-doh-client"
)

type dohResolver struct {
	server string
	r      *doh.Resolver
}

func (r *dohResolver) New(server string) resolver {
	nr := &doh.Resolver{
		Host:  server,
		Class: doh.IN,
	}

	return &dohResolver{server: server, r: nr}
}

func (r *dohResolver) GetIP(ctx context.Context, domain string) (ip net.IP, err error) {
	a, _, err := r.r.LookupA(domain)
	if err != nil {
		return
	}
	if len(a) <= 0 {
		err = ErrNoValidIP
		return
	}

	return net.ParseIP(a[0].IP4), nil
}
