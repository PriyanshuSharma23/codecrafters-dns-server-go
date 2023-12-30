package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type Message struct {
	Header Header
}

// header stucture
/*
   	                              1  1  1  1  1  1
   	0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                      ID                       |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    QDCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    ANCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    NSCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
   |                    ARCOUNT                    |
   +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type Header struct {
	ID      uint16
	FLAGS   uint16
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

// func newHeader() *Header {
// 	return &Header{
// 		ID:      1234,
// 		FLAGS:   0b1000000000000000,
// 		QDCOUNT: 1,
// 		ANCOUNT: 0,
// 		NSCOUNT: 0,
// 		ARCOUNT: 0,
// 	}
// }

func newHeader() *Header {
	return &Header{
		ID:      1234,
		FLAGS:   0b1000000000000000,
		QDCOUNT: 1,
		ANCOUNT: 1,
		NSCOUNT: 0,
		ARCOUNT: 0,
	}
}

func (h *Header) toBytes() []byte {

	buf := make([]byte, 12)

	binary.BigEndian.PutUint16(buf[0:2], h.ID)
	binary.BigEndian.PutUint16(buf[2:4], h.FLAGS)
	binary.BigEndian.PutUint16(buf[4:6], h.QDCOUNT)
	binary.BigEndian.PutUint16(buf[6:8], h.ANCOUNT)
	binary.BigEndian.PutUint16(buf[8:10], h.NSCOUNT)
	binary.BigEndian.PutUint16(buf[10:12], h.ARCOUNT)

	return buf
}

/*
	                              1  1  1  1  1  1
	0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5

+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                                               |
/                     QNAME                     /
/                                               /
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     QTYPE                     |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
|                     QCLASS                    |
+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
*/
type Question struct {
	QNAME  string
	QTYPE  uint16
	QCLASS uint16
}

func labelEncoder(s string) []byte {
	buf := make([]byte, 0)
	for _, part := range strings.Split(s, ".") {
		buf = append(buf, byte(len(part)))
		buf = append(buf, part...)
	}
	buf = append(buf, '\x00')
	return buf
}

func newQuestion() *Question {
	return &Question{
		QNAME:  "codecrafters.io",
		QTYPE:  1,
		QCLASS: 1,
	}
}

func (q *Question) toBytes() []byte {
	buf := make([]byte, 0, len(q.QNAME)+4) // 4 => 2 bytes qclass 2bytes qtype
	buf = append(buf, labelEncoder(q.QNAME)...)

	buf = binary.BigEndian.AppendUint16(buf, q.QTYPE)
	buf = binary.BigEndian.AppendUint16(buf, q.QCLASS)

	return buf
}

//                                    1  1  1  1  1  1
//      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                                               |
//    /                                               /
//    /                      NAME                     /
//    |                                               |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                      TYPE                     |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                     CLASS                     |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                      TTL                      |
//    |                                               |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
//    |                   RDLENGTH                    |
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--|
//    /                     RDATA                     /
//    /                                               /
//    +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+

type ResourceRecord struct {
	NAME     string
	TYPE     uint16
	CLASS    uint16
	TTL      uint16
	RDLENGTH uint16
	RDATA    []byte
}

// func newResourceRecord() *ResourceRecord {
// 	return &ResourceRecord{
// 		NAME:     "codecrafters.io",
// 		TYPE:     1,
// 		CLASS:    1,
// 		TTL:      60,
// 		RDLENGTH: 4,
// 		RDATA:    []byte{8, 8, 8, 8},
// 	}
// }

func (rr *ResourceRecord) toBytes() []byte {
	buf := make([]byte, 0)

	buf = append(buf, labelEncoder(rr.NAME)...)
	buf = binary.BigEndian.AppendUint16(buf, rr.TYPE)
	buf = binary.BigEndian.AppendUint16(buf, rr.CLASS)
	buf = binary.BigEndian.AppendUint16(buf, rr.TTL)
	buf = binary.BigEndian.AppendUint16(buf, rr.RDLENGTH)
	buf = append(buf, rr.RDATA...)

	return buf
}

type Answer []ResourceRecord

func newAnswer() *Answer {
	return &Answer{
		ResourceRecord{
			NAME:     "codecrafters.io",
			TYPE:     1,
			CLASS:    1,
			TTL:      60,
			RDLENGTH: 4,
			RDATA:    []byte{8, 8, 8, 8},
		},
	}
}

func (a *Answer) toBytes() []byte {
	buf := make([]byte, 0)
	for _, rr := range *a {
		buf = append(buf, rr.toBytes()...)
	}
	return buf
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		response := make([]byte, 0)
		header := newHeader()
		response = header.toBytes()

		question := newQuestion()
		response = append(response, question.toBytes()...)

		answer := newAnswer()
		response = append(response, answer.toBytes()...)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
