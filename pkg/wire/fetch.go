package wire

import (
	"encoding/binary"
	"fmt"
)

func FetchDevices(c *Client) ([]*Device, error) {
	deviceCount, err := fetchDeviceCount(c)
	if err != nil {
		return []*Device{}, fmt.Errorf("Couldn't fetch Device count: %w", err)
	}

	ds := make([]*Device, deviceCount)
	for i := uint32(0); i < deviceCount; i++ {
		ds[i], err = fetchDevice(c, i)
		if err != nil {
			return []*Device{}, fmt.Errorf("Couldn't fetch Device %d: %w", i, err)
		}
	}

	return ds, nil
}

/* API is so shit.
* Everything has a length header.
* Devices' length is fetched by a separate command.
* Then each device by a command.
* Within a device, zones etc aren't API commands, they're packed into the binary blob, with thier length preceeding them */
func fetchDeviceCount(c *Client) (uint32, error) {
	if err := c.sendCommand(0, cmdGetDevCnt, []byte{}); err != nil {
		return 0, fmt.Errorf("Couldn't fetch Device count: %w", err)
	}

	body, err := c.readMessage()
	if err != nil {
		return 0, fmt.Errorf("Couldn't fetch Device count: %w", err)
	}

	deviceCount := binary.LittleEndian.Uint32(body)

	return deviceCount, nil
}

func fetchDevice(c *Client, i uint32) (*Device, error) {
	if err := c.sendCommand(uint32(i), cmdGetDevData, []byte{}); err != nil {
		return &Device{}, fmt.Errorf("Couldn't fetch Device %d: %w", i, err)
	}

	body, err := c.readMessage()
	if err != nil {
		return &Device{}, fmt.Errorf("Couldn't fetch Device %d: %w", i, err)
	}

	device := extractDevice(body, i)

	return device, nil
}
