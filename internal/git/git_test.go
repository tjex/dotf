package git_test

import (
	"testing"

	"git.sr.ht/~tjex/dotf/internal/git"
)

func TestModuleRemote(t *testing.T) {
	tests := []struct{ origin, suffix, path, repo, want string }{
		{"git@git.sr.ht:~tjex/dotfiles", "-rc", "/home/u/.config/nvim", "", "git@git.sr.ht:~tjex/nvim-rc"},
		{"git@git.sr.ht:~tjex/dotfiles", "-rc", "~/.local/share/timewarrior/", "", "git@git.sr.ht:~tjex/timewarrior-rc"},
		{"https://git.sr.ht/~tjex/dotfiles/", "", "~/notes", "", "https://git.sr.ht/~tjex/notes"},
		{"git@host:dotfiles", "-rc", "~/notes", "", "git@host:notes-rc"},
		{"git@git.sr.ht:~tjex/dotfiles", "-rc", "~/.local/share/task", "task-data-rc", "git@git.sr.ht:~tjex/task-data-rc"},
		{"git@git.sr.ht:~tjex/dotfiles", "-rc", "~/x", "https://github.com/me/x.git", "https://github.com/me/x.git"},
		{"git@git.sr.ht:~tjex/dotfiles", "-rc", "~/x", "git@github.com:me/x.git", "git@github.com:me/x.git"},
	}
	for _, tt := range tests {
		got, err := git.ModuleRemote(tt.origin, tt.suffix, tt.path, tt.repo)
		if err != nil || got != tt.want {
			t.Errorf("ModuleRemote(%q,%q,%q) = %q, %v; want %q", tt.origin, tt.suffix, tt.path, got, err, tt.want)
		}
	}
}
