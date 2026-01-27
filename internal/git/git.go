package git

import (
	"fmt"

	"git.sr.ht/~tjex/dotf/cmd"
)

// Returns whether local branch wants a pull / push
func SyncState(repoPath string) (bool, bool) {
	cmd.Cmd("git", []string{"-C", repoPath, "fetch", "--quiet"})

	fetch := []string{"-C", repoPath, "rev-list", "--left-right", "--count", "@{upstream}...HEAD"}
	out := cmd.Cmd("git", fetch)

	var wantsPull, wantsPush int
	fmt.Sscanf(out, "%d %d", &wantsPull, &wantsPush)

	return wantsPull > 0, wantsPush > 0
}

func UncommittedChanges(repo, worktree string) bool {
	out := cmd.Cmd("git", []string{"-C", repo, "--work-tree", worktree, "status", "--porcelain"})
	return len(out) > 0
}

func Cmd(args []string) string {
	out := cmd.Cmd("git", args)
	return out
}

func Push(repo string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "push"})
	return out
}

func Pull(repo string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "pull"})
	return out
}

func Commit(repo string, message *string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "commit", "-m", *message})
	return out
}

func AddAll(repo string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "add", "-A"})
	return out
}

func Status(repo string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "status", "--porcelain"})
	return out
}
