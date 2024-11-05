package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	netcat "netcat/ressources"
)

const defaultPort = "8989"

// Program entry point
func main() {
	port := defaultPort
	if len(os.Args) == 2 {
		port = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.Create("server.log")
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	log.Println("Chat server starting on port :" + port)
	
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go netcat.HandleClient(conn)
	}
}
