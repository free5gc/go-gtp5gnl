package gtp5gnl

import (
	"fmt"
	"syscall"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"
)

func GetReport(c *Client, link *Link, urrid uint64, seid uint64) ([]USAReport, error) {
	return GetReportOID(c, link, OID{uint64(urrid), seid})
}

func GetReportOID(c *Client, link *Link, oid OID) ([]USAReport, error) {
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_REPORT})
	if err != nil {
		return nil, err
	}
	urrid, ok := oid.ID()
	if !ok {
		return nil, fmt.Errorf("invalid oid: %v", oid)
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  URR_ID,
			Value: nl.AttrU32(urrid),
		},
	})
	if err != nil {
		return nil, err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  URR_SEID,
			Value: nl.AttrU64(seid),
		})
		if err != nil {
			return nil, err
		}
	}
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if len(rsps) < 1 {
		return nil, err
	}
	report, err := DecodeReport(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return nil, err
	}
	return report, err
}
