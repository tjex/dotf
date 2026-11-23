package git

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestRepositoryStderrKeepsContextWhenWritesInterleave(t *testing.T) {
	read, write, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = write
	t.Cleanup(func() {
		os.Stderr = old
		_ = read.Close()
		_ = write.Close()
	})

	first := &repoWriter{repo: "/first", lineStart: true}
	second := &repoWriter{repo: "/second", lineStart: true}
	_, _ = first.Write([]byte("fatal: first "))
	_, _ = second.Write([]byte("fatal: second\n"))
	_, _ = first.Write([]byte("continued\n"))
	_ = write.Close()
	output, err := io.ReadAll(read)
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		if !strings.HasPrefix(line, "(/first) ") && !strings.HasPrefix(line, "(/second) ") {
			t.Errorf("line lacks repository context: %q", line)
		}
	}
	if !strings.Contains(string(output), "(/first) continued") {
		t.Errorf("continuation lacks repository context: %s", output)
	}
}
