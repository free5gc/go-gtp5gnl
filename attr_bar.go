package gtp5gnl

import (
	"github.com/khirono/go-nl"
)

const (
	BAR_ID = iota + 3
	BAR_DOWNLINK_DATA_NOTIFICATION_DELAY
	BAR_BUFFERING_PACKETS_COUNT
	BAR_SEID
)

type BAR struct {
	ID    uint8
	Delay *uint8
	Count *uint16
	SEID  *uint64
}

func DecodeBAR(b []byte) (*BAR, error) {
	bar := new(BAR)
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case BAR_ID:
			bar.ID = b[n]
		case BAR_DOWNLINK_DATA_NOTIFICATION_DELAY:
			v := b[n]
			bar.Delay = &v
		case BAR_BUFFERING_PACKETS_COUNT:
			v := native.Uint16(b[n:attrLen])
			bar.Count = &v
		case BAR_SEID:
			v := native.Uint64(b[n:attrLen])
			bar.SEID = &v
		}
		b = b[hdr.Len.Align():]
	}
	return bar, nil
}
