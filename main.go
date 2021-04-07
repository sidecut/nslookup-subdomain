package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"sync"
)

type Results struct {
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
		go func(c chan Results) {
			defer wgLookups.Done()

			names, err := net.LookupAddr(ipAddress)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error looking up %v: %v\n", ipAddress, err)
			} else {
				results := Results{ipAddress: ipAddress, names: names}
				c <- results
			}
		}(resultsChannel)
	}

	wgLookups.Wait()
	close(resultsChannel)
}

func consumeResults(resultsChannel chan Results) {
	fmt.Println("Address,Name")
	for result := range resultsChannel {
		for _, name := range result.names {
			fmt.Printf("%s,\"%s\"\n", result.ipAddress, name)
		}
	}
}

func main() {
	prefix3octet := flag.String("3", "", "3-octet Address prefix, e.g. 192.168.1. or 192.168.1")
	flag.Parse()

	var prefix3octetString string
	if matched, _ := regexp.MatchString(octet3Regexp, *prefix3octet); matched {
		prefix3octetString = *prefix3octet + "."
	} else if matched, _ := regexp.MatchString(octet3TrailingRegexp, *prefix3octet); matched {
		prefix3octetString = *prefix3octet
	} else {
		flag.Usage()
		os.Exit(1)
	}

	resultsChannel := make(chan Results)
	go produceResults(prefix3octetString, resultsChannel)
	consumeResults(resultsChannel)
}
