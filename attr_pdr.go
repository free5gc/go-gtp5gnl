package gtp5gnl

import (
	"log"
	"net"

	"github.com/khirono/go-nl"
)

const (
	PDR_ID = iota + 3
	PDR_PRECEDENCE
	PDR_PDI
	PDR_OUTER_HEADER_REMOVAL
	PDR_FAR_ID
	PDR_ROLE_ADDR_IPV4
	PDR_UNIX_SOCKET_PATH
	PDR_QER_ID
	PDR_SEID
	PDR_URR_ID
)

type PDR struct {
	ID              uint16
	Precedence      *uint32
	PDI             *PDI
	OuterHdrRemoval *uint8
	FARID           *uint32
	QERID           []uint32
	URRID           []uint32
	SEID            *uint64
}

func DecodePDR(b []byte) (*PDR, error) {
	pdr := new(PDR)
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case PDR_ID:
			pdr.ID = native.Uint16(b[n:attrLen])
		case PDR_PRECEDENCE:
			v := native.Uint32(b[n:attrLen])
			pdr.Precedence = &v
		case PDR_PDI:
			pdi, err := DecodePDI(b[n:attrLen])
			if err != nil {
				return nil, err
			}
			pdr.PDI = &pdi
		case PDR_OUTER_HEADER_REMOVAL:
			v := uint8(b[n])
			pdr.OuterHdrRemoval = &v
		case PDR_FAR_ID:
			v := native.Uint32(b[n:attrLen])
			pdr.FARID = &v
		case PDR_QER_ID:
			v := native.Uint32(b[n:attrLen])
			pdr.QERID = append(pdr.QERID, v)
		case PDR_URR_ID:
			v := native.Uint32(b[n:attrLen])
			pdr.URRID = append(pdr.URRID, v)
		case PDR_SEID:
			v := native.Uint64(b[n:attrLen])
			pdr.SEID = &v
		default:
			log.Printf("unknown type: %v\n", hdr.Type)
		}
		b = b[hdr.Len.Align():]
	}
	return pdr, nil
}

const (
	PDI_UE_ADDR_IPV4 = iota + 1
	PDI_F_TEID
	PDI_SDF_FILTER
	PDI_SRC_INTF
)

type PDI struct {
	UEAddr net.IP
	FTEID  *FTEID
	SDF    *SDFFilter
}

func DecodePDI(b []byte) (PDI, error) {
	var pdi PDI
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return pdi, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case PDI_UE_ADDR_IPV4:
			pdi.UEAddr = make([]byte, 4)
			copy(pdi.UEAddr, b[n:n+4])
		case PDI_F_TEID:
			fteid, err := DecodeFTEID(b[n:attrLen])
			if err != nil {
				return pdi, err
			}
			pdi.FTEID = &fteid
		case PDI_SDF_FILTER:
			sdf, err := DecodeSDFFilter(b[n:attrLen])
			if err != nil {
				return pdi, err
			}
			pdi.SDF = &sdf
		}
		b = b[hdr.Len.Align():]
	}
	return pdi, nil
}

const (
	F_TEID_I_TEID = iota + 1
	F_TEID_GTPU_ADDR_IPV4
)

type FTEID struct {
	TEID     uint32
	GTPuAddr net.IP
}

func DecodeFTEID(b []byte) (FTEID, error) {
	var fteid FTEID
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return fteid, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case F_TEID_I_TEID:
			fteid.TEID = native.Uint32(b[n:attrLen])
		case F_TEID_GTPU_ADDR_IPV4:
			fteid.GTPuAddr = make([]byte, 4)
			copy(fteid.GTPuAddr, b[n:n+4])
		}
		b = b[hdr.Len.Align():]
	}
	return fteid, nil
}

const (
	SDF_FILTER_FLOW_DESCRIPTION = iota + 1
	SDF_FILTER_TOS_TRAFFIC_CLASS
	SDF_FILTER_SECURITY_PARAMETER_INDEX
	SDF_FILTER_FLOW_LABEL
	SDF_FILTER_SDF_FILTER_ID
)

type SDFFilter struct {
	FD  *FlowDesc
	TTC *uint16
	SPI *uint32
	FL  *uint32
	BID *uint32
}

func DecodeSDFFilter(b []byte) (SDFFilter, error) {
	var sdf SDFFilter
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return sdf, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case SDF_FILTER_FLOW_DESCRIPTION:
			fd, err := DecodeFlowDesc(b[n:attrLen])
			if err != nil {
				return sdf, err
			}
			sdf.FD = &fd
		case SDF_FILTER_TOS_TRAFFIC_CLASS:
			v := native.Uint16(b[n:attrLen])
			sdf.TTC = &v
		case SDF_FILTER_SECURITY_PARAMETER_INDEX:
			v := native.Uint32(b[n:attrLen])
			sdf.SPI = &v
		case SDF_FILTER_FLOW_LABEL:
			v := native.Uint32(b[n:attrLen])
			sdf.FL = &v
		case SDF_FILTER_SDF_FILTER_ID:
			v := native.Uint32(b[n:attrLen])
			sdf.BID = &v
		default:
			log.Printf("unknown type: %v\n", hdr.Type)
		}
		b = b[hdr.Len.Align():]
	}
	return sdf, nil
}

const (
	FLOW_DESCRIPTION_ACTION = iota + 1
	FLOW_DESCRIPTION_DIRECTION
	FLOW_DESCRIPTION_PROTOCOL
	FLOW_DESCRIPTION_SRC_IPV4
	FLOW_DESCRIPTION_SRC_MASK
	FLOW_DESCRIPTION_DEST_IPV4
	FLOW_DESCRIPTION_DEST_MASK
	FLOW_DESCRIPTION_SRC_PORT
	FLOW_DESCRIPTION_DEST_PORT
)

const (
	SDF_FILTER_ACTION_UNSPEC = iota
	SDF_FILTER_PERMIT
)

const (
	SDF_FILTER_DIRECTION_UNSPEC = iota
	SDF_FILTER_IN
	SDF_FILTER_OUT
)

type FlowDesc struct {
	Action   uint8
	Dir      uint8
	Proto    uint8
	Src      net.IPNet
	Dst      net.IPNet
	SrcPorts [][]uint16
	DstPorts [][]uint16
}

func DecodeFlowDesc(b []byte) (FlowDesc, error) {
	var fd FlowDesc
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return fd, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case FLOW_DESCRIPTION_ACTION:
			fd.Action = b[n]
		case FLOW_DESCRIPTION_DIRECTION:
			fd.Dir = b[n]
		case FLOW_DESCRIPTION_PROTOCOL:
			fd.Proto = b[n]
		case FLOW_DESCRIPTION_SRC_IPV4:
			fd.Src.IP = make([]byte, 4)
			copy(fd.Src.IP, b[n:n+4])
		case FLOW_DESCRIPTION_SRC_MASK:
			fd.Src.Mask = make([]byte, 4)
			copy(fd.Src.Mask, b[n:n+4])
		case FLOW_DESCRIPTION_DEST_IPV4:
			fd.Dst.IP = make([]byte, 4)
			copy(fd.Dst.IP, b[n:n+4])
		case FLOW_DESCRIPTION_DEST_MASK:
			fd.Dst.Mask = make([]byte, 4)
			copy(fd.Dst.Mask, b[n:n+4])
		case FLOW_DESCRIPTION_SRC_PORT:
			for n < attrLen {
				v := native.Uint32(b[n:attrLen])
				lb := uint16(v >> 16)
				ub := uint16(v)
				x := []uint16{lb}
				if ub != lb {
					x = append(x, ub)
				}
				fd.SrcPorts = append(fd.SrcPorts, x)
				n += 4
			}
		case FLOW_DESCRIPTION_DEST_PORT:
			for n < attrLen {
				v := native.Uint32(b[n:attrLen])
				lb := uint16(v >> 16)
				ub := uint16(v)
				x := []uint16{lb}
				if ub != lb {
					x = append(x, ub)
				}
				fd.DstPorts = append(fd.DstPorts, x)
				n += 4
			}
		default:
			log.Printf("unknown type: %v\n", hdr.Type)
		}
		b = b[hdr.Len.Align():]
	}
	return fd, nil
}
