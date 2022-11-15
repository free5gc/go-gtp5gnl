package gtp5gnl

import (
	"fmt"
	"syscall"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"
)

// m LINK: ifindex
// o URR_SEID: u64
// m URR_ID: u32
// o URR_MEASUREMENT_METHOD: u64
// o URR_REPORTING_TRIGGER: u64
// o URR_MEASUREMENT_PERIOD: u64
// o URR_MEASUREMENT_INFO: u64
// o URR_SEQ: u64
func CreateURR(c *Client, link *Link, urrid int, attrs []nl.Attr) error {
	return CreateURROID(c, link, OID{uint64(urrid)}, attrs)
}

func CreateURROID(c *Client, link *Link, oid OID, attrs []nl.Attr) error {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_URR})
	if err != nil {
		return err
	}
	urrid, ok := oid.ID()
	if !ok {
		return fmt.Errorf("invalid oid: %v", oid)
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
		return err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  URR_SEID,
			Value: nl.AttrU64(seid),
		})
		if err != nil {
			return err
		}
	}
	err = req.Append(nl.AttrList(attrs))
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func UpdateURR(c *Client, link *Link, urrid int, attrs []nl.Attr) ([]USAReport, error) {
	return UpdateURROID(c, link, OID{uint64(urrid)}, attrs)
}

func UpdateURROID(c *Client, link *Link, oid OID, attrs []nl.Attr) ([]USAReport, error) {
	flags := syscall.NLM_F_REPLACE
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_URR})
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
	err = req.Append(nl.AttrList(attrs))
	if err != nil {
		return nil, err
	}
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if len(rsps) < 1 {
		return nil, nil
	}
	reports, err := DecodeAllUSAReports(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return nil, err
	}
	return reports, err
}

func RemoveURR(c *Client, link *Link, urrid int) ([]USAReport, error) {
	return RemoveURROID(c, link, OID{uint64(urrid)})
}

func RemoveURROID(c *Client, link *Link, oid OID) ([]USAReport, error) {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_DEL_URR})
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
		return nil, fmt.Errorf("RemoveURROID(%v): no usage report", oid)
	}
	reports, err := DecodeAllUSAReports(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return nil, err
	}
	return reports, err
}

func GetURR(c *Client, link *Link, urrid int) (*URR, error) {
	return GetURROID(c, link, OID{uint64(urrid)})
}

func GetURROID(c *Client, link *Link, oid OID) (*URR, error) {
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_URR})
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
		return nil, fmt.Errorf("GetURROID: nil URR of oid(%v)", oid)
	}
	urr, err := DecodeURR(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return nil, err
	}
	return urr, err
}

func GetURRAll(c *Client) ([]URR, error) {
	flags := syscall.NLM_F_DUMP
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_URR})
	if err != nil {
		return nil, err
	}
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var urrs []URR
	for _, rsp := range rsps {
		urr, err := DecodeURR(rsp.Body[genl.SizeofHeader:])
		if err != nil {
			return nil, err
		}
		urrs = append(urrs, *urr)
	}
	return urrs, err
}
