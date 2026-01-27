package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/cmd/sync"
	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/git"
	"git.sr.ht/~tjex/dotf/internal/printer"
	"github.com/alexflint/go-arg"
)

var args struct {
	ModuleCmd *dotf.ModuleCmd `arg:"subcommand:m" help:"operations for git modules."`
	SyncCmd   *sync.SyncCmd   `arg:"subcommand:sync" help:"Sync remote and local with latest changes (i.e, pull all -> push all)."`
	Quiet     bool            `arg:"-q,--quiet" help:"Only display error messages."`
}

var (
	Version = "dev"
)

func main() {
	err := config.ReadUserConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

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
	case args.SyncCmd != nil:
		s := &sync.Sync{Printer: printer, Cmd: args.SyncCmd}

		if err := s.Run(printer); err != nil {
			printer.Println(fmt.Sprintf("Error: %v", err))
			p.WriteHelp(os.Stdout)
		}
	case args.ModuleCmd != nil:
		m := &dotf.Module{Printer: printer, Cmd: args.ModuleCmd}

		if err := m.Run(printer); err != nil {
			printer.Println(fmt.Sprintf("Error: %v", err))
			p.WriteHelp(os.Stdout)
		}
	default:
		// If arguments above aren't called, nor is `--help`
		if os.Args[1] == "--version" || os.Args[1] == "-v" {
			v, _ := strings.CutPrefix(Version, "v")
			fmt.Println("dotf version", v)
			out := git.Dotf(stdinArgs) // prints git version
			printer.Println(out)
		} else if os.Args[1] != "--help" && os.Args[1] != "-h" {
			out := git.Dotf(stdinArgs)
			printer.Println(out)
		} else {
			var choice string
			fmt.Println("dotf wraps around git. \nDisplay help for dotf (d) or git (g)?")
			fmt.Scan(&choice)
			switch choice {
			case "d":
				p.WriteHelp(os.Stdout)
			case "g":
				out := git.Dotf(stdinArgs)
				printer.Println(out)
			}
		}
	}

}
