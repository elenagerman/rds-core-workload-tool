package protocols

import (
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"
)

const (
	// ProtocolTCP the name of the protocol
	ProtocolTCP       = "tcp"
	packagesNumberTCP = 5
	timeoutTCP        = 2
	timeoutDialTCP    = 10
)

// TCPTest define, run and process return code of tcp test command
type TCPTest struct {
	CommonTest
	ServerPort    int
	InterfaceName *net.Interface
}

// RunTest runs the test
func (test *TCPTest) RunTest() {
	err := test.testTCP()
	if test.Negative {
		if err != nil {
			log.Print("Negative TCP test passed")
			os.Exit(0)
		}
		log.Print("Negative TCP test failed")
		os.Exit(1)
	}

	if err == nil {
		log.Print("TCP test passed as expected")
		os.Exit(0)
	}
	log.Print("TCP test failed")
	os.Exit(1)
}

// NewTCPTest creates new instance of ConnectivityTestParameters
func NewTCPTest(mtu int, protocolVersion int, serverIP string, serverPort int, negative bool, interfaceName string) *TCPTest {
	intFace, err := net.InterfaceByName(interfaceName)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	return &TCPTest{CommonTest{mtu, serverIP, protocolVersion, negative}, serverPort, intFace}
}

func (test *TCPTest) testTCP() error {
	raddr := test.resolveAddress()
	dialer := net.Dialer{Timeout: timeoutDialTCP * time.Second, Control: controlOnConnSetup(test.InterfaceName.Name)}
	connection, err := dialer.Dial(
		fmt.Sprintf("%s%d", ProtocolTCP, test.ProtocolVersion),
		raddr.String())

	if err != nil {
		return err
	}
	var testString string
	for i := 0; i < test.MTU; i++ {
		testString += "a"
	}

	fmt.Printf(fmt.Sprintf("TCP PING %s %d(%d) bytes of data.\n", test.ServerIP, test.MTU, test.MTU+28))
	var (
		statTotalTime      int64
		exitCode           int
		statPacketLost     int
		statPacketReceived int
	)
	for i := 1; i <= packagesNumberTCP; i++ {
		byteTestString := []byte(testString)
		test.runTCPPing(connection, i, byteTestString, &statPacketLost, &statPacketReceived, &statTotalTime, &exitCode)
	}

	fmt.Printf("--- %s TCP statistics ---\n", test.ServerIP)
	fmt.Printf(
		"%d packets transmitted, %d received, %d packet loss, time %dms\n",
		packagesNumberTCP,
		statPacketReceived,
		totalPackageLoss(packagesNumberTCP, statPacketLost),
		statTotalTime)

	if exitCode == 1 {
		return fmt.Errorf("TCP connectivity test returns error code")
	}
	return nil
}

func (test *TCPTest) runTCPPing(
	conn net.Conn,
	packetNumber int,
	byteTestString []byte,
	statPacketLost *int,
	statPacketReceived *int,
	statTotalTime *int64,
	exitCode *int) {

	time.Sleep(1 * time.Second)

	deadline := time.Now().Add(timeoutTCP * time.Second)
	conn.SetDeadline(deadline)
	startTime := time.Now()
	_, err := conn.Write(byteTestString)
	if err != nil {
		fmt.Print(err)
	}

	buffer := make([]byte, test.MTU)
	readBufferSized, err := conn.Read(buffer)
	elapsed := time.Since(startTime)
	if err != nil {
		fmt.Printf("Package lost\n")
		*statPacketLost++
		*exitCode = 1
		return
	}

	if string(buffer) == string(byteTestString) {
		*statTotalTime += elapsed.Microseconds()
		fmt.Printf("%d bytes from %s: tcp_seq=%d time=%dms\n",
			readBufferSized, conn.RemoteAddr(), packetNumber, elapsed.Microseconds())
		*statPacketReceived++
	} else {
		*exitCode = 1
	}
}

func (test *TCPTest) resolveAddress() *net.TCPAddr {
	addr, err := net.ResolveTCPAddr(fmt.Sprintf("%s%d", ProtocolTCP, test.ProtocolVersion),
		fmt.Sprintf("[%s]:%d", test.ServerIP, test.ServerPort))
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	return addr
}

func controlOnConnSetup(vrfName string) func(network string, address string, c syscall.RawConn) error {
	return func(network string, address string, c syscall.RawConn) error {
		if vrfName == "" {
			return nil
		}
		var operr error
		fn := func(fd uintptr) {
			operr = syscall.SetsockoptString(int(fd), syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, vrfName)
		}
		if err := c.Control(fn); err != nil {
			return err
		}
		if operr != nil {
			return operr
		}
		return nil
	}
}
