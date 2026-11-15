package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"git.sr.ht/~tjex/dotf/cmd/clone"
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
	CloneCmd  *clone.CloneCmd `arg:"subcommand:clone" help:"Clone the bare repository and/or all modules."`
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
		os.Exit(1)
	}

	printer := printer.NewPrinter(args.Quiet)
	err = p.Parse(stdinArgs)
	var runner interface{ Run() error }
	var name string
	switch {
	case args.SyncCmd != nil:
		runner, name = &sync.Sync{Printer: printer, Cmd: args.SyncCmd}, "sync"
	case args.CloneCmd != nil:
		runner, name = &clone.Clone{Printer: printer, Cmd: args.CloneCmd}, "clone"
	case args.ModuleCmd != nil:
		runner, name = &dotf.Modules{Printer: printer, Cmd: args.ModuleCmd}, "m"
	}

	if runner != nil {
		if err := runner.Run(); err != nil {
			printer.Println(fmt.Sprintf("Error: %v", err))
			p.WriteHelpForSubcommand(os.Stdout, name)
		}
		return
	}

	switch {
	case err == arg.ErrHelp: // Means user passed: --help or -h
		p.WriteHelp(os.Stdout)
	case err == arg.ErrVersion: // Means user passed: --version or -v
		// FIXME: Not working. Passing --version just fals to default case.
		v, _ := strings.CutPrefix(Version, "v")
		fmt.Println("dotf version", v)
	default:
		cmd := git.Dotf(stdinArgs)
		cmd.Run()
	}

}
