package gtp5gnl

import (
	"encoding/json"
	"log"
	"syscall"
	"testing"

	"github.com/khirono/go-nl"
)

func TestCreateFAR(t *testing.T) {
	mux, err := nl.NewMux()
	if err != nil {
		t.Fatal(err)
	}
	defer mux.Close()
	go mux.Serve()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c, err := NewClient(conn, mux)
	if err != nil {
		t.Fatal(err)
	}

	link, err := GetLink("upfgtp")
	if err != nil {
		t.Fatal(err)
	}

	attrs := []nl.Attr{
		{
			Type:  FAR_APPLY_ACTION,
			Value: nl.AttrU8(1),
		},
		{
			Type: FAR_FORWARDING_PARAMETER,
			Value: nl.AttrList{
				{
					Type: FORWARDING_PARAMETER_OUTER_HEADER_CREATION,
					Value: nl.AttrList{
						{
							Type:  OUTER_HEADER_CREATION_DESCRIPTION,
							Value: nl.AttrU16(255),
						},
						{
							Type:  OUTER_HEADER_CREATION_O_TEID,
							Value: nl.AttrU32(10),
						},
						{
							Type:  OUTER_HEADER_CREATION_PEER_ADDR_IPV4,
							Value: nl.AttrBytes([]byte{20, 20, 20, 2}),
						},
						{
							Type:  OUTER_HEADER_CREATION_PORT,
							Value: nl.AttrU16(3),
						},
					},
				},
				{
					Type:  FORWARDING_PARAMETER_PFCPSM_REQ_FLAGS,
					Value: nl.AttrU8(2),
				},
			},
		},
	}

	err = CreateFAR(c, link, 1, attrs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetFAR(t *testing.T) {
	mux, err := nl.NewMux()
	if err != nil {
		t.Fatal(err)
	}
	defer mux.Close()
	go mux.Serve()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c, err := NewClient(conn, mux)
	if err != nil {
		t.Fatal(err)
	}

	link, err := GetLink("upfgtp")
	if err != nil {
		t.Fatal(err)
	}

	far, err := GetFAR(c, link, 1)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(far, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("FAR: %s\n", j)
}

func TestGetFARAll(t *testing.T) {
	mux, err := nl.NewMux()
	if err != nil {
		t.Fatal(err)
	}
	defer mux.Close()
	go mux.Serve()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c, err := NewClient(conn, mux)
	if err != nil {
		t.Fatal(err)
	}

	fars, err := GetFARAll(c)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(fars, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("FARs: %s\n", j)
}

func TestRemoveFAR(t *testing.T) {
	mux, err := nl.NewMux()
	if err != nil {
		t.Fatal(err)
	}
	defer mux.Close()
	go mux.Serve()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c, err := NewClient(conn, mux)
	if err != nil {
		t.Fatal(err)
	}

	link, err := GetLink("upfgtp")
	if err != nil {
		t.Fatal(err)
	}

	err = RemoveFAR(c, link, 1)
	if err != nil {
		t.Fatal(err)
	}
}
