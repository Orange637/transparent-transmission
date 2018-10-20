package main

import (
	"flag"
	"log"
	"net"
)

var clients map[string]*net.UDPAddr

func main() {
	address := *flag.String("a", ":5555", "usage:「-a=:5555」to specify ip and port")
	flag.Parse()

	// parse server address
	serverAddr, err := net.ResolveUDPAddr("udp4", address)
	if nil != err {
		log.Printf("Failed to parse address:%v%\n", err)
		return
	}

	// listen on udp
	listener, err := net.ListenUDP("udp4", serverAddr)
	if nil != err {
		log.Printf("Failed to listen on %s:%v\n", serverAddr, err)
		return
	} else {
		log.Printf("Listen on %s\n", listener.LocalAddr().String())
	}
	defer listener.Close()

	clients = make(map[string]*net.UDPAddr)
	// udp transparent transmission
	for {
		data := make([]byte, 4096)
		num, readAddr, err := listener.ReadFromUDP(data)
		if nil != err {
			log.Printf("Failed to read data:%v\n", err)
			continue
		}

		clientAddress := readAddr.String()
		if _, exist := clients[clientAddress]; !exist {
			clients[clientAddress] = readAddr
			log.Printf("Added new client %s\n", clientAddress)
		}

		for key, value := range clients {
			if key != clientAddress {
				_, err = listener.WriteToUDP(data[:num], value)
				if nil != err {
					log.Printf("Failed to send data to %s:%v", key, err)
					delete(clients, key)
				} else {
					log.Printf("Transport %d bytes from %s to %s\n", num, clientAddress, key)
				}
			}
		}
	}
}
