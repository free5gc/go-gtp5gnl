package gtp5gnl

import (
	"time"
	"unsafe"

	"github.com/khirono/go-nl"
)

// for UPF Usage Statistic
const (
	USTAT_UL_VOL_RX = iota + 1
	USTAT_UL_VOL_TX
	USTAT_DL_VOL_RX
	USTAT_DL_VOL_TX

	USTAT_UL_PKT_RX
	USTAT_UL_PKT_TX
	USTAT_DL_PKT_RX
	USTAT_DL_PKT_TX
)

type UsageStatistic struct {
	TotalVolRx uint64
	TotalVolTx uint64
	UlVolRx    uint64
	UlVolTx    uint64
	DlVolRx    uint64
	DlVolTx    uint64
	TotalPktRx uint64
	TotalPktTx uint64
	UlPktRx    uint64
	UlPktTx    uint64
	DlPktRx    uint64
	DlPktTx    uint64
}

const (
	UR = iota + 5
)

const (
	UR_URRID = iota + 3
	UR_USAGE_REPORT_TRIGGER
	UR_URSEQN
	UR_VOLUME_MEASUREMENT
	UR_QUERY_URR_REFERENCE
	UR_START_TIME
	UR_END_TIME
	UR_SEID
)

const (
	UR_VOLUME_MEASUREMENT_FLAGS = iota + 1

	UR_VOLUME_MEASUREMENT_TOVOL
	UR_VOLUME_MEASUREMENT_UVOL
	UR_VOLUME_MEASUREMENT_DVOL

	UR_VOLUME_MEASUREMENT_TOPACKET
	UR_VOLUME_MEASUREMENT_UPACKET
	UR_VOLUME_MEASUREMENT_DPACKET
)

const (
	TOVOL uint8 = 1 << iota
	ULVOL
	DLVOL
	TONOP
	ULNOP
	DLNOP
)

type USAReport struct {
	URRID          uint32
	URSEQN         uint32
	USARTrigger    uint32
	VolMeasurement VolumeMeasurement
	QueryUrrRef    uint32
	StartTime      time.Time
	EndTime        time.Time
	SEID           uint64
}

type VolumeMeasurement struct {
	Flag           uint8
	TotalVolume    uint64
	UplinkVolume   uint64
	DownlinkVolume uint64
	TotalPktNum    uint64
	UplinkPktNum   uint64
	DownlinkPktNum uint64
}

// The maximun netlink message size is 16K, and the body for the attibutes are 7856 Bytes
const (
	MAX_NETLINK_MSG_BODY_SIZE = 7856
	NETLIMK_ATTR_HDR_SIZE     = 4
)

// The netlink attribute size of UR need to count the UR header(4) + size of the attributes (and it's header) in UR
func MaxNetlinkUsageReportNum() int {
	size := NETLIMK_ATTR_HDR_SIZE // UR attr header

	size += NETLIMK_ATTR_HDR_SIZE         // UR_URRID attr header
	size += int(unsafe.Sizeof(uint32(0))) // UR_URRID attr data

	size += NETLIMK_ATTR_HDR_SIZE         // UR_USAGE_REPORT_TRIGGER attr header
	size += int(unsafe.Sizeof(uint32(0))) // UR_USAGE_REPORT_TRIGGER attr data

	size += NETLIMK_ATTR_HDR_SIZE         // UR_URSEQN attr header
	size += int(unsafe.Sizeof(uint32(0))) // UR_URSEQN attr data

	size += NETLIMK_ATTR_HDR_SIZE         // UR_VOLUME_MEASUREMENT attr header
	size += NETLIMK_ATTR_HDR_SIZE         // UR_VOLUME_MEASUREMENT_TOVOL attr header
	size += int(unsafe.Sizeof(uint64(0))) // UR_VOLUME_MEASUREMENT_TOVOL attr data
	size += NETLIMK_ATTR_HDR_SIZE         // UR_VOLUME_MEASUREMENT_UVOL attr header
	size += int(unsafe.Sizeof(uint64(0))) // UR_VOLUME_MEASUREMENT_UVOL attr data
	size += NETLIMK_ATTR_HDR_SIZE         // UR_VOLUME_MEASUREMENT_DVOL attr header
	size += int(unsafe.Sizeof(uint64(0))) // UR_VOLUME_MEASUREMENT_DVOL attr data
	size += NETLIMK_ATTR_HDR_SIZE         // UR_VOLUME_MEASUREMENT_TOPACKET attr header
	size += int(unsafe.Sizeof(uint64(0))) // UR_VOLUME_MEASUREMENT_TOPACKET attr data
	size += NETLIMK_ATTR_HDR_SIZE         // UR_VOLUME_MEASUREMENT_UPACKET attr header
	size += int(unsafe.Sizeof(uint64(0))) // UR_VOLUME_MEASUREMENT_UPACKET attr data
	size += NETLIMK_ATTR_HDR_SIZE         // UR_VOLUME_MEASUREMENT_DPACKET attr header
	size += int(unsafe.Sizeof(uint64(0))) // UR_VOLUME_MEASUREMENT_DPACKET attr data

	size += NETLIMK_ATTR_HDR_SIZE         // UR_START_TIME attr header
	size += int(unsafe.Sizeof(uint64(0))) // UR_START_TIME attr data

	size += NETLIMK_ATTR_HDR_SIZE         // UR_END_TIME attr header
	size += int(unsafe.Sizeof(uint64(0))) // UR_END_TIME attr data

	size += NETLIMK_ATTR_HDR_SIZE         // UR_SEID attr header
	size += int(unsafe.Sizeof(uint64(0))) // UR_SEID attr data

	return MAX_NETLINK_MSG_BODY_SIZE / size
}

func decodeVolumeMeasurement(b []byte) (VolumeMeasurement, error) {
	var VolMeasurement VolumeMeasurement
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return VolMeasurement, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case UR_VOLUME_MEASUREMENT_TOVOL:
			v := native.Uint64(b[n:attrLen])
			VolMeasurement.TotalVolume = v
			VolMeasurement.Flag |= TOVOL
		case UR_VOLUME_MEASUREMENT_UVOL:
			v := native.Uint64(b[n:attrLen])
			VolMeasurement.UplinkVolume = v
			VolMeasurement.Flag |= ULVOL
		case UR_VOLUME_MEASUREMENT_DVOL:
			v := native.Uint64(b[n:attrLen])
			VolMeasurement.DownlinkVolume = v
			VolMeasurement.Flag |= DLVOL
		case UR_VOLUME_MEASUREMENT_TOPACKET:
			v := native.Uint64(b[n:attrLen])
			VolMeasurement.TotalPktNum = v
			VolMeasurement.Flag |= TONOP
		case UR_VOLUME_MEASUREMENT_UPACKET:
			v := native.Uint64(b[n:attrLen])
			VolMeasurement.UplinkPktNum = v
			VolMeasurement.Flag |= ULNOP
		case UR_VOLUME_MEASUREMENT_DPACKET:
			v := native.Uint64(b[n:attrLen])
			VolMeasurement.DownlinkPktNum = v
			VolMeasurement.Flag |= DLNOP
		default:
			return VolMeasurement, nil
		}

		b = b[hdr.Len.Align():]
	}
	return VolMeasurement, nil
}

func DecodeUsageStatistic(b []byte) (*UsageStatistic, error) {
	ustat := new(UsageStatistic)

	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case USTAT_UL_VOL_RX:
			ustat.UlVolRx = native.Uint64(b[n:attrLen])
		case USTAT_UL_VOL_TX:
			ustat.UlVolTx = native.Uint64(b[n:attrLen])
		case USTAT_DL_VOL_RX:
			ustat.DlVolRx = native.Uint64(b[n:attrLen])
		case USTAT_DL_VOL_TX:
			ustat.DlVolTx = native.Uint64(b[n:attrLen])
		case USTAT_UL_PKT_RX:
			ustat.UlPktRx = native.Uint64(b[n:attrLen])
		case USTAT_UL_PKT_TX:
			ustat.UlPktTx = native.Uint64(b[n:attrLen])
		case USTAT_DL_PKT_RX:
			ustat.DlPktRx = native.Uint64(b[n:attrLen])
		case USTAT_DL_PKT_TX:
			ustat.DlPktTx = native.Uint64(b[n:attrLen])
		}

		b = b[hdr.Len.Align():]
	}

	// total volume count
	ustat.TotalVolRx = ustat.UlVolRx + ustat.DlVolRx
	ustat.TotalVolTx = ustat.UlVolTx + ustat.DlVolTx

	// total packet count
	ustat.TotalPktRx = ustat.UlPktRx + ustat.DlPktRx
	ustat.TotalPktTx = ustat.UlPktTx + ustat.DlPktTx

	return ustat, nil
}

func DecodeAllUSAReports(b []byte) ([]USAReport, error) {
	var usars []USAReport

	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case UR:
			r, err := decodeUSAReport(b[n:attrLen])
			if err != nil {
				return nil, err
			}
			usars = append(usars, *r)
		}

		b = b[hdr.Len.Align():]
	}
	return usars, nil
}

func decodeUSAReport(b []byte) (*USAReport, error) {
	report := new(USAReport)

	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		attrLen := int(hdr.Len)
		switch hdr.MaskedType() {
		case UR_URRID:
			report.URRID = native.Uint32(b[n:attrLen])
		case UR_USAGE_REPORT_TRIGGER:
			report.USARTrigger = native.Uint32(b[n:attrLen])
		case UR_URSEQN:
			report.URSEQN = native.Uint32(b[n:attrLen])
		case UR_VOLUME_MEASUREMENT:
			volMeasurement, err := decodeVolumeMeasurement(b[n:attrLen])
			if err != nil {
				return nil, err
			}
			report.VolMeasurement = volMeasurement
		case UR_START_TIME:
			v := native.Uint64(b[n:attrLen])
			report.StartTime = time.Unix(0, int64(v))
		case UR_END_TIME:
			v := native.Uint64(b[n:attrLen])
			report.EndTime = time.Unix(0, int64(v))
		case UR_SEID:
			report.SEID = native.Uint64(b[n:attrLen])
		}

		b = b[hdr.Len.Align():]
	}
	return report, nil
}
