package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/kononovn/testcmd/protocols"
	"github.com/kononovn/testcmd/servers"
)

const (
	ipv4BroadcastAddress = "255.255.255.255"
)

var (
	supportedProtocols       = []string{protocols.ProtocolICMP, protocols.ProtocolUDP, protocols.ProtocolTCP, protocols.ProtocolSCTP}
	supportedServerProtocols = []string{protocols.ProtocolUDP, protocols.ProtocolSCTP}
)

func validateIP(host string, multicast bool) error {
	ip := net.ParseIP(host)
	if multicast {
		if ip.IsMulticast() {
			return nil
		}
		return fmt.Errorf("Unsupported parameter server ip=%s is not mulicast address", host)
	}
	if ip != nil {
		return nil
	}
	return fmt.Errorf("Unsupported parameter server ip=%s", host)
}

func ipProtocolVersion(host string) int {
	if strings.Contains(host, ":") {
		return 6
	}
	return 4
}

func validateIntInRange(testInt int, rangeStart int, rangeStop int) error {
	if testInt >= rangeStart && testInt <= rangeStop {
		return nil
	}
	return fmt.Errorf("value=%d not in range %d...%d", testInt, rangeStart, rangeStop)
}

func validateMtu(mtuSize int) error {
	err := validateIntInRange(mtuSize, 50, 9000)
	if err != nil {
		return fmt.Errorf("unsupported parameter mtu=%d %s", mtuSize, err)
	}
	return nil
}

func validatePort(portNumber int) error {
	err := validateIntInRange(portNumber, 1, 65534)
	if err != nil {
		return fmt.Errorf("unsupported parameter port=%d %s", portNumber, err)
	}
	return nil
}

func validateProtocol(protocolName string) error {
	for _, item := range supportedProtocols {
		if protocolName == item {
			return nil
		}
	}
	return fmt.Errorf("Unsupported parameter protocol=%s", protocolName)
}

func main() {
	serverMode := flag.Bool("listen", false, "Insert this flag in order to run server")
	interfaceName := flag.String("interface", "", "Interface name. Examples: ens33/eth0/net1")
	multicast := flag.Bool("multicast", false, "Insert this flag in order to run udp multicast server")
	broadcast := flag.Bool("broadcast", false, "Insert this flag in order to run udp broadcast server")
	protocol := flag.String("protocol", "", "Protocol name. Options: tcp/udp/icmp/sctp")
	mtu := flag.Int("mtu", 1450, "MTU Size. Options: Any int in range 50-9000")
	dstAddress := flag.String("server", "", "Destination ip address IPv4/IPv6")
	serverPort := flag.Int("port", 80, "Port number. Options: Any int in range 1-65534")
	negative := flag.Bool("negative", false, "Insert this flag if no connectivity expected")
	flag.Parse()

	err := validateProtocol(*protocol)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err = validateMtu(*mtu)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if *serverMode {
		err = validatePort(*serverPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if *multicast {
			err = validateIP(*dstAddress, *multicast)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			protocolVersion := ipProtocolVersion(*dstAddress)
			servers.RunMulticastUDPServer(*serverPort, *dstAddress, protocolVersion, *mtu, *interfaceName)
		} else if *broadcast {
			servers.RunBroadcastUDPServer(*serverPort, ipv4BroadcastAddress, *mtu, *interfaceName)
		} else {
			switch *protocol {
			case protocols.ProtocolUDP:
				if *dstAddress != "" {
					log.Printf("Parameter -server=%s ignored in server UDP unicast mode. Use all interfaces 0.0.0.0", *dstAddress)
				}
				servers.RunUDPServer(*serverPort, *mtu)
			case protocols.ProtocolSCTP:
				servers.RunSCTP(*dstAddress, *serverPort, *mtu, *interfaceName, ipProtocolVersion(*dstAddress))
			case protocols.ProtocolTCP:
				servers.RunTCPServer(*dstAddress, *serverPort, *interfaceName, *mtu)
			}
		}
		return
	}

	err = validateIP(*dstAddress, *multicast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	protocolVersion := ipProtocolVersion(*dstAddress)

	switch *protocol {
	case protocols.ProtocolICMP:
		test := protocols.NewICMPTest(*mtu, protocolVersion, *dstAddress, *interfaceName, *negative)
		test.RunTest()

	case protocols.ProtocolTCP:
		err = validatePort(*serverPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		test := protocols.NewTCPTest(*mtu, protocolVersion, *dstAddress, *serverPort, *negative, *interfaceName)
		test.RunTest()

	case protocols.ProtocolUDP:
		err = validatePort(*serverPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		test := protocols.NewUDPTest(*mtu, protocolVersion, *dstAddress, *serverPort, *negative, *multicast, *broadcast, *interfaceName)
		test.RunTest()

	case protocols.ProtocolSCTP:
		err = validatePort(*serverPort)
		if err != nil {
			log.Fatalf("port validation error: %v\n", err)
		}
		test := protocols.NewSCTPTest(*mtu, *dstAddress, protocolVersion, *serverPort, *negative)
		test.RunTest()
	}
}
