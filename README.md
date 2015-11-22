# DNS Ping

Simple Go program to ping DNS servers.

It sends periodically DNS A requests to a resolver, for a given domain, and
displays the time needed to get the reply.

## Usage

First:

    go build

Then:

    ./dnsping
    dnsping - monitor response time for DNS servers
    Send DNS requests periodically to monitor a DNS server response time.

    Usage: dnsping [option ...] targetdomain
      -count=0: Number of requests to send
      -d=1000: Interval to wait between two pings
      -i=false: Do an iterative query instead of recursive (to stress authoritative nameservers)
      -r="127.0.0.1:53": Resolver to test against
      -random=false: Use random Request Identifiers for each query
      -v=false: Verbose logging
