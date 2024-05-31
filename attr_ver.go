package gtp5gnl

import (
	"bytes"

	"github.com/khirono/go-nl"
)

func DecodeVersion(b []byte) (string, error) {
	hdr, n, err := nl.DecodeAttrHdr(b)
	if err != nil {
		return "", err
	}

	attrLen := int(hdr.Len)
	ver := string(bytes.Trim(b[n:attrLen], "\x00"))

	return ver, err
}
