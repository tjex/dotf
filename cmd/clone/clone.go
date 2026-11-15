package clone

import (
	"errors"

	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/internal/printer"
)

type CloneCmd struct {
	CloneAll     bool `arg:"-a,--all" help:"Clone both the bare repository and all modules."`
	CloneBare    bool `arg:"-b,--bare" help:"Clone the bare repository."`
	CloneModules bool `arg:"-m,--modules" help:"Clone all modules."`
}

type Clone struct {
	Printer *printer.Printer
	Cmd     *CloneCmd
}

func (c *Clone) Run() error {
	switch {
	case c.Cmd.CloneAll:
		return c.cloneAll()
	case c.Cmd.CloneBare:
		return c.cloneBare()
	case c.Cmd.CloneModules:
		return c.cloneModules()
	default:
		return errors.New("No flag passed to clone.")
	}
}

func (c *Clone) cloneAll() error {
	err := c.cloneBare()
	if err != nil {
		return err
	}

	err = c.cloneModules()
	if err != nil {
		return err
	}

	return nil
}

func (c *Clone) cloneBare() error {
	c.Printer.Println("Cloning bare repository...")
	bare := &dotf.Bare{Printer: c.Printer}
	return bare.Clone()
}

func (c *Clone) cloneModules() error {
	c.Printer.Println("Cloning modules...")
	modules := &dotf.Modules{Printer: c.Printer}
	return modules.Clone()
}
