package types

import (
	"encoding/hex"
	"fmt"
)

var rawPollMessage = []byte{
	0x36, 0x38, 0x30, 0x30, 0x31, 0x30, 0x36, 0x38, 0x31, 0x30, 0x37, 0x37,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // inverter serial
	0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x39, 0x66, 0x31, 0x36,
}

func NewPollMessage(sn string) ([]byte, error) {
	id, err := hex.DecodeString(sn)
	if err != nil {
		return nil, err
	}

	if len(sn) != 8 || len(id) != 4 {
		return nil, fmt.Errorf("illegal inverter serial number: %s", sn)
	}

	msg := make([]byte, 32)
	copy(msg, rawPollMessage)
	copy(msg[12:20], []byte(sn))

	return msg, nil
}
