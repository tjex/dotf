package git_test

import (
	"testing"

	"git.sr.ht/~tjex/dotf/internal/git"
)

func TestModuleRemote(t *testing.T) {
	tests := []struct{ origin, suffix, path, want string }{
		{"git@git.sr.ht:~tjex/dotfiles", "-rc", "/home/u/.config/nvim", "git@git.sr.ht:~tjex/nvim-rc"},
		{"git@git.sr.ht:~tjex/dotfiles", "-rc", "~/.local/share/timewarrior/", "git@git.sr.ht:~tjex/timewarrior-rc"},
		{"https://git.sr.ht/~tjex/dotfiles/", "", "~/notes", "https://git.sr.ht/~tjex/notes"},
		{"git@host:dotfiles", "-rc", "~/notes", "git@host:notes-rc"},
	}
	for _, tt := range tests {
		got, err := git.ModuleRemote(tt.origin, tt.suffix, tt.path)
		if err != nil || got != tt.want {
			t.Errorf("ModuleRemote(%q,%q,%q) = %q, %v; want %q", tt.origin, tt.suffix, tt.path, got, err, tt.want)
		}
	}
}
