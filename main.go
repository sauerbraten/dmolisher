package main

import (
	"encoding/binary"
	"errors"
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

func readDemoHeader(stream io.Reader) error {
	buf := make([]byte, 24)
	_, err := stream.Read(buf)
	if err != nil {
		return fmt.Errorf("reading Sauerbraten demo header: %w", err)
	}
	// todo: actually verify header
	return nil
}

func main() {
	err := readDemoHeader(os.Stdin)
	if err != nil {
		fmt.Printf("error parsing demo stream: %v\n", err)
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
