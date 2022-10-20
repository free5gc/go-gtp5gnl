package gtp5gnl

import (
	"net"
	"sync"
	"syscall"
	"testing"

	"github.com/khirono/go-nl"
	"github.com/khirono/go-rtnllink"
)

func TestCreateLink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	var wg sync.WaitGroup
	mux, err := nl.NewMux()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		mux.Close()
		wg.Wait()
	}()
	wg.Add(1)
	go func() {
		mux.Serve()
		wg.Done()
	}()

	conn, err := nl.Open(syscall.NETLINK_ROUTE)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := nl.NewClient(conn, mux)

	laddr, err := net.ResolveUDPAddr("udp", ":2152")
	if err != nil {
		t.Fatal(err)
	}
	conn2, err := net.ListenUDP("udp", laddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn2.Close()
	f, err := conn2.File()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	linkinfo := &nl.Attr{
		Type: syscall.IFLA_LINKINFO,
		Value: nl.AttrList{
			{
				Type:  rtnllink.IFLA_INFO_KIND,
				Value: nl.AttrString("gtp5g"),
			},
			{
				Type: rtnllink.IFLA_INFO_DATA,
				Value: nl.AttrList{
					{
						Type:  IFLA_FD1,
						Value: nl.AttrU32(f.Fd()),
					},
					{
						Type:  IFLA_HASHSIZE,
						Value: nl.AttrU32(1024),
					},
				},
			},
		},
	}
	err = rtnllink.Create(c, "upfgtp", linkinfo)
	if err != nil {
		t.Fatal(err)
	}

	err = rtnllink.Up(c, "upfgtp")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveLink(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}

	var wg sync.WaitGroup
	mux, err := nl.NewMux()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		mux.Close()
		wg.Wait()
	}()
	wg.Add(1)
	go func() {
		mux.Serve()
		wg.Done()
	}()

	conn, err := nl.Open(syscall.NETLINK_ROUTE)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c := nl.NewClient(conn, mux)

	err = rtnllink.Remove(c, "upfgtp")
	if err != nil {
		t.Fatal(err)
	}
}
