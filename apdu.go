package apdu

import (
	"encoding/hex"
	"errors"
	"strings"
)

// APDU represents ISO-7816 APDU command
type APDU struct {
	Cla  byte
	Ins  byte
	P1   byte
	P2   byte
	Data []byte
	Le   byte
}

var (
	errInvalidDataLength = errors.New("Lc byte does not match data length")
	errInvalidLength     = errors.New("invalid length")
)

// Bytes return APDU encoded into a byte slice
func (a APDU) Bytes() []byte {
	command := []byte{a.Cla, a.Ins, a.P1, a.P2}

	lc := byte(len(a.Data))
	if lc > 0 {
		command = append(command, lc)
	}

	command = append(command, a.Data...)
	return append(command, a.Le)
}

// String implements fmt.Stringer
func (a APDU) String() string {
	return hex.EncodeToString(a.Bytes())
}

var spaceReplacer = strings.NewReplacer(" ", "", "\t", "", "\n", "", "\r", "")

// FromString creates APDU command from hexadecimal string
func FromString(s string) (*APDU, error) {
	b, err := hex.DecodeString(spaceReplacer.Replace(s))
	if err != nil {
		return nil, err
	}

	return FromBytes(b)
}

// FromBytes creates APDU command from the byte slice
func FromBytes(b []byte) (*APDU, error) {
	if len(b) < 4 {
		return nil, errInvalidLength
	}

	apdu := APDU{Cla: b[0], Ins: b[1], P1: b[2], P2: b[3]}

	switch len(b) {
	case 4:
		return &apdu, nil
	case 5:
		apdu.Le = b[4]
		return &apdu, nil
	}

	lc := int(b[4])
	if len(b)-5 < lc {
		return nil, errInvalidDataLength
	}

	if len(b) > 6+lc {
		return nil, errInvalidLength
	}

	apdu.Data = b[5 : 5+lc]

	if len(b) > 5+lc {
		apdu.Le = b[5+lc]
	}

	return &apdu, nil
}

// MustFromString creates APDU command from hexadecimal string and panics on error
func MustFromString(s string) APDU {
	apdu, err := FromString(s)
	if err != nil {
		panic(err)
	}

	return *apdu
}

// Select returns APDU command for selecting requested application ID
func Select(aid []byte) APDU {
	return APDU{
		Cla:  0x00,
		Ins:  0xa4,
		P1:   0x04,
		P2:   0x00,
		Data: aid,
	}
}
