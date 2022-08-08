package gtp5gnl

import (
	"bytes"

	"github.com/khirono/go-nl"
)

func DecodeVersion(b []byte) (string, error) {
	_, n, err := nl.DecodeAttrHdr(b)
	if err != nil {
		return "", err
	}

	ver := string(bytes.Trim(b[n:], "\x00"))

	return ver, err
}
