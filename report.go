package gtp5gnl

import (
	"fmt"
	"syscall"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"
)

func GetReport(c *Client, link *Link, urrids []uint64, seids []uint64) ([]USAReport, error) {
	var oids []OID
	for i, lseid := range seids {
		oids = append(oids, OID{lseid, uint64(urrids[i])})
	}
	return GetReportOID(c, link, oids)
}

func GetReportOID(c *Client, link *Link, oids []OID) ([]USAReport, error) {
	// var attrs nl.AttrList
	var attrs []nl.Attr
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_REPORT})
	if err != nil {
		return nil, err
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  URR_NUM,
			Value: nl.AttrU32(len(oids)),
		},
	})
	if err != nil {
		return nil, err
	}

	for _, oid := range oids {
		urrid, ok := oid.ID()
		if !ok {
			return nil, fmt.Errorf("invalid oid: %v", oid)
		}

		seid, ok := oid.SEID()
		if ok {
			attrs = append(attrs, nl.Attr{
				Type: SESS_URRS,
				Value: nl.AttrList{
					{
						Type:  URR_ID,
						Value: nl.AttrU32(urrid),
					},
					{
						Type:  URR_SEID,
						Value: nl.AttrU64(seid),
					},
				},
			},
			)
		}
	}
	err = req.Append(nl.AttrList(attrs))
	if err != nil {
		return nil, err
	}

	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if len(rsps) < 1 {
		return nil, fmt.Errorf("nil Report")
	}
	reports, err := DecodeAllUSAReports(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return nil, err
	}
	return reports, err
}
