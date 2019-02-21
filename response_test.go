package apdu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	tests := []struct {
		name string
		b    []byte

		want1      *Response
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name:    "invalid length",
			b:       []byte{0x90},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errInvalidLength, err)
			},
		},

		{
			name:    "status word only",
			b:       []byte{0x90, 0x00},
			want1:   &Response{StatusWord: StatusWord{0x90, 0x00}},
			wantErr: false,
		},
		{
			name: "with data",
			b:    []byte{0x11, 0x22, 0x90, 0x00},
			want1: &Response{
				StatusWord: StatusWord{0x90, 0x00},
				Data:       []byte{0x11, 0x22},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := NewResponse(tt.b)

			assert.Equal(t, tt.want1, got1, "NewResponse returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestParseResponse(t *testing.T) {
	tests := []struct {
		name string
		b    []byte

		want1      []byte
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name:    "invalid length",
			b:       []byte{0x00},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errInvalidLength, err)
			},
		},
		{
			name:    "not 0x900",
			b:       []byte{0x11, 0x22, 0x33, 0x9F, 0x00},
			want1:   []byte{0x11, 0x22, 0x33},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, StatusWord{0x9f, 0x00}, err)
			},
		},
		{
			name:    "success",
			b:       []byte{0x11, 0x22, 0x33, 0x90, 0x00},
			want1:   []byte{0x11, 0x22, 0x33},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := ParseResponse(tt.b)

			assert.Equal(t, tt.want1, got1, "ParseResponse returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestStatusWord_Error(t *testing.T) {
	tests := []struct {
		name  string
		sw    StatusWord
		want1 string
	}{
		{
			name:  "0x6A82",
			sw:    StatusWord{0x6a, 0x82},
			want1: "6A82 (File not found)",
		},
		{
			name:  "unknown",
			sw:    StatusWord{0x11, 0x11},
			want1: "1111 (Unknown)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := tt.sw.Error()
			assert.Equal(t, tt.want1, got1, "StatusWord.Error returned unexpected result")
		})
	}
}

func TestStatusWord_IsError(t *testing.T) {
	tests := []struct {
		name  string
		sw    StatusWord
		want1 bool
	}{
		{
			name:  "0x9000",
			sw:    StatusWord{0x90, 0x00},
			want1: false,
		},
		{
			name:  "0x9f32",
			sw:    StatusWord{0x9f, 0x11},
			want1: false,
		},
		{
			name:  "0x6a82",
			sw:    StatusWord{0x6a, 0x82},
			want1: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := tt.sw.IsError()
			assert.Equal(t, tt.want1, got1, "StatusWord.IsError returned unexpected result")
		})
	}
}
