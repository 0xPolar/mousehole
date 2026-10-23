package net

import (
	"encoding/binary"
	"fmt"
)

type IPv4Info struct {
	HeaderLength       int
	TotalLength        int
	SourceAddress      [4]byte
	DestinationAddress [4]byte
}

func bitshift(raw byte) byte {
	return 0xF0
}

func ParseIPv4(packet []byte) (IPv4Info, error) {
	if len(packet) < 20 {
		return IPv4Info{}, fmt.Errorf("packet shorter than minimum packet size")
	}

	packetVersion := packet[0] >> 4
	if packetVersion != 4 {
		return IPv4Info{}, fmt.Errorf("wrong packet version only allowed packet version if 4")
	}
	headerLength := int(packet[0]&0x0F) * 4
	if headerLength < 20 {
		return IPv4Info{}, fmt.Errorf("header length too short")
	}

	packetLen := binary.BigEndian.Uint16(packet[2:4])
	if int(packetLen) < headerLength {
		return IPv4Info{}, fmt.Errorf("header length londer than total packet length")
	}

	if int(packetLen) != len(packet) {
		return IPv4Info{}, fmt.Errorf("packet supplied and real size and not the same")
	}
	packetSource := packet[12:16]
	packetDestionation := packet[16:20]
	return IPv4Info{
		HeaderLength:       headerLength,
		TotalLength:        int(packetLen),
		SourceAddress:      [4]byte(packetSource),
		DestinationAddress: [4]byte(packetDestionation),
	}, nil
}
