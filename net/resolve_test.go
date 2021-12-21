package manet

import (
	"testing"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

func TestResolvingAddrs(t *testing.T) {
	unspec := []ma.Multiaddr{
		newMultiaddr(t, "/ip4/0.0.0.0/tcp/1234"),
		newMultiaddr(t, "/ip4/1.2.3.4/tcp/1234"),
		newMultiaddr(t, "/ip6/::/tcp/1234"),
		newMultiaddr(t, "/ip6/::100/tcp/1234"),
	}

	iface := []ma.Multiaddr{
		newMultiaddr(t, "/ip4/127.0.0.1"),
		newMultiaddr(t, "/ip4/10.20.30.40"),
		newMultiaddr(t, "/ip6/::1"),
		newMultiaddr(t, "/ip6/::f"),
	}

	spec := []ma.Multiaddr{
		newMultiaddr(t, "/ip4/127.0.0.1/tcp/1234"),
		newMultiaddr(t, "/ip4/10.20.30.40/tcp/1234"),
		newMultiaddr(t, "/ip4/1.2.3.4/tcp/1234"),
		newMultiaddr(t, "/ip6/::1/tcp/1234"),
		newMultiaddr(t, "/ip6/::f/tcp/1234"),
		newMultiaddr(t, "/ip6/::100/tcp/1234"),
	}

	actual, err := ResolveUnspecifiedAddresses(unspec, iface)
	require.NoError(t, err)
	require.Equal(t, actual, spec)

	ip4u := []ma.Multiaddr{newMultiaddr(t, "/ip4/0.0.0.0")}
	ip4i := []ma.Multiaddr{newMultiaddr(t, "/ip4/1.2.3.4")}

	ip6u := []ma.Multiaddr{newMultiaddr(t, "/ip6/::")}
	ip6i := []ma.Multiaddr{newMultiaddr(t, "/ip6/::1")}

	if _, err := ResolveUnspecifiedAddress(ip4u[0], ip6i); err == nil {
		t.Fatal("should have failed")
	}
	if _, err := ResolveUnspecifiedAddress(ip6u[0], ip4i); err == nil {
		t.Fatal("should have failed")
	}

	if _, err := ResolveUnspecifiedAddresses(ip6u, ip4i); err == nil {
		t.Fatal("should have failed")
	}
	if _, err := ResolveUnspecifiedAddresses(ip4u, ip6i); err == nil {
		t.Fatal("should have failed")
	}
}

func TestAddrOverNonLocalIP(t *testing.T) {
	bad := []ma.Multiaddr{
		newMultiaddr(t, "/ip6/fe80::1/tcp/1234"),   // link local
		newMultiaddr(t, "/ip6/fe80::100/tcp/1234"), // link local
	}
	good := []ma.Multiaddr{
		newMultiaddr(t, "/ip4/127.0.0.1/tcp/1234"),
		newMultiaddr(t, "/ip6/::1/tcp/1234"),
		newMultiaddr(t, "/ip4/1.2.3.4/udp/1234/utp"),
	}
	for _, addr := range bad {
		require.Falsef(t, IsIpv6LinkLocal(addr), "%s is a link local addr", addr)
	}
	for _, addr := range good {
		require.Truef(t, IsIpv6LinkLocal(addr), "%s is not a link local addr", addr)
	}
}
