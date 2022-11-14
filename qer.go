package gtp5gnl

import (
	"fmt"
	"syscall"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"
)

// m QER_ID: u32
// o QER_GATE: u8
// o QER_MBR {
// m   QER_MBR_UL_HIGH32: u32
// m   QER_MBR_UL_LOW8: u8
// m   QER_MBR_DL_HIGH32: u32
// m   QER_MBR_DL_LOW8: u8
//   }
// o QER_GBR {
// m   QER_GBR_UL_HIGH32: u32
// m   QER_GBR_UL_LOW8: u8
// m   QER_GBR_DL_HIGH32: u32
// m   QER_GBR_DL_LOW8: u8
//   }
// o QER_CORR_ID: u32
// o QER_RQI: u8
// o QER_QFI: u8
// o QER_PPI: u8
//
func CreateQER(c *Client, link *Link, qerid int, attrs []nl.Attr) error {
	return CreateQEROID(c, link, OID{uint64(qerid)}, attrs)
}

func CreateQEROID(c *Client, link *Link, oid OID, attrs []nl.Attr) error {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_QER})
	if err != nil {
		return err
	}
	qerid, ok := oid.ID()
	if !ok {
		return fmt.Errorf("invalid oid: %v", oid)
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  QER_ID,
			Value: nl.AttrU32(qerid),
		},
	})
	if err != nil {
		return err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  QER_SEID,
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

func UpdateQER(c *Client, link *Link, qerid int, attrs []nl.Attr) error {
	return UpdateQEROID(c, link, OID{uint64(qerid)}, attrs)
}

func UpdateQEROID(c *Client, link *Link, oid OID, attrs []nl.Attr) error {
	flags := syscall.NLM_F_REPLACE
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_QER})
	if err != nil {
		return err
	}
	qerid, ok := oid.ID()
	if !ok {
		return fmt.Errorf("invalid oid: %v", oid)
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  QER_ID,
			Value: nl.AttrU32(qerid),
		},
	})
	if err != nil {
		return err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  QER_SEID,
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

func RemoveQER(c *Client, link *Link, qerid int) error {
	return RemoveQEROID(c, link, OID{uint64(qerid)})
}

func RemoveQEROID(c *Client, link *Link, oid OID) error {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_DEL_QER})
	if err != nil {
		return err
	}
	qerid, ok := oid.ID()
	if !ok {
		return fmt.Errorf("invalid oid: %v", oid)
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  QER_ID,
			Value: nl.AttrU32(qerid),
		},
	})
	if err != nil {
		return err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  QER_SEID,
			Value: nl.AttrU64(seid),
		})
		if err != nil {
			return err
		}
	}
	_, err = c.Do(req)
	return err
}

func GetQER(c *Client, link *Link, qerid int) (*QER, error) {
	return GetQEROID(c, link, OID{uint64(qerid)})
}

func GetQEROID(c *Client, link *Link, oid OID) (*QER, error) {
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_QER})
	if err != nil {
		return nil, err
	}
	qerid, ok := oid.ID()
	if !ok {
		return nil, fmt.Errorf("invalid oid: %v", oid)
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  QER_ID,
			Value: nl.AttrU32(qerid),
		},
	})
	if err != nil {
		return nil, err
	}
	seid, ok := oid.SEID()
	if ok {
		err = req.Append(&nl.Attr{
			Type:  QER_SEID,
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
		return nil, fmt.Errorf("nil QER of oid(%v)", oid)
	}
	qer, err := DecodeQER(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return nil, err
	}
	return qer, err
}

func GetQERAll(c *Client) ([]QER, error) {
	flags := syscall.NLM_F_DUMP
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_QER})
	if err != nil {
		return nil, err
	}
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var qers []QER
	for _, rsp := range rsps {
		qer, err := DecodeQER(rsp.Body[genl.SizeofHeader:])
		if err != nil {
			return nil, err
		}
		qers = append(qers, *qer)
	}
	return qers, err
}
