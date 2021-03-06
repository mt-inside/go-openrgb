package wire

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-logr/logr"
)

// Client represents an open connection to an OpenRGB server.
type Client struct {
	log  logr.Logger
	sock net.Conn
}

// NewClient returns a new Client that is connected to the server at addr. If a
// connection can't be established, or there's a protocol or other problem, an
// error is returned and Client in nil. userAgent is the string given to the
// server to identify this client software; it does not need to be unique.
func NewClient(log logr.Logger, addr, userAgent string) (*Client, error) {
	c := &Client{log: log.WithName("client")}

	sock, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Couldn't connect to OpenRGB server at %s: %w", addr, err)
	}

	c.sock = sock

	err = c.sendCommand(0, cmdGetProtocolVersion, []byte{})
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("Couldn't get protocol version: %w", err)
	}
	body, err := c.readMessage()
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("Couldn't get protocol version: %w", err)
	}
	offset := 0
	protoVer := extractUint32(body, &offset)
	if protoVer != knownProtoVer {
		c.Close()
		return nil, fmt.Errorf("Server protocol version: %d; we support: %d", protoVer, knownProtoVer)
	}

	err = c.sendCommand(0, cmdSetClientName, []byte(userAgent+"\000")) // Not a headed bstring, but does need the terminator
	if err != nil {
		c.Close()
		return nil, fmt.Errorf("Couldn't set client name: %w", err)
	}

	c.log.V(2).Info("Connected", "user agent", userAgent, "protocol version", protoVer)

	return c, nil
}

// Close terminates the connection to the OpenRGB server, by amougst other
// things closing the TCP stream.
func (c *Client) Close() error {
	c.log.V(2).Info("Disconnected")

	return c.sock.Close()
}

func (c *Client) sendCommand(deviceID uint32, commandID Command, body []byte) error {
	if body == nil {
		panic("body mustn't be nil; please use empty object")
	}

	header := encodeHeader(uint32(deviceID), uint32(commandID), uint32(len(body)))

	c.log.V(2).Info(">> header", "device", deviceID, "command", commandID)
	if c.log.V(4).Enabled() {
		spew.Dump(header)
	}
	_, err := c.sock.Write(header)
	if err != nil {
		return fmt.Errorf("Couldn't send message header: %w", err)
	}

	c.log.V(2).Info(">> body", "len", len(body))
	if c.log.V(4).Enabled() {
		spew.Dump(body)
	}
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
	_, _, bodyLen := decodeHeader(headerBytes)
	c.log.V(2).Info("<< header")
	if c.log.V(4).Enabled() {
		spew.Dump(headerBytes)
	}

	bodyBytes := make([]byte, bodyLen)
	_, err = c.sock.Read(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("Couldn't read from server: %w", err)
	}
	c.log.V(2).Info("<< body", "len", len(bodyBytes))
	if c.log.V(4).Enabled() {
		spew.Dump(bodyBytes)
	}
	return bodyBytes, nil
}

const headerMagic = "ORGB"
const headerLen = 16

func encodeHeader(deviceID, commandID, bodyLen uint32) []byte {
	header := [headerLen]byte{}

	copy(header[:], headerMagic)
	binary.LittleEndian.PutUint32(header[4:8], deviceID)
	binary.LittleEndian.PutUint32(header[8:12], commandID)
	binary.LittleEndian.PutUint32(header[12:16], bodyLen)

	return header[:]
}

func decodeHeader(header []byte) (deviceID, commandID, bodyLen uint32) {
	l := len(header)
	if l != headerLen {
		panic(fmt.Sprintf("Header length incorrect. Expected: %d, got: %d", headerLen, l))
	}

	magic := string(header[0:4])
	if magic != headerMagic {
		// Too deep to start returning errors-as-values
		panic(fmt.Sprintf("Header missing magic. Expected: %s, got: %s", headerMagic, magic))
	}

	return binary.LittleEndian.Uint32(header[4:8]),
		binary.LittleEndian.Uint32(header[8:12]),
		binary.LittleEndian.Uint32(header[12:16])
}
