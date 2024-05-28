package main

import (
	"fmt"
	"os"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/internal/config"
	"github.com/alexflint/go-arg"
)

type smCmd struct {
	Prime bool `arg:"-p,--prime" default:"false" help:"add and commit all changes to all submodules"`
	List  bool `arg:"-l,--list" default:"false" help:"list all tracked submodules"`
}

type pushCmd struct {
	Remote string `arg:"positional"`
}

var args struct {
	Push *pushCmd `arg:"subcommand:push" help:"push to origin and mirror concurrently"`
	Sm   *smCmd   `arg:"subcommand:sm" help:"operations for git submodules"`
}

func main() {
	config.ReadUserConfig()
	stdinArgs := os.Args[1:]
	p, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		fmt.Println(err)
	}
	err = p.Parse(stdinArgs)

	if len(os.Args) < 2 {
		// Print dotf help if no subcommands given
		p.WriteHelp(os.Stdout)
		return
	}

	switch {
	case args.Push != nil:
		dotf.Push(args.Push.Remote)
	case args.Sm != nil:
		// positional flags for `sm`
		switch {
		case args.Sm.Prime:
			dotf.Prime()
		case args.Sm.List:
			dotf.List()
		default:
			p.WriteHelp(os.Stdout)
		}
	default:
		// If arguments above aren't called, nor is `--help`
		if arg.ErrHelp == nil {
			cmd.DotfExecute(stdinArgs)
		} else {
			var choice string
			fmt.Println("Display help for dotf (d) or git (g)?")
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
