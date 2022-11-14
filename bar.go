package gtp5gnl

import (
	"fmt"
	"syscall"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"
)

// m LINK: ifindex
// o BAR_SEID: u64
// m BAR_ID: u8
// o BAR_DOWNLINK_DATA_NOTIFICATION_DELAY: u8
// o BAR_BUFFERING_PACKETS_COUNT: u16
func CreateBAR(c *Client, link *Link, barid int, attrs []nl.Attr) error {
	return CreateBAROID(c, link, OID{uint64(barid)}, attrs)
}

func CreateBAROID(c *Client, link *Link, oid OID, attrs []nl.Attr) error {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_BAR})
	if err != nil {
		return err
	}
	barid, ok := oid.ID()
	if !ok {
		return fmt.Errorf("invalid oid: %v", oid)
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  BAR_ID,
			Value: nl.AttrU8(barid),
		},
	})
	if err != nil {
		return err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  BAR_SEID,
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

func UpdateBAR(c *Client, link *Link, barid int, attrs []nl.Attr) error {
	return UpdateBAROID(c, link, OID{uint64(barid)}, attrs)
}

func UpdateBAROID(c *Client, link *Link, oid OID, attrs []nl.Attr) error {
	flags := syscall.NLM_F_REPLACE
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_BAR})
	if err != nil {
		return err
	}
	barid, ok := oid.ID()
	if !ok {
		return fmt.Errorf("invalid oid: %v", oid)
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  BAR_ID,
			Value: nl.AttrU8(barid),
		},
	})
	if err != nil {
		return err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  BAR_SEID,
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

func RemoveBAR(c *Client, link *Link, barid int) error {
	return RemoveBAROID(c, link, OID{uint64(barid)})
}

func RemoveBAROID(c *Client, link *Link, oid OID) error {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_DEL_BAR})
	if err != nil {
		return err
	}
	barid, ok := oid.ID()
	if !ok {
		return fmt.Errorf("invalid oid: %v", oid)
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  BAR_ID,
			Value: nl.AttrU8(barid),
		},
	})
	if err != nil {
		return err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  BAR_SEID,
			Value: nl.AttrU64(seid),
		})
		if err != nil {
			return err
		}
	}
	_, err = c.Do(req)
	return err
}

func GetBAR(c *Client, link *Link, barid int) (*BAR, error) {
	return GetBAROID(c, link, OID{uint64(barid)})
}

func GetBAROID(c *Client, link *Link, oid OID) (*BAR, error) {
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_BAR})
	if err != nil {
		return nil, err
	}
	barid, ok := oid.ID()
	if !ok {
		return nil, fmt.Errorf("invalid oid: %v", oid)
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  BAR_ID,
			Value: nl.AttrU8(barid),
		},
	})
	if err != nil {
		return nil, err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  BAR_SEID,
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
		return nil, fmt.Errorf("nil BAR of oid(%v)", oid)
	}
	bar, err := DecodeBAR(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return nil, err
	}
	return bar, err
}

func GetBARAll(c *Client) ([]BAR, error) {
	flags := syscall.NLM_F_DUMP
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_BAR})
	if err != nil {
		return nil, err
	}
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var bars []BAR
	for _, rsp := range rsps {
		bar, err := DecodeBAR(rsp.Body[genl.SizeofHeader:])
		if err != nil {
			return nil, err
		}
		bars = append(bars, *bar)
	}
	return bars, err
}
