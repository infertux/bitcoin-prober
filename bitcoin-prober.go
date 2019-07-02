package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/gcash/bchd/chaincfg"
	"github.com/gcash/bchd/peer"
	"github.com/gcash/bchd/wire"

	"github.com/oschwald/geoip2-golang"
)

func probePeer(host string, bitcoinNet wire.BitcoinNet) (*wire.MsgVersion, error) {
	vermsg := make(chan *wire.MsgVersion)

	chainParams := chaincfg.MainNetParams
	chainParams.Net = bitcoinNet

	peerConfig := &peer.Config{
		UserAgentName:    "bitcoin-prober",
		UserAgentVersion: "1.0.0",
		ChainParams:      &chainParams,
		Services:         0,
		DisableRelayTx:   true,
		Listeners: peer.MessageListeners{
			OnVersion: func(peer *peer.Peer, msg *wire.MsgVersion) *wire.MsgReject {
				vermsg <- msg
				return nil
			},
		},
	}
	peer, err := peer.NewOutboundPeer(peerConfig, host)
	if err != nil {
		return nil, err
	}

	connection, err := net.DialTimeout("tcp", peer.Addr(), time.Second*5)
	if err != nil {
		return nil, err
	}
	peer.AssociateConnection(connection)

	select {
	case version := <-vermsg:
		peer.Disconnect()
		peer.WaitForDisconnect()
		return version, nil
	case <-time.After(time.Second * 5):
		peer.Disconnect()
		peer.WaitForDisconnect()
		errorString := fmt.Sprintf("timeout for peer %v", peer)
		return nil, errors.New(errorString)
	}
}

func outputPeer(address string, msg *wire.MsgVersion, verbose bool) error {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return err
	}

	geoipDB, err := geoip2.Open("/usr/share/GeoIP/GeoLite2-City.mmdb")
	if err == nil {
		defer geoipDB.Close()

		ip := net.ParseIP(host)
		geoipResult, err2 := geoipDB.City(ip)
		if err2 != nil {
			return err2
		}

		country := geoipResult.Country.IsoCode
		city := geoipResult.City.Names["en"]

		fmt.Printf("%s is located in %s", host, country)
		if city != "" {
			fmt.Printf(", %s", city)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("UserAgent: %v\n", msg.UserAgent)
	fmt.Printf("Services: %v\n", msg.Services)
	fmt.Printf("ProtocolVersion: %v\n", msg.ProtocolVersion)
	fmt.Printf("LastBlock: %v\n", msg.LastBlock)
	fmt.Printf("RelayTx: %v\n", !msg.DisableRelayTx)

	if verbose {
		fmt.Printf("%+v\n", msg)

		// TODO: output actual hex payload
		fmt.Printf("Replay with netcat: echo \"%X\" | nc -v %s %s\n",
			"payload", host, port)
	}

	return nil
}

func main() {
	networks := map[string]wire.BitcoinNet{
		"BCH": 0xe8f3e1e3,
		"BTC": 0xd9b4bef9,
	}

	address := flag.String("address", "", "Address to probe")
	network := flag.String("network", "BCH", "Network (BCH or BTC)")
	verbose := flag.Bool("verbose", false, "Be verbose")
	flag.Parse()

	bitcoinNet, ok := networks[*network]
	if !ok {
		fmt.Printf("Unknown network %v", *network)
		return
	}

	normalizedAddress := normalizeAddress(*address, "8333")

	connection, err := net.DialTimeout("tcp", normalizedAddress, time.Second*5)
	if err != nil {
		fmt.Printf("Error looking up %v: %v\n", normalizedAddress, err)
		return
	}

	ip := connection.RemoteAddr().String()
	if normalizedAddress != ip {
		fmt.Printf("Resolved %s to %s\n", *address, ip)
	}

	fmt.Printf("Probing %v on the %v network...\n", ip, *network)
	msg, err := probePeer(ip, bitcoinNet)
	if err != nil {
		fmt.Printf("Error probing %v: %v\n", ip, err)
		return
	}

	if err := outputPeer(ip, msg, *verbose); err != nil {
		fmt.Printf("Could not output peer %v: %v", ip, err)
	}
}

// normalizeAddress returns address with the passed default port appended if
// there is not already a port specified.
func normalizeAddress(address, defaultPort string) string {
	_, _, err := net.SplitHostPort(address)
	if err != nil {
		return net.JoinHostPort(address, defaultPort)
	}
	return address
}
