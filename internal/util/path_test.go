package util_test

import (
	"os"
	"testing"

	"git.sr.ht/~tjex/dotf/internal/util"
)

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path string
		want string
	}{
		{
			"Expand ${HOME}",
			"${HOME}",
			home,
		},
		{
			"Expand tilde",
			"~",
			home,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := util.ExpandPath(tt.path)
			if gotErr != nil {
				t.Errorf("ExpandPath() failed: %v", gotErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExpandPath() = %q, want %q", got, tt.want)
			}
		})
	}
}
