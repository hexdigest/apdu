package apdu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromBytes(t *testing.T) {
	tests := []struct {
		name string
		b    []byte

		want1      *APDU
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name:    "length is too short",
			b:       []byte{0x11, 0x22},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errInvalidLength, err)
			},
		},
		{
			name:    "4 byte",
			b:       []byte{0x11, 0x22, 0x33, 0x44},
			want1:   &APDU{Cla: 0x11, Ins: 0x22, P1: 0x33, P2: 0x44},
			wantErr: false,
		},
		{
			name:    "5 byte (with Le)",
			b:       []byte{0x11, 0x22, 0x33, 0x44, 0x55},
			want1:   &APDU{Cla: 0x11, Ins: 0x22, P1: 0x33, P2: 0x44, Le: 0x55},
			wantErr: false,
		},
		{
			name:    "with data but without Le",
			b:       []byte{0x11, 0x22, 0x33, 0x44, 0x02, 0x77, 0x88},
			want1:   &APDU{Cla: 0x11, Ins: 0x22, P1: 0x33, P2: 0x44, Data: []byte{0x77, 0x88}},
			wantErr: false,
		},
		{
			name:    "with data and with Le",
			b:       []byte{0x11, 0x22, 0x33, 0x44, 0x02, 0x77, 0x88, 0x99},
			want1:   &APDU{Cla: 0x11, Ins: 0x22, P1: 0x33, P2: 0x44, Le: 0x99, Data: []byte{0x77, 0x88}},
			wantErr: false,
		},
		{
			name:    "data is too short",
			b:       []byte{0x11, 0x22, 0x33, 0x44, 0x03, 0x77},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errInvalidDataLength, err)
			},
		},
		{
			name:    "data is too long",
			b:       []byte{0x11, 0x22, 0x33, 0x44, 0x02, 0x77, 0x88, 0x99, 0xFF},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errInvalidLength, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := FromBytes(tt.b)

			assert.Equal(t, tt.want1, got1, "FromBytes returned unexpected result")

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

func TestFromString(t *testing.T) {
	tests := []struct {
		name string
		s    string

		want1      *APDU
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name:    "parse error",
			s:       "000W",
			want1:   nil,
			wantErr: true,
		},
		{
			name:    "success",
			s:       "11223344",
			want1:   &APDU{Cla: 0x11, Ins: 0x22, P1: 0x33, P2: 0x44},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1, err := FromString(tt.s)

			assert.Equal(t, tt.want1, got1, "FromString returned unexpected result")

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

func TestMustFromString(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		assert.Panics(t, func() {
			MustFromString("panic is good for you")
		})
	})

	t.Run("success", func(t *testing.T) {
		apdu := MustFromString("11223344")
		assert.Equal(t, APDU{Cla: 0x11, Ins: 0x22, P1: 0x33, P2: 0x44}, apdu)
	})
}

func TestSelect(t *testing.T) {
	apdu := Select([]byte{0x11, 0x22})
	assert.Equal(t, APDU{Cla: 0x00, Ins: 0xa4, P1: 0x04, P2: 0x00, Data: []byte{0x11, 0x22}}, apdu)
}

func TestAPDU_Bytes(t *testing.T) {
	b := APDU{Cla: 0x11, Ins: 0x22, P1: 0x33, P2: 0x44, Data: []byte{0x88, 0x99}, Le: 0x77}.Bytes()
	assert.Equal(t, []byte{0x11, 0x22, 0x33, 0x44, 0x02, 0x88, 0x99, 0x77}, b)
}

func TestAPDU_String(t *testing.T) {
	s := APDU{Cla: 0x11, Ins: 0x22, P1: 0x33, P2: 0x44, Data: []byte{0x88, 0x99}, Le: 0x77}.String()
	assert.Equal(t, "1122334402889977", s)
}
