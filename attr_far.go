package gtp5gnl

import (
	"net"

	"github.com/khirono/go-nl"
)

const (
	FAR_ID = iota + 3
	FAR_APPLY_ACTION
	FAR_FORWARDING_PARAMETER
	FAR_RELATED_TO_PDR
	FAR_UL_OR_DL
	FAR_SEID
	FAR_BAR_ID
)

type FAR struct {
	ID       uint32
	Action   uint16
	Param    *ForwardParam
	PDRIDs   []uint16
	BARID    *uint8
	SEID     *uint64
	UL_OR_DL uint8
}

func DecodeFAR(b []byte) (*FAR, error) {
	far := new(FAR)
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case FAR_ID:
			far.ID = native.Uint32(b[n:attrLen])
		case FAR_APPLY_ACTION:
			far.Action = native.Uint16(b[n:attrLen])
		case FAR_FORWARDING_PARAMETER:
			param, err := DecodeForwardParam(b[n:attrLen])
			if err != nil {
				return nil, err
			}
			far.Param = &param
		case FAR_RELATED_TO_PDR:
			d := b[n:attrLen]
			for len(d) > 0 {
				v := native.Uint16(d)
				far.PDRIDs = append(far.PDRIDs, v)
				d = d[2:]
			}
		case FAR_BAR_ID:
			v := b[n]
			far.BARID = &v
		case FAR_SEID:
			v := native.Uint64(b[n:attrLen])
			far.SEID = &v
		case FAR_UL_OR_DL:
			far.UL_OR_DL = b[n]
		}
		b = b[hdr.Len.Align():]
	}
	return far, nil
}

const (
	FORWARDING_PARAMETER_OUTER_HEADER_CREATION = iota + 1
	FORWARDING_PARAMETER_FORWARDING_POLICY
	FORWARDING_PARAMETER_PFCPSM_REQ_FLAGS
	FORWARDING_PARAMETER_TOS_TC
)

type ForwardParam struct {
	Creation *HeaderCreation
	Policy   *string
	TosTc    uint8
}

func DecodeForwardParam(b []byte) (ForwardParam, error) {
	var param ForwardParam
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return param, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case FORWARDING_PARAMETER_OUTER_HEADER_CREATION:
			hc, err := DecodeHeaderCreation(b[n:attrLen])
			if err != nil {
				return param, err
			}
			param.Creation = &hc
		case FORWARDING_PARAMETER_FORWARDING_POLICY:
			s, _, _ := nl.DecodeAttrString(b[n:attrLen])
			param.Policy = &s
		case FORWARDING_PARAMETER_TOS_TC:
			param.TosTc = b[n]
		}
		b = b[hdr.Len.Align():]
	}
	return param, nil
}

const (
	OUTER_HEADER_CREATION_DESCRIPTION = iota + 1
	OUTER_HEADER_CREATION_O_TEID
	OUTER_HEADER_CREATION_PEER_ADDR_IPV4
	OUTER_HEADER_CREATION_PORT
)

type HeaderCreation struct {
	Desc     uint16
	TEID     uint32
	PeerAddr net.IP
	Port     uint16
}

func DecodeHeaderCreation(b []byte) (HeaderCreation, error) {
	var hc HeaderCreation
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return hc, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case OUTER_HEADER_CREATION_DESCRIPTION:
			hc.Desc = native.Uint16(b[n:attrLen])
		case OUTER_HEADER_CREATION_O_TEID:
			hc.TEID = native.Uint32(b[n:attrLen])
		case OUTER_HEADER_CREATION_PEER_ADDR_IPV4:
			hc.PeerAddr = make([]byte, 4)
			copy(hc.PeerAddr, b[n:n+4])
		case OUTER_HEADER_CREATION_PORT:
			hc.Port = native.Uint16(b[n:attrLen])
		}
		b = b[hdr.Len.Align():]
	}
	return hc, nil
}
