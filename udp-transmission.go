package main

import (
	"flag"
	"log"
	"net"
)

var inConn *net.UDPConn
var outConn *net.UDPConn
var device *net.UDPAddr
var clients map[string]*net.UDPAddr

func main() {
	inAddress := *flag.String("i", ":5555", "usage:「-i=:5555」to specify in ip and port")
	outAddress := *flag.String("o", ":5557", "usage:「-o=:5557」to specify out ip and port")
	flag.Parse()

	// parse server address
	inAddr, err := net.ResolveUDPAddr("udp4", inAddress)
	// client, err = net.ResolveUDPAddr("udp4", inAddress)
	if nil != err {
		log.Printf("Failed to parse in address:%v \n", err)
		return
	}
	outAddr, err := net.ResolveUDPAddr("udp4", outAddress)
	if nil != err {
		log.Printf("Failed to parse out address:%v \n", err)
		return
	}

	// listen on udp
	inConn, err = net.ListenUDP("udp4", inAddr)
	if nil != err {
		log.Printf("Failed to listen on %s:%v\n", inAddr, err)
		return
	} else {
		log.Printf("Listen on %s\n", inConn.LocalAddr().String())
	}
	outConn, err = net.ListenUDP("udp4", outAddr)
	if nil != err {
		log.Printf("Failed to listen on %s:%v\n", outAddr, err)
		return
	} else {
		log.Printf("Listen on %s\n", outConn.LocalAddr().String())
	}
	defer inConn.Close()
	defer outConn.Close()

	clients = make(map[string]*net.UDPAddr)
	go connectClient()

	// udp transparent transmission
	for {
		data := make([]byte, 65535)
		num, readAddr, err := inConn.ReadFromUDP(data)
		if nil != err {
			log.Printf("Failed to read data:%v\n", err)
			continue
		}
		if 0 >= len(clients) {
			continue
		}

		for _, client := range clients {
			num, err = outConn.WriteToUDP(data[:num], client)
			if nil != err {
				log.Printf("Failed to send data to %s:%v", client.String(), err)
			} else if false {
				log.Printf("Transport %d bytes from %s to %s\n", num, readAddr.String(), client.String())
			}
		}
	}
}

func connectClient() {
	for {
		data := make([]byte, 65535)
		_, clientAddr, err := outConn.ReadFromUDP(data)
		if nil != err {
			log.Printf("Failed to connect client:%v\n", err)
			continue
		}

		if _, added := clients[clientAddr.String()]; !added {
			clients[clientAddr.String()] = clientAddr
		}
	}
}
