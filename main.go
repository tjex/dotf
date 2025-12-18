package main

import (
	"fmt"
	"os"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/printer"
	"github.com/alexflint/go-arg"
)

var args struct {
	ModuleCmd *dotf.ModuleCmd `arg:"subcommand:m" help:"operations for git modules."`
	Quiet     bool            `arg:"-q,--quiet" help:"Only display error messages."`
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

	printer := printer.NewPrinter(args.Quiet)

	switch {
	case args.ModuleCmd != nil:
		m := &dotf.Module{Printer: printer, Cmd: args.ModuleCmd}

		if err := m.Run(printer); err != nil {
			printer.Println(fmt.Sprintf("Error: %v", err))
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
			switch choice {
			case "d":
				p.WriteHelp(os.Stdout)
			case "g":
				cmd.DotfExecute(stdinArgs)
			}
		}
	}

}
