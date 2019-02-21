package apdu

import "fmt"

// Response to APDU command from the target
type Response struct {
	Data       []byte
	StatusWord StatusWord
}

// StatusWord is a mandatory part of the response to APDU command
type StatusWord [2]byte

// IsError returns true if the status word is not 0x9fXX or 0x9000
func (sw StatusWord) IsError() bool {
	return !(sw[0] == 0x9f || (sw[0] == 0x90 && sw[1] == 0x00))
}

// Error implements error interface
func (sw StatusWord) Error() string {
	description, ok := statusWords[uint16(sw[0])<<8|uint16(sw[1])]
	if !ok {
		description = "Unknown"
	}

	return fmt.Sprintf("%X (%s)", [2]byte(sw), description)
}

// NewResponse parses byte slice and returns Response struct or error
func NewResponse(b []byte) (*Response, error) {
	if len(b) < 2 {
		return nil, errInvalidLength
	}

	response := Response{
		StatusWord: [2]byte{b[len(b)-2], b[len(b)-1]},
	}

	if len(b) > 2 {
		response.Data = make([]byte, len(b)-2)
		copy(response.Data, b[0:len(b)-2])
	}

	return &response, nil
}

// ParseResponse parses response and returns response data and/or error
// if the status word is not 0x9000
func ParseResponse(b []byte) (data []byte, err error) {
	r, err := NewResponse(b)
	if err != nil {
		return nil, err
	}

	if r.StatusWord[0] != 0x90 || r.StatusWord[1] != 0x00 {
		return r.Data, r.StatusWord
	}

	return r.Data, nil
}
