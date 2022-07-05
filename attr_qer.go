package gtp5gnl

import (
	"github.com/khirono/go-nl"
)

const (
	QER_ID = iota + 3
	QER_GATE
	QER_MBR
	QER_GBR
	QER_CORR_ID
	QER_RQI
	QER_QFI
	QER_PPI
	QER_RCSR // deplicated
	QER_RELATED_TO_PDR
	QER_SEID
)

type QER struct {
	ID     uint32
	Gate   uint8
	MBR    MBR
	GBR    GBR
	CorrID uint32
	RQI    uint8
	QFI    uint8
	PPI    uint8
	PDRIDs []uint16
	SEID   *uint64
}

func DecodeQER(b []byte) (*QER, error) {
	qer := new(QER)
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		switch hdr.MaskedType() {
		case QER_ID:
			qer.ID = native.Uint32(b[n:])
		case QER_GATE:
			qer.Gate = b[n]
		case QER_MBR:
			mbr, err := DecodeMBR(b[n:])
			if err != nil {
				return nil, err
			}
			qer.MBR = mbr
		case QER_GBR:
			gbr, err := DecodeGBR(b[n:])
			if err != nil {
				return nil, err
			}
			qer.GBR = gbr
		case QER_CORR_ID:
			qer.CorrID = native.Uint32(b[n:])
		case QER_RQI:
			qer.RQI = b[n]
		case QER_QFI:
			qer.QFI = b[n]
		case QER_PPI:
			qer.PPI = b[n]
		case QER_RELATED_TO_PDR:
			d := b[n:hdr.Len]
			for len(d) > 0 {
				v := native.Uint16(d)
				qer.PDRIDs = append(qer.PDRIDs, v)
				d = d[2:]
			}
		case QER_SEID:
			v := native.Uint64(b[n:])
			qer.SEID = &v
		}
		b = b[hdr.Len.Align():]
	}
	return qer, nil
}

const (
	QER_MBR_UL_HIGH32 = iota + 1
	QER_MBR_UL_LOW8
	QER_MBR_DL_HIGH32
	QER_MBR_DL_LOW8
)

type MBR struct {
	ULHigh uint32
	ULLow  uint8
	DLHigh uint32
	DLLow  uint8
}

func DecodeMBR(b []byte) (MBR, error) {
	var mbr MBR
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return mbr, err
		}
		switch hdr.MaskedType() {
		case QER_MBR_UL_HIGH32:
			mbr.ULHigh = native.Uint32(b[n:])
		case QER_MBR_UL_LOW8:
			mbr.ULLow = b[n]
		case QER_MBR_DL_HIGH32:
			mbr.DLHigh = native.Uint32(b[n:])
		case QER_MBR_DL_LOW8:
			mbr.DLLow = b[n]
		}
		b = b[hdr.Len.Align():]
	}
	return mbr, nil
}

const (
	QER_GBR_UL_HIGH32 = iota + 1
	QER_GBR_UL_LOW8
	QER_GBR_DL_HIGH32
	QER_GBR_DL_LOW8
)

type GBR struct {
	ULHigh uint32
	ULLow  uint8
	DLHigh uint32
	DLLow  uint8
}

func DecodeGBR(b []byte) (GBR, error) {
	var gbr GBR
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return gbr, err
		}
		switch hdr.MaskedType() {
		case QER_GBR_UL_HIGH32:
			gbr.ULHigh = native.Uint32(b[n:])
		case QER_GBR_UL_LOW8:
			gbr.ULLow = b[n]
		case QER_GBR_DL_HIGH32:
			gbr.DLHigh = native.Uint32(b[n:])
		case QER_GBR_DL_LOW8:
			gbr.DLLow = b[n]
		}
		b = b[hdr.Len.Align():]
	}
	return gbr, nil
}
