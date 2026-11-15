package sync

import (
	"errors"

	"git.sr.ht/~tjex/dotf/cmd/dotf"
	"git.sr.ht/~tjex/dotf/internal/printer"
)

type SyncCmd struct {
	SyncAll     bool `arg:"-a,--all" help:"Sync both the bare repository and all modules."`
	SyncBare    bool `arg:"-b,--bare" help:"Sync the bare repository."`
	SyncModules bool `arg:"-m,--modules" help:"Sync all modules."`
}

type Sync struct {
	Printer *printer.Printer
	Cmd     *SyncCmd
}

func (s *Sync) Run() error {
	switch {
	case s.Cmd.SyncAll:
		return s.syncAll()
	case s.Cmd.SyncBare:
		return s.syncBare()
	case s.Cmd.SyncModules:
		return s.syncModules()
	default:
		return errors.New("No flag passed to sync.")
	}
}

func (s *Sync) syncAll() error {
	err := s.syncBare()
	if err != nil {
		return err
	}

	err = s.syncModules()
	if err != nil {
		return err
	}

	return nil
}

func (s *Sync) syncBare() error {
	s.Printer.Println("Syncing bare repository...")
	bare := &dotf.Bare{Printer: s.Printer}
	return bare.Sync()
}

func (s *Sync) syncModules() error {
	s.Printer.Println("Syncing modules...")
	modules := &dotf.Modules{Printer: s.Printer}

	if err := modules.Pull(); err != nil {
		return err
	}
	if err := modules.Prime(); err != nil {
		return err
	}
	return modules.Push()
}
