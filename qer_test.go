package gtp5gnl

import (
	"encoding/json"
	"log"
	"syscall"
	"testing"

	"github.com/khirono/go-nl"
)

func TestCreateQER(t *testing.T) {
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

	attrs := nl.AttrList{
		{
			Type:  QER_GATE,
			Value: nl.AttrU8(0),
		},
		{
			Type: QER_MBR,
			Value: nl.AttrList{
				{
					Type:  QER_MBR_UL_HIGH32,
					Value: nl.AttrU32(123),
				},
				{
					Type:  QER_MBR_UL_LOW8,
					Value: nl.AttrU8(4),
				},
				{
					Type:  QER_MBR_DL_HIGH32,
					Value: nl.AttrU32(765),
				},
				{
					Type:  QER_MBR_DL_LOW8,
					Value: nl.AttrU8(8),
				},
			},
		},
		{
			Type:  QER_QFI,
			Value: nl.AttrU8(9),
		},
	}

	err = CreateQER(c, link, 1, attrs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetQER(t *testing.T) {
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

	qer, err := GetQER(c, link, 1)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(qer, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("QER: %s\n", j)
}

func TestGetQERAll(t *testing.T) {
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

	qers, err := GetQERAll(c)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(qers, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("QERs: %s\n", j)
}

func TestRemoveQER(t *testing.T) {
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

	err = RemoveQER(c, link, 1)
	if err != nil {
		t.Fatal(err)
	}
}
