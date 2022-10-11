package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var (
	message = flag.String("message", "Dies ist ein Test", "echo-Nachricht")
)

func main() {
	flag.Parse()

	conn, err := net.DialTimeout("tcp", "localhost:56789", 3*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	const maxBufferSize = 1024
	cnt := 0
	ch := make(chan string)
	done := make(chan struct{})
	go func() {
		recv := make([]byte, maxBufferSize)
		for {
			select {
			case <-done:
				close(ch)
				return

			default:
				n, err := conn.Read(recv)
				if err != nil {
					log.Printf("error receiving form server %w", err)
				}
				ch <- string(recv[:n])
			}
		}
	}()
	repeat := time.Tick(time.Second)
	runtime := time.After(10 * time.Second)
	for {
		timeout := time.After(2 * time.Second)
		select {
		case s := <-ch:
			fmt.Printf("Received: %s\n", s)

		case <-repeat:
			cnt++
			msg := strings.TrimSpace(*message) + fmt.Sprintf("%5d", cnt)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				log.Fatalf("failed to write to server %w", err)
			}

		case <-timeout:
			log.Print("timeout for write-read")

		case <-runtime:
			log.Print("ending connection")
			close(done)
			return
		}
	}
}
