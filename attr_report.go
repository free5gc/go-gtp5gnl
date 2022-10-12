package gtp5gnl

import (
	"time"

	"github.com/khirono/go-nl"
)

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

func DecodeVolumeMeasurement(b []byte) (VolumeMeasurement, error) {
	var VolMeasurement VolumeMeasurement

	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return VolMeasurement, err
		}

		switch hdr.MaskedType() {
		case UR_VOLUME_MEASUREMENT_TOVOL:
			v := native.Uint64(b[n:])
			VolMeasurement.TotalVolume = v
			VolMeasurement.Flag |= TOVOL
		case UR_VOLUME_MEASUREMENT_UVOL:
			v := native.Uint64(b[n:])
			VolMeasurement.UplinkVolume = v
			VolMeasurement.Flag |= ULVOL
		case UR_VOLUME_MEASUREMENT_DVOL:
			v := native.Uint64(b[n:])
			VolMeasurement.DownlinkVolume = v
			VolMeasurement.Flag |= DLVOL
		case UR_VOLUME_MEASUREMENT_TOPACKET:
			v := native.Uint64(b[n:])
			VolMeasurement.TotalPktNum = v
			VolMeasurement.Flag |= TONOP
		case UR_VOLUME_MEASUREMENT_UPACKET:
			v := native.Uint64(b[n:])
			VolMeasurement.UplinkPktNum = v
			VolMeasurement.Flag |= ULNOP
		case UR_VOLUME_MEASUREMENT_DPACKET:
			v := native.Uint64(b[n:])
			VolMeasurement.DownlinkPktNum = v
			VolMeasurement.Flag |= DLNOP
		}

		b = b[hdr.Len.Align():]
	}
	return VolMeasurement, nil
}

func DecodeAllUSAReports(b []byte) ([]USAReport, error) {
	var usars []USAReport

	for len(b) > 0 {
		hdr, n, err := nl.DecodeAttrHdr(b)
		if err != nil {
			return nil, err
		}

		switch hdr.MaskedType() {
		case UR:
			r, err := DecodeUSAReport(b[n:])
			if err != nil {
				return nil, err
			}
			usars = append(usars, *r)
		}

		b = b[hdr.Len.Align():]
	}
	return usars, nil
}

func DecodeUSAReport(b []byte) (*USAReport, error) {
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
				return nil, err
			}
			report.VolMeasurement = volMeasurement
		case UR_START_TIME:
			v := native.Uint64(b[n:])
			report.StartTime = time.Unix(0, int64(v))
		case UR_END_TIME:
			v := native.Uint64(b[n:])
			report.EndTime = time.Unix(0, int64(v))
		}

		b = b[hdr.Len.Align():]
	}
	return report, nil
}
