package git

import (
	"fmt"
	"os/exec"

	"git.sr.ht/~tjex/dotf/internal/config"
	"git.sr.ht/~tjex/dotf/internal/proc"
)

// Base helper to run a git command with -C <repo> automatically
func RepoCmd(repo string, args ...string) *exec.Cmd {
	fullArgs := append([]string{"-C", repo}, args...)
	return proc.Stream("git", fullArgs)
}

func RepoCmdRoutine(repo string, args ...string) (string, error) {
	fullArgs := append([]string{"-C", repo}, args...)
	return proc.Output("git", fullArgs)
}

// Plain git command with no repo flags
func Cmd(args []string) *exec.Cmd {
	return proc.Stream("git", args)
}

// Git command with flags set as per the user's config
func Dotf(args []string) *exec.Cmd {
	bareRepoArgs := buildArgsArray(args)
	return proc.Stream("git", bareRepoArgs)
}

// Returns whether local branch wants a pull/push
func SyncState(repo string) (bool, bool, error) {
	RepoCmd(repo, "fetch", "--quiet").Run()

	out, err := RepoCmdRoutine(repo, "rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	if err != nil {
		return false, false, err
	}

	var wantsPull, wantsPush int
	fmt.Sscanf(out, "%d %d", &wantsPull, &wantsPush)

	return wantsPull > 0, wantsPush > 0, nil
}

func UncommittedChanges(repo, worktree string) (bool, error) {
	out, err := proc.Output("git", []string{"-C", repo, "--work-tree", worktree, "status", "--porcelain"})
	if err != nil {
		return false, err
	}
	return len(out) > 0, nil
}

func Clone(repo string) (string, error) {
	return RepoCmdRoutine(repo, "clone")
}

func Push(repo string) (string, error) {
	return RepoCmdRoutine(repo, "push")
}

func Pull(repo string) (string, error) {
	return RepoCmdRoutine(repo, "pull")
}

func Commit(repo string, message string) (string, error) {
	return RepoCmdRoutine(repo, "commit", "-m", message)
}

func AddAll(repo string) (string, error) {
	return RepoCmdRoutine(repo, "add", "-A")
}

func Status(repo string) (string, error) {
	return RepoCmdRoutine(repo, "-c", "color.ui=always", "status", "-s")
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
