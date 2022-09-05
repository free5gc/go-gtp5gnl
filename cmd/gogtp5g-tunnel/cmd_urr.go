package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"syscall"

	"github.com/free5gc/go-gtp5gnl"
	"github.com/khirono/go-nl"
)

func ParseURROptions(args []string) ([]nl.Attr, error) {
	var attrs []nl.Attr
	p := NewCmdParser(args)
	for {
		opt, ok := p.GetToken()
		if !ok {
			break
		}
		switch opt {
		//TODO:
		}
	}

	return attrs, nil
}

// add urr <ifname> <oid> [options...]
func CmdAddURR(args []string) error {
	if len(args) < 2 {
		return errors.New("too few parameter")
	}
	ifname := args[0]
	oid, err := ParseOID(args[1])
	if err != nil {
		return err
	}
	attrs, err := ParseURROptions(args[2:])
	if err != nil {
		return err
	}

	mux, err := nl.NewMux()
	if err != nil {
		return err
	}
	defer mux.Close()
	go mux.Serve()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := gtp5gnl.NewClient(conn, mux)
	if err != nil {
		return err
	}

	link, err := gtp5gnl.GetLink(ifname)
	if err != nil {
		return err
	}

	return gtp5gnl.CreateURROID(c, link, oid, attrs)
}

// mod urr <ifname> <oid> [options...]
func CmdModURR(args []string) ([]gtp5gnl.USAReport, error) {
	if len(args) < 2 {
		return nil, errors.New("too few parameter")
	}
	ifname := args[0]
	oid, err := ParseOID(args[1])
	if err != nil {
		return nil, err
	}
	attrs, err := ParseURROptions(args[2:])
	if err != nil {
		return nil, err
	}

	mux, err := nl.NewMux()
	if err != nil {
		return nil, err
	}
	defer mux.Close()
	go mux.Serve()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c, err := gtp5gnl.NewClient(conn, mux)
	if err != nil {
		return nil, err
	}

	link, err := gtp5gnl.GetLink(ifname)
	if err != nil {
		return nil, err
	}

	return gtp5gnl.UpdateURROID(c, link, oid, attrs)
}

// delete urr <ifname> <oid>
func CmdDeleteURR(args []string) ([]gtp5gnl.USAReport, error) {
	if len(args) < 2 {
		return nil, errors.New("too few parameter")
	}
	ifname := args[0]
	oid, err := ParseOID(args[1])
	if err != nil {
		return nil, err
	}

	mux, err := nl.NewMux()
	if err != nil {
		return nil, err
	}
	defer mux.Close()
	go mux.Serve()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c, err := gtp5gnl.NewClient(conn, mux)
	if err != nil {
		return nil, err
	}

	link, err := gtp5gnl.GetLink(ifname)
	if err != nil {
		return nil, err
	}

	return gtp5gnl.RemoveURROID(c, link, oid)
}

// get urr <ifname> <oid>
func CmdGetURR(args []string) error {
	if len(args) < 2 {
		return errors.New("too few parameter")
	}
	ifname := args[0]
	oid, err := ParseOID(args[1])
	if err != nil {
		return err
	}

	mux, err := nl.NewMux()
	if err != nil {
		return err
	}
	defer mux.Close()
	go mux.Serve()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := gtp5gnl.NewClient(conn, mux)
	if err != nil {
		return err
	}

	link, err := gtp5gnl.GetLink(ifname)
	if err != nil {
		return err
	}

	urr, err := gtp5gnl.GetURROID(c, link, oid)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(urr, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", j)
	return nil
}

// list urr
func CmdListURR(args []string) error {
	mux, err := nl.NewMux()
	if err != nil {
		return err
	}
	defer mux.Close()
	go mux.Serve()

	conn, err := nl.Open(syscall.NETLINK_GENERIC)
	if err != nil {
		return err
	}
	defer conn.Close()

	c, err := gtp5gnl.NewClient(conn, mux)
	if err != nil {
		return err
	}

	urrs, err := gtp5gnl.GetURRAll(c)
	if err != nil {
		return err
	}

	j, err := json.MarshalIndent(urrs, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", j)
	return nil
}
