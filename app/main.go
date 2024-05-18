package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
)

func encodeDomain(domain string) []byte {
	encodes := []byte{}
	for _, seg := range strings.Split(domain, ".") {
		n := len(seg)
		encodes = append(encodes, byte(n))
		encodes = append(encodes, []byte(seg)...)
	}
	return append(encodes, 0x00)
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		_, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		// Process received question
		receivedQuestion := []byte{}
		i := 12
		for i < len(buf) {
			fmt.Println("nikhilk2")
			fmt.Println(buf[i])
			length := int(buf[i])
			fmt.Println(length)
			receivedQuestion = append(receivedQuestion, buf[i])
			if length == 0 {
				break
			}
			i++ // move to the start of the segment
			receivedQuestion = append(receivedQuestion, buf[i:i+length]...)
			i += length // move to the next length prefix
		}

		receivedQuestion = append(receivedQuestion, 0, 1, 0, 1)

		response := []byte{}
		response = append(response, buf[:2]...)
		response = append(response,
			4, 210,
			128,
			0,
			1, 0,
			0, 0,
			0, 0,
			0, 0)

		fmt.Println(len(receivedQuestion))

		response = append(response, receivedQuestion...)

		fmt.Println(len(response))

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
