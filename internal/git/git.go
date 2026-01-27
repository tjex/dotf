package git

import (
	"fmt"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

func Cmd(args []string) string {
	out := cmd.Cmd("git", args)
	return out
}

// Git command with flags set as per the user's config
func Dotf(args []string) string {
	bareRepoArgs := buildArgsArray(args)
	out := cmd.Cmd("git", bareRepoArgs)

	return out
}

// Returns whether local branch wants a pull / push
func SyncState(repo string) (bool, bool) {
	cmd.Cmd("git", []string{"-C", repo, "fetch", "--quiet"})

	fetch := []string{"-C", repo, "rev-list", "--left-right", "--count", "@{upstream}...HEAD"}
	out := cmd.Cmd("git", fetch)

	var wantsPull, wantsPush int
	fmt.Sscanf(out, "%d %d", &wantsPull, &wantsPush)

	return wantsPull > 0, wantsPush > 0
}

func UncommittedChanges(repo, worktree string) bool {
	out := cmd.Cmd("git", []string{"-C", repo, "--work-tree", worktree, "status", "--porcelain"})
	return len(out) > 0
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
	out := cmd.Cmd("git", []string{"-C", repo, "-c", "color.ui=always", "status", "-s"})
	return out
}

// Build the arguments array for dotf git call
func buildArgsArray(gitArgs []string) []string {
	conf := config.UserConfig()
	var argsArray []string

	// NOTE: Color flag needs to be set here for correct flag ordering
	argsArray = append(argsArray, "-c", "color.ui=always")
	argsArray = append(argsArray, conf.RepoFlags...)
	argsArray = append(argsArray, gitArgs...)

	return argsArray
}
