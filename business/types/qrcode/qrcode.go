// Package qrcode contains types related to QR codes.
package qrcode

import "github.com/yeqown/go-qrcode/v2"

// QRCode represents a QR code.
type QRCode struct {
	value string
}

// String returns the value of the QR code.
func (q QRCode) String() string {
	return q.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (q QRCode) MarshalText() ([]byte, error) {
	return []byte(q.value), nil
}

// NewQRCode creates a new QR code from a URL.
// TODO: implement this
func NewQRCode(value string) QRCode {
	_, err := qrcode.New(value)
	if err != nil {
		panic(err)
	}

	return QRCode{value: value}
}
