package servers

import (
	"context"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"
)

// RunTCPServer runs tcp server
func RunTCPServer(address string, port int, intFace string, bufferSize int) {
	listen(fmt.Sprintf("%s:%d", address, port), intFace, bufferSize)
}

func listen(address, vrfName string, bufferSize int) {
	lc := net.ListenConfig{Control: controlOnConnSetup(vrfName)}

	ln, err := lc.Listen(context.Background(), "tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConnection(conn, bufferSize)
	}
}

func handleConnection(conn net.Conn, bufferSize int) {
	log.Print("Start TCP Server")

	for {
		time.Sleep(1 * time.Second)
		buffer := make([]byte, bufferSize)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Print(fmt.Sprintf("Connection from client: %s closed", conn.RemoteAddr()))
			return
		}
		log.Printf("packet-received: bytes=%d from=%s\n",
			n, conn.RemoteAddr())
		conn.Write([]byte(buffer[:n]))
	}

	conn.Close()
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
