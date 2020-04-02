package base

import (
	"flag"
	"fmt"
	"os"
)

type Command struct {
	Name  string
	Doc   string
	Run   func(args []string) error
	Flags flag.FlagSet
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s\n\n", c.Name)
	c.Flags.PrintDefaults()
	os.Exit(2)
}
