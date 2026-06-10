package git

import (
	"fmt"
	"os/exec"

	"git.sr.ht/~tjex/dotf/cmd"
	"git.sr.ht/~tjex/dotf/internal/config"
)

// Base helper to run a git command with -C <repo> automatically
func RepoCmd(repo string, args ...string) *exec.Cmd {
	fullArgs := append([]string{"-C", repo}, args...)
	return cmd.Cmd("git", fullArgs)
}

func RepoCmdRoutine(repo string, args ...string) (string, error) {
	fullArgs := append([]string{"-C", repo}, args...)
	return cmd.CmdRoutine("git", fullArgs)
}

// Full error-returning helper
func Cmd(args []string) *exec.Cmd {
	return cmd.Cmd("git", args)
}

// Git command with flags set as per the user's config
func Dotf(args []string) *exec.Cmd {
	bareRepoArgs := buildArgsArray(args)
	return cmd.Cmd("git", bareRepoArgs)
}

// Git command with flags set as per the user's config
func DotfInstance(args []string) *exec.Cmd {
	bareRepoArgs := buildArgsArray(args)
	return cmd.Cmd("git", bareRepoArgs)
}

// Returns whether local branch wants a pull/push
func SyncState(repo string) (bool, bool, error) {
	cmd := RepoCmd(repo, "fetch", "--quiet")
	cmd.Run()

	out, err := RepoCmdRoutine(repo, "rev-list", "--left-right", "--count", "@{upstream}...HEAD")
	if err != nil {
		return false, false, err
	}

	var wantsPull, wantsPush int
	fmt.Sscanf(out, "%d %d", &wantsPull, &wantsPush)

	return wantsPull > 0, wantsPush > 0, nil
}

func UncommittedChanges(repo, worktree string) (bool, error) {
	out, err := cmd.CmdRoutine("git", []string{"-C", repo, "--work-tree", worktree, "status", "--porcelain"})
	if err != nil {
		return false, err
	}
	return len(out) > 0, nil
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
