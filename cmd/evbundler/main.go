package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-loadtest/evbundler/cmd/evbundler/internal/base"
	"github.com/go-loadtest/evbundler/cmd/evbundler/internal/ping"
)

type Commands []*base.Command

func (cs Commands) Usage() {
	usage := `
Usage of evbundler:

		evbundler <command> [arguments]

The commands are:

`
	fmt.Fprint(os.Stderr, usage)
	for i := range cs {
		fmt.Fprintf(os.Stderr, "\t\t%s\t%s\n", cs[i].Name, cs[i].Doc)
	}
	os.Exit(2)
}

func main() {
	commands := Commands{
		ping.Cmd,
	}

	flag.Usage = commands.Usage
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		commands.Usage()
	}

	for _, c := range commands {
		if c.Name == args[0] {
			c.Flags.Usage = c.Usage
			c.Flags.Parse(args[1:])
			c.Run(c.Flags.Args())
		}
	}
}
