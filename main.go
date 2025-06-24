package main

import (
	"fmt"
	"os"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/internal/config"
	"github.com/alexflint/go-arg"
)

type moduleCmd struct {
	Prime bool `arg:"-p,--prime" default:"false" help:"add and commit all changes to all modules"`
	List  bool `arg:"-l,--list" default:"false" help:"list all tracked modules"`
	Edit  bool `arg:"-e, --edit" default:"false" help:"cd into selected module via fzf"`
}


var args struct {
	Module   *moduleCmd   `arg:"subcommand:m" help:"operations for git modules."`
}

func main() {
	config.ReadUserConfig()
	stdinArgs := os.Args[1:]
	p, err := arg.NewParser(arg.Config{Program: "dotf", Exit: os.Exit}, &args)
	if err != nil {
		fmt.Println(err)
	}

	if len(os.Args) < 2 {
		// Print dotf help if no subcommands given
		p.WriteHelp(os.Stdout)
		os.Exit(0)
	}

	p.Parse(stdinArgs)

	switch {
	case args.Module != nil:
		// positional flags for `sm`
		switch {
		case args.Module.Prime:
			dotf.Prime()
		case args.Module.List:
			dotf.List()
		case args.Module.Edit:
			dotf.Edit()
		default:
			p.WriteHelp(os.Stdout)
		}
	default:
		// If arguments above aren't called, nor is `--help`
		if os.Args[1] != "--help" && os.Args[1] != "-h" {
			cmd.DotfExecute(stdinArgs)
		} else {
			var choice string
			fmt.Println("dotf wraps around git. \nDisplay help for dotf (d) or git (g)?")
			fmt.Scan(&choice)
			switch {
			case choice == "d":
				p.WriteHelp(os.Stdout)
			case choice == "g":
				cmd.DotfExecute(stdinArgs)
			}
		}
	}

}
