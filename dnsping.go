package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/miekg/dns"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

// Runtime options
var (
	count        int
	pingInterval int
	verbose      bool
	iterative    bool
	resolver     string
	randomIds    bool
)

func init() {
	flag.IntVar(&pingInterval, "d", 1000,
		"Interval to wait between two pings")
	flag.BoolVar(&verbose, "v", false,
		"Verbose logging")
	flag.IntVar(&count, "count", 0,
		"Number of requests to send")
	flag.BoolVar(&randomIds, "random", false,
		"Use random Request Identifiers for each query")
	flag.BoolVar(&iterative, "i", false,
		"Do an iterative query instead of recursive (to stress authoritative nameservers)")
	flag.StringVar(&resolver, "r", "127.0.0.1:53",
		"Resolver to test against")
}

func main() {
	fmt.Printf("dnsping - monitor response time for DNS servers\n")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, strings.Join([]string{
			"Send DNS requests periodically to monitor a DNS server response time.",
			"",
			"Usage: dnsping [option ...] targetdomain",
			"",
		}, "\n"))
		flag.PrintDefaults()
	}

	flag.Parse()

	// We need exactly one target domain
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	targetDomain := flag.Args()[0]
	fmt.Printf("Pinging resolver %s with domain %s\n", resolver, targetDomain)

	sent, errors := pinguer(targetDomain)
	fmt.Printf(
		"Statistics: %d requests sent, %d received (%.0f%% error)\n",
		sent,
		sent-errors,
		float64(errors)/float64(sent),
	)
}

func pinguer(domain string) (int, int) {

	// Every N steps, we will tell the stats module how many requests we sent
	maxRequestID := big.NewInt(65536)
	errors := 0
	totalSent := 0

	questionRecord := dns.TypeA

	message := new(dns.Msg).SetQuestion(domain, questionRecord)
	if iterative {
		message.RecursionDesired = false
	}

	for reqnumber := 0; count == 0 || reqnumber < count; reqnumber++ {

		// Try to resolve the domain
		if randomIds {
			// Regenerate message Id to avoid servers dropping (seemingly) duplicate messages
			newid, _ := rand.Int(rand.Reader, maxRequestID)
			message.Id = uint16(newid.Int64())
		}

		start := time.Now()
		err := dnsExchange(resolver, message)
		elapsedMilliSeconds := float64(time.Since(start)) / float64(time.Millisecond)

		if err != nil {
			if verbose {
				fmt.Printf("%s error: % (%s)\n", domain, err, resolver)
			}
			errors++
		}

		totalSent++

		// Display results of the ping
		fmt.Printf(
			"ping %s with %s %s: %.3fms\n",
			resolver,
			dns.TypeToString[questionRecord],
			domain,
			elapsedMilliSeconds,
		)

		time.Sleep(time.Duration(pingInterval) * time.Millisecond)
	}

	return totalSent, errors
}

func dnsExchange(resolver string, message *dns.Msg) error {
	dnsconn, err := net.Dial("udp", resolver)
	if err != nil {
		return err
	}
	co := &dns.Conn{Conn: dnsconn}
	defer co.Close()

	// Actually send the message and wait for answer
	co.WriteMsg(message)

	_, err = co.ReadMsg()
	return err
}
