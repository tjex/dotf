package git

import (
	"fmt"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

// Base helper to run a git command with -C <repo> automatically
func gitRepoCmd(repo string, args ...string) (string, error) {
	fullArgs := append([]string{"-C", repo}, args...)
	return cmd.Cmd("git", fullArgs)
}

// Full error-returning helper
func Cmd(args []string) (string, error) {
	return cmd.Cmd("git", args)
}

// Git command with flags set as per the user's config
func Dotf(args []string) (string, error) {
	bareRepoArgs := buildArgsArray(args)
	return cmd.Cmd("git", bareRepoArgs)
}

// Returns whether local branch wants a pull/push
func SyncState(repo string) (bool, bool, error) {
	_, _ = gitRepoCmd(repo, "fetch", "--quiet")

	out, err := gitRepoCmd(repo, "rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	if err != nil {
		return false, false, err
	}

	var wantsPull, wantsPush int
	fmt.Sscanf(out, "%d %d", &wantsPull, &wantsPush)

	return wantsPull > 0, wantsPush > 0, nil
}

func UncommittedChanges(repo, worktree string) (bool, error) {
	out, err := cmd.Cmd("git", []string{"-C", repo, "--work-tree", worktree, "status", "--porcelain"})
	if err != nil {
		return false, err
	}
	return len(out) > 0, nil
}

func Push(repo string) (string, error) {
	return gitRepoCmd(repo, "push")
}

func Pull(repo string) string {
	out, _ := gitRepoCmd(repo, "pull")
	return out
}

func Commit(repo string, message *string) string {
	out, _ := gitRepoCmd(repo, "commit", "-m", *message)
	return out
}

func AddAll(repo string) string {
	out, _ := gitRepoCmd(repo, "add", "-A")
	return out
}

func Status(repo string) string {
	out, _ := gitRepoCmd(repo, "-c", "color.ui=always", "status", "-s")
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
