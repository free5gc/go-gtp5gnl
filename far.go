package gtp5gnl

import (
	"syscall"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"
)

// m LINK: ifindex
// m FAR_ID: u32
// m FAR_APPLY_ACTION: u8
// m FAR_FORWARDING_PARAMETER {
// m   FORWARDING_PARAMETER_OUTER_HEADER_CREATION {
//       OUTER_HEADER_CREATION_DESCRIPTION: u16
//       OUTER_HEADER_CREATION_O_TEID: u32
//       OUTER_HEADER_CREATION_PEER_ADDR_IPV4: [4]byte IP
//       OUTER_HEADER_CREATION_PORT: u16
//     }
// o   FORWARDING_PARAMETER_FORWARDING_POLICY: string
//   }
func CreateFAR(c *Client, link *Link, farid int, attrs []nl.Attr) error {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_FAR})
	if err != nil {
		return err
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  FAR_ID,
			Value: nl.AttrU32(farid),
		},
	})
	if err != nil {
		return err
	}
	err = req.Append(nl.AttrList(attrs))
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func UpdateFAR(c *Client, link *Link, farid int, attrs []nl.Attr) error {
	flags := syscall.NLM_F_REPLACE
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_ADD_FAR})
	if err != nil {
		return err
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  FAR_ID,
			Value: nl.AttrU32(farid),
		},
	})
	if err != nil {
		return err
	}
	err = req.Append(nl.AttrList(attrs))
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func RemoveFAR(c *Client, link *Link, farid int) error {
	flags := syscall.NLM_F_EXCL
	flags |= syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_DEL_FAR})
	if err != nil {
		return err
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  FAR_ID,
			Value: nl.AttrU32(farid),
		},
	})
	if err != nil {
		return err
	}
	_, err = c.Do(req)
	return err
}

func GetFAR(c *Client, link *Link, farid int) (*FAR, error) {
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_FAR})
	if err != nil {
		return nil, err
	}
	err = req.Append(nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  FAR_ID,
			Value: nl.AttrU32(farid),
		},
	})
	if err != nil {
		return nil, err
	}
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if len(rsps) < 1 {
		return nil, err
	}
	far, err := DecodeFAR(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return nil, err
	}
	return far, err
}

func GetFARAll(c *Client) ([]FAR, error) {
	flags := syscall.NLM_F_DUMP
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_FAR})
	if err != nil {
		return nil, err
	}
	rsps, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	var fars []FAR
	for _, rsp := range rsps {
		far, err := DecodeFAR(rsp.Body[genl.SizeofHeader:])
		if err != nil {
			return nil, err
		}
		fars = append(fars, *far)
	}
	return fars, err
}
