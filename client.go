package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"
)

type Client struct {
	sock net.Conn
}

//go:generate stringer -type=Command
type Command int

const (
	cmdGetDevCnt              = 0
	cmdGetDevData             = 1
	cmdSetClientName  Command = 50
	cmdUpdateLEDs             = 1050
	cmdUpdateZoneLEDs         = 1051
	cmdSetCustomMode          = 1100
)

func NewClient(addr, userAgent string) (*Client, error) {
	sock, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to OpenRGB server at %s: %w", addr, err)
	}

	c := &Client{sock: sock}

	err = c.sendCommand(0, cmdSetClientName, []byte(userAgent+" (go-openrgb)"))
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("Couldn't set client name: %w", err)
	}

	return c, nil
}

func (c *Client) Close() error {
	return c.sock.Close()
}

func (c *Client) sendCommand(deviceID uint32, commandID Command, body []byte) error {
	if body == nil {
		panic("body mustn't be nil; please use empty object")
	}

	header := encodeHeader(uint32(deviceID), uint32(commandID), uint32(len(body)))

	fmt.Printf(">>> ")
	spew.Dump(header)
	_, err := c.sock.Write(header)
	if err != nil {
		return fmt.Errorf("Couldn't send message header: %w", err)
	}

	fmt.Printf(">++ ")
	spew.Dump(body)
	_, err = c.sock.Write(body)
	if err != nil {
		return fmt.Errorf("Couldn't send message body: %w", err)
	}

	return nil
}

func (c *Client) readMessage() (body []byte, err error) {
	headerBytes := make([]byte, headerLen)
	_, err = c.sock.Read(headerBytes)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read from server: %w", err)
	}
	fmt.Printf("<<< ")
	spew.Dump(headerBytes)
	_, _, bodyLen := decodeHeader(headerBytes)

	bodyBytes := make([]byte, bodyLen)
	_, err = c.sock.Read(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read from server: %w", err)
	}
	fmt.Printf("<++ ")
	spew.Dump(bodyBytes)
	return bodyBytes, nil
}

const headerMagic = "ORGB"
const headerLen = 16

/* If this turns into a bottleneck, use a byte *array*, calculate offsets, copy() */
func encodeHeader(deviceID, commandID, bodyLen uint32) []byte {
	header := new(bytes.Buffer)

	header.Write([]byte(headerMagic))
	binary.Write(header, binary.LittleEndian, deviceID)
	binary.Write(header, binary.LittleEndian, commandID)
	binary.Write(header, binary.LittleEndian, bodyLen)

	if header.Len() != headerLen {
		panic(fmt.Sprintf("Assertion failed: header len %d should be %d", header.Len(), headerLen))
	}

	return header.Bytes()
}

func decodeHeader(header []byte) (commandID, deviceID, bodyLen uint32) {
	magic := string(header[0:4])
	if magic != headerMagic {
		// Too deep to start returning errors-as-values
		panic(fmt.Sprintf("Header missing magic. Expected: %s, got: %s", headerMagic, header))
	}

	return binary.LittleEndian.Uint32(header[4:8]),
		binary.LittleEndian.Uint32(header[8:12]),
		binary.LittleEndian.Uint32(header[12:16])
}
