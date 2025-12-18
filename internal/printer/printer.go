package printer

import (
	"fmt"
)

type Printer struct {
	Quiet bool
}

func NewPrinter(quiet bool) *Printer {
	return &Printer{Quiet: quiet}
}

func (p *Printer) Println(v... any) {
	if !p.Quiet {
		fmt.Println(v...)
	}
}
