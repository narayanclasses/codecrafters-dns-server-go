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
		answerSection := []byte{}
		i := 12
		qcount := int(buf[5])
		j := 0
		for j < qcount {
			tempi := -1
			for i < len(buf) {
				length := int(buf[i])
				if length == 192 {
					tempi = i + 1
					i = int(buf[i+1])
					continue
				}
				receivedQuestion = append(receivedQuestion, buf[i])
				answerSection = append(answerSection, buf[i])
				if length == 0 {
					break
				}
				i++ // move to the start of the segment
				receivedQuestion = append(receivedQuestion, buf[i:i+length]...)
				answerSection = append(answerSection, buf[i:i+length]...)
				i += length // move to the next length prefix
			}
			if tempi != -1 {
				i = tempi
			}
			receivedQuestion = append(receivedQuestion, buf[i+1], buf[i+2], buf[i+3], buf[i+4])
			answerSection = append(answerSection, buf[i+1], buf[i+2], buf[i+3], buf[i+4])
			answerSection = append(answerSection,
				0, 0, 0, 0,
				0, 4,
				0x08, 0x08, 0x08, 0x08,
			)
			i = i + 5
			j++
		}

		opcode := buf[2] & (121)
		rcode := make([]byte, 1)
		if opcode != 0 {
			rcode[0] = 4
		}

		response := []byte{}
		response = append(response,
			buf[0], buf[1],
			(buf[2]&(121))|(128),
			rcode[0],
			0, byte(qcount),
			0, byte(qcount),
			0, 0,
			0, 0)

		response = append(response, receivedQuestion...)
		response = append(response, answerSection...)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
