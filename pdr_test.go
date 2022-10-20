package gtp5gnl

import (
	"encoding/json"
	"log"
	"sync"
	"syscall"
	"testing"

	"github.com/khirono/go-nl"
)

func TestCreatePDR(t *testing.T) {
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
			Type:  PDR_PRECEDENCE,
			Value: nl.AttrU32(254),
		},
		{
			Type:  PDR_OUTER_HEADER_REMOVAL,
			Value: nl.AttrU8(1),
		},
		{
			Type:  PDR_FAR_ID,
			Value: nl.AttrU32(1),
		},
		{
			Type:  PDR_QER_ID,
			Value: nl.AttrU32(1),
		},
		{
			Type: PDR_PDI,
			Value: nl.AttrList{
				{
					Type:  PDI_UE_ADDR_IPV4,
					Value: nl.AttrBytes([]byte{60, 60, 0, 1}),
				},
				{
					Type: PDI_F_TEID,
					Value: nl.AttrList{
						{
							Type:  F_TEID_I_TEID,
							Value: nl.AttrU32(5),
						},
						{
							Type:  F_TEID_GTPU_ADDR_IPV4,
							Value: nl.AttrBytes([]byte{40, 40, 40, 2}),
						},
					},
				},
			},
		},
	}

	err = CreatePDR(c, link, 2, attrs)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetPDR(t *testing.T) {
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

	pdr, err := GetPDR(c, link, 2)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(pdr, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("PDR: %s\n", j)
}

func TestGetPDRAll(t *testing.T) {
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

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	c, err := NewClient(conn, mux)
	if err != nil {
		t.Fatal(err)
	}

	pdrs, err := GetPDRAll(c)
	if err != nil {
		t.Fatal(err)
	}

	j, err := json.MarshalIndent(pdrs, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("PDRs: %s\n", j)
}

func TestRemovePDR(t *testing.T) {
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

	err = RemovePDR(c, link, 2)
	if err != nil {
		t.Fatal(err)
	}
}
