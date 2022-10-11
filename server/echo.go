package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:56789")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("starting connection with %s", conn.RemoteAddr().String())
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			log.Print("ending connection")
			break
		}
		if err != nil {
			log.Print(err)
		}
		m, err := conn.Write(buffer[:n])
		if err != nil {
			log.Print(err)
		}
		if n != m {
			log.Print(fmt.Errorf("read-write unequal length"))
		}
	}
}
