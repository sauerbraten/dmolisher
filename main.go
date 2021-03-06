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
	stamp := Stamp{}
	err := binary.Read(stream, binary.LittleEndian, &stamp)
	if err != nil {
		return nil, fmt.Errorf("reading stamp: %w", err)
	}
	return &stamp, nil
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

type DemoHeader struct {
	Magic           [16]byte
	FileVersion     int32
	ProtocolVersion int32
}

func readDemoHeader(stream io.Reader) (*DemoHeader, error) {
	hdr := DemoHeader{}
	err := binary.Read(stream, binary.LittleEndian, &hdr)
	if err != nil {
		return nil, fmt.Errorf("reading demo header: %w", err)
	}
	if string(hdr.Magic[:]) != "SAUERBRATEN_DEMO" {
		return nil, fmt.Errorf("reading demo header: wrong magic (not a demo file?)")
	}
	return &hdr, nil
}

var (
	filterChannel = flag.Int("channel", -1, "print only packets sent on channel (0/1/2)")
	printHex      = flag.Bool("hex", false, "print data bytes in hexadecimal instead of decimal")
	printVersions = flag.Bool("versions", false, "print file and protocol versions")
)

func init() {
	flag.CommandLine.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Reads uncompressed demo file on stdin and emits CSV on stdout.\n")
		flag.Usage()
	}
}

func main() {
	flag.Parse()

	hdr, err := readDemoHeader(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing demo stream: %v\n", err)
		os.Exit(1)
	}
	if *printVersions {
		fmt.Printf("# file version: %d, protocol version: %d\n", hdr.FileVersion, hdr.ProtocolVersion)
	}
	if hdr.FileVersion != 1 {
		fmt.Fprintln(os.Stderr, "error: unsupported file version (only version 1 is supported)")
		os.Exit(2)
	}

	fmt.Println("gamemillis, channel, data length, data")

	stamp, data, err := readPacket(os.Stdin)
	for err == nil {
		if *filterChannel == -1 || int(stamp.Channel) == *filterChannel {
			fmt.Printf("%6d, %d, %2d,", stamp.Time, stamp.Channel, stamp.Length)
			for _, b := range data {
				if *printHex {
					fmt.Printf(" %02x", b)
				} else {
					fmt.Printf(" %d", b)
				}
			}
			fmt.Println()
		}

		stamp, data, err = readPacket(os.Stdin)
	}

	if !errors.Is(err, io.EOF) {
		fmt.Fprintf(os.Stderr, "error parsing demo stream: %v\n", err)
		os.Exit(3)
	}
}
