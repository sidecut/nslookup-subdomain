package main

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"sync"

	"github.com/spf13/pflag"
)

type Results struct {
	index     int
	ipAddress string
	names     []string
}

const octetRegexp = `\d{1,3}`
const periodRegexp = "[.]"
const octet3Regexp = "^" + octetRegexp + periodRegexp + octetRegexp + periodRegexp + octetRegexp + "$"
const octet3TrailingRegexp = "^" + octetRegexp + periodRegexp + octetRegexp + periodRegexp + octetRegexp + periodRegexp + "$"

func produceResults(addrPrefix string, resultsChannel chan Results) {
	var wgLookups sync.WaitGroup

	for i := 0; i < 256; i++ {
		ipAddress := fmt.Sprintf("%s%03d", addrPrefix, i)

		wgLookups.Add(1)
		go func(i int, c chan Results) {
			defer wgLookups.Done()

			names, err := net.LookupAddr(ipAddress)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error looking up %v: %v\n", ipAddress, err)
			} else {
				results := Results{index: i, ipAddress: ipAddress, names: names}
				c <- results
			}
		}(i, resultsChannel)
	}

	wgLookups.Wait()
	close(resultsChannel)
}

func consumeAndOutputResults(resultsChannel chan Results) {
	var results [256]Results
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
	prefix3octet := pflag.StringP("prefix3", "3", "", "3-octet Address prefix, e.g. 192.168.1. or 192.168.1")
	pflag.Parse()

	var prefix3octetString string
	if matched, _ := regexp.MatchString(octet3Regexp, *prefix3octet); matched {
		prefix3octetString = *prefix3octet + "."
	} else if matched, _ := regexp.MatchString(octet3TrailingRegexp, *prefix3octet); matched {
		prefix3octetString = *prefix3octet
	} else {
		pflag.Usage()
		os.Exit(1)
	}

	resultsChannel := make(chan Results)
	go produceResults(prefix3octetString, resultsChannel)
	consumeAndOutputResults(resultsChannel)
}
