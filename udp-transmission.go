package main

import (
	"flag"
	"log"
	"net"
)

var clients map[string]*net.UDPAddr
var server *net.UDPConn

func main() {
	address := flag.String("a", ":5555", "usage:「-a=:5555」to specify ip and port")
	flag.Parse()

	// parse server address
	serverAddr, err := net.ResolveUDPAddr("udp4", *address)
	if nil != err {
		log.Printf("Failed to parse address:%v\n", err)
		return
	}

	// listen on udp
	server, err = net.ListenUDP("udp4", serverAddr)
	if nil != err {
		log.Printf("Failed to listen on %s:%v\n", serverAddr, err)
		return
	} else {
		log.Printf("Listen on %s\n", server.LocalAddr().String())
	}
	defer server.Close()

	clients = make(map[string]*net.UDPAddr)
	// udp transparent transmission
	for {
		data := make([]byte, 65535)
		num, readAddr, err := server.ReadFromUDP(data)
		if nil != err {
			log.Printf("Failed to read data:%v\n", err)
			continue
		}
		handle(readAddr, data[:num])
	}
}

func handle(sourceAddr *net.UDPAddr, data []byte) {
	sourceAddress := sourceAddr.String()
	if _, exist := clients[sourceAddress]; !exist {
		clients[sourceAddress] = sourceAddr
		log.Printf("Added new client %s\n", sourceAddress)
	}

	for key, value := range clients {
		if key != sourceAddress {
			go transport(value, data, sourceAddress)
		}
	}
}

func transport(targetAddr *net.UDPAddr, data []byte, sourceAddress string) {
	targetAddress := targetAddr.String()
	num, err := server.WriteToUDP(data, targetAddr)
	if nil != err {
		log.Printf("Failed to send data to %s:%v", targetAddress, err)
		delete(clients, targetAddress)
	} else {
		log.Printf("Transport %d bytes from %s to %s\n", num, sourceAddress, targetAddress)
	}
}
