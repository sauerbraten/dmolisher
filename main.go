package main

import (
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type Stamp struct {
	Time    int32
	Channel int32
	Length  int32
}

func readStamp(stream io.Reader) (*Stamp, error) {
	buf := make([]int32, 3)
	err := binary.Read(stream, binary.LittleEndian, buf)
	if err != nil {
		return nil, fmt.Errorf("reading stamp: %w", err)
	}
	return &Stamp{
		Time:    buf[0],
		Channel: buf[1],
		Length:  buf[2],
	}, nil
}

func readPacket(stream io.Reader) (*Stamp, []byte, error) {
	stamp, err := readStamp(stream)
	if err != nil {
		return nil, nil, err
	}

	buf := make([]byte, stamp.Length)
	_, err = stream.Read(buf)
	if err != nil {
		return nil, nil, fmt.Errorf("reading packet data: %w", err)
	}

	return stamp, buf, nil
}

func readDemoHeader(stream io.Reader) (int32, int32, error) {
	magic := make([]byte, 16)
	_, err := stream.Read(magic)
	if err != nil {
		return -1, -1, fmt.Errorf("reading demo header: reading magic: %w", err)
	}
	if string(magic) != "SAUERBRATEN_DEMO" {
		return -1, -1, fmt.Errorf("reading demo header: wrong magic (probably not a demo file)")
	}

	versions := make([]int32, 2)
	err = binary.Read(stream, binary.LittleEndian, versions)
	if err != nil {
		return -1, -1, fmt.Errorf("reading demo header: reading file and protocol versions: %w", err)
	}

	return versions[0], versions[1], nil
}

var printVersions = flag.Bool("versions", false, "print file and protocol versions")

func main() {
	flag.Parse()

	fileVersion, protocolVersion, err := readDemoHeader(os.Stdin)
	if err != nil {
		fmt.Printf("error parsing demo stream: %v\n", err)
	}
	if *printVersions {
		fmt.Printf("# file version: %d, protocol version: %d\n", fileVersion, protocolVersion)
	}
	if fileVersion != 1 {
		fmt.Fprintf(os.Stderr, "error: unsupported file version (only version 1 is supported)")
		os.Exit(1)
	}

	fmt.Println("gamemillis, channel, data length, data (bytes in decimal)")
	stamp, data, err := readPacket(os.Stdin)
	for err == nil {
		fmt.Printf("%d, %d, %d, %v\n", stamp.Time, stamp.Channel, stamp.Length, data)
		stamp, data, err = readPacket(os.Stdin)
	}

	if !errors.Is(err, io.EOF) {
		fmt.Printf("error parsing demo stream: %v\n", err)
	}
}
