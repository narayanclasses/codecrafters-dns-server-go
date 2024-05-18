package main

import (
	"encoding/binary"
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

		// receivedBytes := []byte(string(buf[:size]))

		// i := 0
		// for i < len(receivedBytes) {
		// 	fmt.Println("nikhilk")
		// 	fmt.Println(receivedBytes[i])
		// 	i++
		// }

		// Process received question
		receivedQuestion := []byte{}
		// i := 12
		// for i < len(receivedBytes) {
		// 	// fmt.Println("nikhilk2")
		// 	// fmt.Println(receivedBytes[i])
		// 	length := int(receivedBytes[i])
		// 	// fmt.Println(length)
		// 	receivedQuestion = append(receivedQuestion, receivedBytes[i])
		// 	if length == 0 {
		// 		break
		// 	}
		// 	i++ // move to the start of the segment
		// 	receivedQuestion = append(receivedQuestion, receivedBytes[i:i+length]...)
		// 	i += length // move to the next length prefix
		// }

		receivedQuestion = append(receivedQuestion, encodeDomain("codecrafters.io")...)

		receivedQuestion = append(receivedQuestion, 0, 1, 0, 1)

		// response := []byte{
		// 	4, 210,
		// 	128,
		// 	0,
		// 	1, 0,
		// 	0, 0,
		// 	0, 0,
		// 	0, 0,
		// }

		response := []byte{}
		response = append(response, buf[:2]...)
		response = append(response, 0b10000000)
		response = append(response, 0b00000000)
		response = binary.BigEndian.AppendUint16(response, uint16(1)) // QDCOUNT
		response = binary.BigEndian.AppendUint16(response, 0x0000)    // ANCOUNT
		response = binary.BigEndian.AppendUint16(response, 0x0000)    // NSCOUNT
		response = binary.BigEndian.AppendUint16(response, 0x0000)    // ARCOUNT
		response = append(response, encodeDomain("codecrafters.io")...)
		response = binary.BigEndian.AppendUint16(response, uint16(1))
		response = binary.BigEndian.AppendUint16(response, uint16(1))

		// fmt.Println(len(receivedQuestion))

		// response = append(response, receivedQuestion...)

		fmt.Println(len(response))

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
