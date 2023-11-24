package protocols

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	// ProtocolICMP the name of the protocol
	ProtocolICMP   = "icmp"
	packagesNumber = 5
)

// ICMPTest define, run and process return code of icmp test command
type ICMPTest struct {
	common        CommonTest
	InterfaceName string
}

// NewICMPTest creates new instance of ConnectivityTestParameters
func NewICMPTest(mtu int, protocolVersion int, serverIP string, intefaceName string, negative bool) *ICMPTest {
	if intefaceName != "" {
		intFace, err := net.InterfaceByName(intefaceName)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		intefaceName = intFace.Name
	}
	return &ICMPTest{
		InterfaceName: intefaceName,
		common: CommonTest{
			MTU:             mtu,
			ServerIP:        serverIP,
			ProtocolVersion: protocolVersion,
			Negative:        negative,
		}}
}

func (test *ICMPTest) defineCommand() string {
	command := []string{"ping", fmt.Sprintf("-%d", test.common.ProtocolVersion),
		test.common.ServerIP, "-c", fmt.Sprintf("%d", packagesNumber), "-w", fmt.Sprintf("%d", packagesNumber),
		"-s", fmt.Sprintf("%d", test.common.MTU), "-M", "do"}
	if test.InterfaceName != "" {
		command = append(command, fmt.Sprintf("-I %s", test.InterfaceName))
	}
	return strings.Join(command, " ")
}

// RunTest runs the test
func (test *ICMPTest) RunTest() {
	_, err := test.common.RunCommand(test.defineCommand())
	if test.common.Negative {
		if err != nil {
			log.Print("ICMP test failed as expected")
			os.Exit(0)
		} else {
			log.Fatalf("negative test failed to return code 1")
			os.Exit(1)
		}
	}
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		os.Exit(1)
	}
}
