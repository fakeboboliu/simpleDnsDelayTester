package main

import (
	"context"
	"errors"
	"net"
)

type resolver interface {
	GetIP(ctx context.Context, domain string) (ip net.IP, err error)
}

var ErrNoValidIP = errors.New("no valid ip")
