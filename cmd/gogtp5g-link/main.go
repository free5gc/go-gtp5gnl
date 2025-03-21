package main

import (
	"fmt"
	"os"
	"path"

	"github.com/free5gc/go-gtp5gnl/linkcmd"
)

func usage(prog string) {
	fmt.Fprintf(os.Stderr, "usage: %v <add|del> <ifname> [--ran]\n", prog)
}

func main() {
	prog := path.Base(os.Args[0])
	if len(os.Args) < 3 {
		usage(prog)
		os.Exit(1)
	}
	cmd := os.Args[1]
	ifname := os.Args[2]
	var role int
	if len(os.Args) > 3 && os.Args[3] == "--ran" {
		role = 1
	}

	switch cmd {
	case "add":
		err := linkcmd.CmdAdd(ifname, role)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", prog, err)
			os.Exit(1)
		}
	case "del":
		err := linkcmd.CmdDel(ifname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v: %v\n", prog, err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "%v: unknown command %q\n", prog, cmd)
		os.Exit(1)
	}
}
