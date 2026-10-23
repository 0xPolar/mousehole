package net

type IPv4Info struct {
	HeaderLength       int
	TotalLength        int
	SourceAddress      [4]byte
	DestinationAddress [4]byte
}
