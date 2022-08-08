package gtp5gnl

import (
	"syscall"

	"github.com/khirono/go-genl"
	"github.com/khirono/go-nl"
)

func GetVersion(c *Client) (string, error) {
	flags := syscall.NLM_F_ACK
	req := nl.NewRequest(c.ID, flags)
	err := req.Append(genl.Header{Cmd: CMD_GET_VERSION})
	if err != nil {
		return "", err
	}

	rsps, err := c.Do(req)
	if err != nil {
		return "", err
	}
	if len(rsps) < 1 {
		return "", err
	}
	ver, err := DecodeVersion(rsps[0].Body[genl.SizeofHeader:])
	if err != nil {
		return "", err
	}
	return ver, err
}
