package main

import (
	"net/netip"
	"testing"
)

func TestAddrInNetwork(t *testing.T) {
	type args struct {
		addr   netip.Addr
		prefix netip.Prefix
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"In the network", args{netip.MustParseAddr("192.168.0.0"), netip.MustParsePrefix("192.168.0.0/24")}, true},
		{"In the network", args{netip.MustParseAddr("192.168.0.254"), netip.MustParsePrefix("192.168.0.0/24")}, true},
		{"In the network", args{netip.MustParseAddr("192.168.1.0"), netip.MustParsePrefix("192.168.0.0/24")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddrInNetwork(tt.args.addr, tt.args.prefix); got != tt.want {
				t.Errorf("AddrInNetwork() = %v, want %v", got, tt.want)
			}
		})
	}
}
