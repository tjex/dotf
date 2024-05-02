package log

import (
	"bytes"
	"fmt"
	"log"
)

var (
	buf    bytes.Buffer
	logger = log.New(&buf, "logger: ", log.Lshortfile)
)

func Print(v any) {
	logger.Print(v)
	fmt.Print(&buf)
}
