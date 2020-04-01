package base

import "flag"

type Command struct {
	Name  string
	Doc   string
	Run   func(args []string) error
	Flags *flag.FlagSet
}

// Commands is register subcommands
var Commands []*Command
