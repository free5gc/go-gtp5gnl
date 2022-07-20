package gtp5gnl

import (
	"github.com/khirono/go-nl"
)

const (
	UR_URRID = iota + 3
	UR_USAGE_REPORT_TRIGGER
	UR_URSEQN
	UR_VOLUME_MEASUREMENT
	UR_QUERY_URR_REFERENCE
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

type USAReport struct {
	URRID          uint32
	URSEQN         uint32
	USARTrigger    uint32
	VolMeasurement VolumeMeasurement
	QueryUrrRef    uint32
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

type UsageReportTrigger struct {
	PERIO uint8
	VOLTH uint8
	TIMTH uint8
	QUHTI uint8
	START uint8
	STOPT uint8
	DROTH uint8
	IMMER uint8
	VOLQU uint8
	TIMQU uint8
	LIUSA uint8
	TERMR uint8
	MONIT uint8
	ENVCL uint8
	MACAR uint8
	EVETH uint8
	EVEQU uint8
	TEBUR uint8
	IPMJL uint8
	QUVTI uint8
	EMRRE uint8
}

func DecodeVolumeMeasurement(b []byte) (VolumeMeasurement, error) {
	var VolMeasurement VolumeMeasurement
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return VolMeasurement, err
		}
		switch hdr.MaskedType() {
		case UR_VOLUME_MEASUREMENT_FLAGS:
			v := uint8(b[n])
			VolMeasurement.Flag = v
		case UR_VOLUME_MEASUREMENT_TOVOL:
			v := native.Uint64(b[n:])
			VolMeasurement.TotalVolume = v
		case UR_VOLUME_MEASUREMENT_UVOL:
			v := native.Uint64(b[n:])
			VolMeasurement.UplinkVolume = v
		case UR_VOLUME_MEASUREMENT_DVOL:
			v := native.Uint64(b[n:])
			VolMeasurement.DownlinkVolume = v
		case UR_VOLUME_MEASUREMENT_TOPACKET:
			v := native.Uint64(b[n:])
			VolMeasurement.TotalPktNum = v
		case UR_VOLUME_MEASUREMENT_UPACKET:
			v := native.Uint64(b[n:])
			VolMeasurement.UplinkPktNum = v
		case UR_VOLUME_MEASUREMENT_DPACKET:
			v := native.Uint64(b[n:])
			VolMeasurement.DownlinkPktNum = v
		}

		b = b[hdr.Len.Align():]
	}
	return VolMeasurement, nil
}

func DecodeReport(b []byte) (*USAReport, error) {
	report := new(USAReport)
	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}
		switch hdr.MaskedType() {
		case UR_URRID:
			report.URRID = native.Uint32(b[n:])
		case UR_USAGE_REPORT_TRIGGER:
			report.USARTrigger = native.Uint32(b[n:])
		case UR_URSEQN:
			report.URSEQN = native.Uint32(b[n:])
		case UR_VOLUME_MEASUREMENT:
			volMeasurement, err := DecodeVolumeMeasurement(b[n:])
			if err != nil {
				return report, err
			}
			report.VolMeasurement = volMeasurement
		}
		b = b[hdr.Len.Align():]
	}
	return report, nil
}
