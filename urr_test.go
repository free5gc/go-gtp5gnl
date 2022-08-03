package gtp5gnl

import (
	"encoding/json"
	"log"
	"sync"
	"syscall"
	"testing"

	"github.com/khirono/go-nl"
)

func TestCreateURR(t *testing.T) {
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
			Type:  URR_MEASUREMENT_METHOD,
			Value: nl.AttrU64(2),
		},
		{
			Type:  URR_REPORTING_TRIGGER,
			Value: nl.AttrU64(3),
		},
		{
			Type:  URR_MEASUREMENT_PERIOD,
			Value: nl.AttrU64(4),
		},
		{
			Type:  URR_MEASUREMENT_INFO,
			Value: nl.AttrU64(5),
		},
		{
			Type:  URR_SEQ,
			Value: nl.AttrU64(6),
		},
	}

	err = CreateURR(c, link, 1, attrs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetURR(t *testing.T) {
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

	urr, err := GetURR(c, link, 1)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(urr, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("URR: %s\n", j)
}

func TestGetURRAll(t *testing.T) {
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

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c, err := NewClient(conn, mux)
	if err != nil {
		t.Fatal(err)
	}

	urrs, err := GetURRAll(c)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(urrs, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("URRs: %s\n", j)
}

func TestRemoveURR(t *testing.T) {
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

	err = RemoveURR(c, link, 1)
	if err != nil {
		t.Fatal(err)
	}
}
