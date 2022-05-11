package gtp5gnl

import (
	"encoding/json"
	"log"
	"syscall"
	"testing"

	"github.com/khirono/go-nl"
)

func TestCreateBAR(t *testing.T) {
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
			Type:  BAR_DOWNLINK_DATA_NOTIFICATION_DELAY,
			Value: nl.AttrU8(2),
		},
		{
			Type:  BAR_BUFFERING_PACKETS_COUNT,
			Value: nl.AttrU16(3),
		},
	}

	err = CreateBAR(c, link, 1, attrs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetBAR(t *testing.T) {
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

	far, err := GetBAR(c, link, 1)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(far, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("BAR: %s\n", j)
}

func TestGetBARAll(t *testing.T) {
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

	fars, err := GetBARAll(c)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(fars, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("BARs: %s\n", j)
}

func TestRemoveBAR(t *testing.T) {
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

	err = RemoveBAR(c, link, 1)
	if err != nil {
		t.Fatal(err)
	}
}
