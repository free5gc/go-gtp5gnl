package gtp5gnl

import (
	"github.com/khirono/go-nl"
)

const (
	IFLA_UNSPEC = iota
	IFLA_FD1
	IFLA_HASHSIZE
	IFLA_ROLE
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
