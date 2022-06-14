package main

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"sync"

	"github.com/spf13/pflag"
)

type Results struct {
	index     int
	ipAddress string
	names     []string
}

const (
	octetRegexp          = `\d{1,3}`
	periodRegexp         = "[.]"
	octet3Regexp         = "^" + octetRegexp + periodRegexp + octetRegexp + periodRegexp + octetRegexp + "$"
	octet3TrailingRegexp = "^" + octetRegexp + periodRegexp + octetRegexp + periodRegexp + octetRegexp + periodRegexp + "$"
)

var (
	cidr *string
)

func produceResults(prefix netip.Prefix, resultsChannel chan Results) {
	var wgLookups sync.WaitGroup

	// fmt.Printf("%v\n", prefix.Masked().Addr())

	addr := prefix.Masked().Addr()
	i := 0
	for AddrInNetwork(addr, prefix) {
		wgLookups.Add(1)
		go func(i int, addr netip.Addr, c chan Results) {
			defer wgLookups.Done()

			names, err := net.LookupAddr(addr.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error looking up %v: %v\n", addr, err)
			} else {
				results := Results{index: i, ipAddress: addr.String(), names: names}
				c <- results
			}
		}(i, addr, resultsChannel)
		i++
		addr = addr.Next()
	}

	wgLookups.Wait()
	close(resultsChannel)
}

func AddrInNetwork(addr netip.Addr, prefix netip.Prefix) bool {
	return prefix.Contains(addr)
}

func consumeAndOutputResults(resultsChannel chan Results) {
	results := make(map[int]Results)
	for result := range resultsChannel {
		results[result.index] = result
	}

	fmt.Println("Address,Name")
	for _, result := range results {
		for _, name := range result.names {
			fmt.Printf("%s,\"%s\"\n", result.ipAddress, name)
		}
	}
}

func main() {
	cidr = pflag.String("cidr", "", "CIDR range")
	pflag.Parse()

	if *cidr == "" {
		pflag.Usage()
		fmt.Fprintf(os.Stderr, `
	See https://golang.org/pkg/net/#hdr-Name_Resolution for details on using
	environment variables to force use of the golang resolver, which will return
	more than one domain name.

	Example:
	$ export GODEBUG=netdns=go    # force pure Go resolver
`)
		os.Exit(1)
	}

	prefix, err := netip.ParsePrefix(*cidr)
	if err != nil {
		log.Fatalln(err)
	}

	// fmt.Printf("%v %v\n", ipv4Addr, ipv4Net)

	resultsChannel := make(chan Results)
	go produceResults(prefix, resultsChannel)
	consumeAndOutputResults(resultsChannel)
}
