package git

import (
	"regexp"
	"git.sr.ht/~tjex/dotf/cmd"
)

// Returns whether local branch wants a push / pull
func SyncState(repoPath string) (bool, bool) {
	args := []string{"-C", repoPath, "rev-list", "--left-right", "--count", "HEAD...@{u}"} 
	out := cmd.Cmd("git", args)

	wantsPush, _ := regexp.MatchString("1.*", out)
	wantsPull, _ := regexp.MatchString(".*1", out)

	return wantsPush, wantsPull
}

func Cmd(args []string) string {
	out := cmd.Cmd("git", args)
	return out
}

func Push(repo string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "push"})
	return out
}
