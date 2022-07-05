package gtp5gnl

import (
	"github.com/khirono/go-nl"
)

const (
	URR_ID = iota + 3
	URR_MEASUREMENT_METHOD
	URR_REPORTING_TRIGGER
	URR_MEASUREMENT_PERIOD
	URR_MEASUREMENT_INFO
	URR_SEQ
	URR_SEID
)

type URR struct {
	ID      uint32
	Method  uint64
	Trigger uint64
	Period  *uint64
	Info    *uint64
	Seq     *uint64
	SEID    *uint64
}

func DecodeURR(b []byte) (*URR, error) {
	urr := new(URR)
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		switch hdr.MaskedType() {
		case URR_ID:
			urr.ID = native.Uint32(b[n:])
		case URR_MEASUREMENT_METHOD:
			urr.Method = native.Uint64(b[n:])
		case URR_REPORTING_TRIGGER:
			urr.Trigger = native.Uint64(b[n:])
		case URR_MEASUREMENT_PERIOD:
			v := native.Uint64(b[n:])
			urr.Period = &v
		case URR_MEASUREMENT_INFO:
			v := native.Uint64(b[n:])
			urr.Info = &v
		case URR_SEQ:
			v := native.Uint64(b[n:])
			urr.Seq = &v
		case URR_SEID:
			v := native.Uint64(b[n:])
			urr.SEID = &v
		}
		b = b[hdr.Len.Align():]
	}
	return urr, nil
}
