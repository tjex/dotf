package git

import (
	"regexp"

	"git.sr.ht/~tjex/dotf/cmd"
)

// Returns whether local branch wants a push / pull
func SyncState(repoPath string) (bool, bool) {
	args := []string{"-C", repoPath, "rev-list", "--left-right", "--count", "HEAD...@{u}"} 
	out := cmd.Cmd("git", args)

	wantsPush, _ := regexp.MatchString(`^([1-9]\d*)\s+[0-9]\d*`, out)
	wantsPull, _ := regexp.MatchString(`^[0-9]\d*\s+([1-9]\d*)`, out)

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

func Pull(repo string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "pull"})
	return out
}

func Commit(repo string , message *string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "commit", "-m", *message})
	return out
}


func Add(repo string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "add", "-A"})
	return out
}


func Status(repo string) string {
	out := cmd.Cmd("git", []string{"-C", repo, "status", "--porcelain"})
	return out
}
