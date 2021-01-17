# simpleDnsDelayTester
A dead simple tool for you to test delay to Internet by making dns query.

## Usage

```
~ sddt -h
  -addr string
        dns server to query (default "1.1.1.1")
  -doh
        use dns-over-https mode
  -domain string
        domain to query (default "github.com")
  -dot
        use dns-over-tls mode
  -h    display help message
  -port uint
        port of dns server (default 53)
  -t duration
        Interval (eg. 1s) (default 1s)
  -tcp
        use tcp mode
  -udp
        use udp mode (default) (default true)
```

### Test w/ udp

```
~ sddt 
```

### Test w/ tcp

confirm your dns server supports tcp before using tcp mode

```
~ sddt -tcp
```

### Test w/ DNS-over-TLS

confirm your dns server supports DoT before using this mode

DoT server typically listens on port 853

```
~ sddt -dot -addr dot.pub -port 853
```

### Test w/ DNS-over-HTTPS

confirm your dns server supports DoH before using this mode

DoH server typically behaves like a web server, but serveing dns content in `/dns-query` path

```
~ sddt -dot -addr doh.pub -port 443
```

## Build

```
go get -u github.com/fakeboboliu/simpleDnsDelayTester
go build -trimpath -ldflags "-s -w" -o sddt
#(optional) upx sddt
```
