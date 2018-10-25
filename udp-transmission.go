package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

var inConn *net.UDPConn
var outConn *net.UDPConn
var device *net.UDPAddr
var client *net.UDPAddr

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

	go connectToClient()

	// udp transparent transmission
	for {
		data := make([]byte, 65535)
		num, readAddr, err := inConn.ReadFromUDP(data)
		if nil != err {
			log.Printf("Failed to read data:%v\n", err)
			continue
		}
		if nil == client {
			continue
		}
		transportToClient(data[:num], readAddr)
	}
}

func connectToClient() {
	for {
		data := make([]byte, 65535)
		_, readAddr, err := outConn.ReadFromUDP(data)
		if nil != err {
			log.Printf("Failed to connect client:%v\n", err)
		} else if nil == client || client.String() != readAddr.String() {
			fmt.Printf("New client connected:%v\n", readAddr)
			client = readAddr
		}
	}
}

func transportToClient(data []byte, deviceAddr *net.UDPAddr) {
	if nil == client {
		return
	}
	num, err := outConn.WriteToUDP(data, client)
	if nil != err {
		log.Printf("Failed to send data to %s:%v", client.String(), err)
	} else if false {
		log.Printf("Transport %d bytes from %s to %s\n", num, deviceAddr.String(), client.String())
	}
}
