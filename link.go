package gtp5gnl

import (
	"syscall"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"
)

const (
	IFLA_UNSPEC = iota
	IFLA_FD1
	IFLA_HASHSIZE
	IFLA_ROLE
	IFLA_ETH_DEVS
)

type Link struct {
	Name  string
	Index int
}

func GetLink(name string) (*Link, error) {
	l := new(Link)
	index, err := nl.IfnameToIndex(name)
	if err != nil {
		return nil, err
	}
	l.Name = name
	l.Index = index
	return l, nil
}

func SetEthDevs(c *Client, link *Link, devs string) error {
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	if err := req.Append(genl.Header{Cmd: CMD_SET_ETH_DEVS}); err != nil {
		return err
	}
	if err := req.Append(&nl.AttrList{
		{
			Type:  LINK,
			Value: nl.AttrU32(link.Index),
		},
		{
			Type:  IFLA_ETH_DEVS,
			Value: nl.AttrString(devs),
		},
	}); err != nil {
		return err
	}
	_, err := c.Do(req)
	return err
}
