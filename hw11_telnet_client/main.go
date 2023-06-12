package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	timeoutFlag := flag.Duration("timeout", 10*time.Second, "timeout for connection to the server")
	flag.Parse()

	args := flag.Args()

	if len(args) < 2 {
		log.Fatal("Usage: go-telnet --timeout=5s host port")
	}

	host := args[0]
	port := args[1]
	address := net.JoinHostPort(host, port)

	client := NewTelnetClient(address, *timeoutFlag, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		cancel()
	}()

	go func() {
		err = client.Receive()
		if err != nil {
			log.Printf("Failed to receive: %v\n", err)
		}
		cancel()
	}()

	go func() {
		err = client.Send()
		if err != nil {
			log.Printf("Failed to send: %v\n", err)
		}
		cancel()
	}()

	<-ctx.Done()
}
