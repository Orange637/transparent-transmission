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

		sourceAddress := readAddr.String()
		if _, exist := clients[sourceAddress]; !exist {
			clients[sourceAddress] = readAddr
			log.Printf("Added new client %s\n", sourceAddress)
		}

		for key, value := range clients {
			if key != sourceAddress {
				go transport(server, data[:num], sourceAddress, value)
			}
		}
	}
}

func transport(server *net.UDPConn, data []byte, sourceAddress string, target *net.UDPAddr) {
	targetAddress := target.String()
	num, err := server.WriteToUDP(data, target)
	if nil != err {
		log.Printf("Failed to send data to %s:%v", targetAddress, err)
		delete(clients, targetAddress)
	} else {
		log.Printf("Transport %d bytes from %s to %s\n", num, sourceAddress, targetAddress)
	}
}
