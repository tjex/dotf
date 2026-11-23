package git

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

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

// RunDotf streams Git diagnostics with repository context, including prompts.
func RunDotf(args []string) error {
	cmd := Dotf(args)
	stderr := &repoWriter{repo: config.UserConfig().GitDir, lineStart: true}
	cmd.Stderr = stderr
	err := cmd.Run()
	stderr.finish()
	return err
}

type repoWriter struct {
	repo      string
	lineStart bool
}

var (
	stderrMu   sync.Mutex
	stderrLast *repoWriter
)

func (w *repoWriter) Write(p []byte) (int, error) {
	stderrMu.Lock()
	defer stderrMu.Unlock()

	if stderrLast != nil && stderrLast != w {
		if _, err := fmt.Fprintln(os.Stderr); err != nil {
			return 0, err
		}
		stderrLast.lineStart = true
		stderrLast = nil
	}

	written := 0
	for len(p) > 0 {
		if w.lineStart {
			if _, err := fmt.Fprintf(os.Stderr, "(%s) ", w.repo); err != nil {
				return written, err
			}
			w.lineStart = false
		}
		end := bytes.IndexByte(p, '\n') + 1
		if end == 0 {
			end = len(p)
		}
		n, err := os.Stderr.Write(p[:end])
		written += n
		if err != nil {
			return written, err
		}
		if n != end {
			return written, io.ErrShortWrite
		}
		w.lineStart = p[end-1] == '\n'
		if w.lineStart {
			stderrLast = nil
		} else {
			stderrLast = w
		}
		p = p[end:]
	}
	return written, nil
}

func (w *repoWriter) finish() {
	stderrMu.Lock()
	defer stderrMu.Unlock()
	if stderrLast == w {
		_, _ = fmt.Fprintln(os.Stderr)
		w.lineStart = true
		stderrLast = nil
	}
}

// Returns whether local branch wants a pull/push
func SyncState(repo string) (bool, bool, error) {
	cmd := RepoCmd(repo, "fetch", "--quiet")
	var stderr bytes.Buffer
	writer := &repoWriter{repo: repo, lineStart: true}
	cmd.Stderr = io.MultiWriter(writer, &stderr)
	err := cmd.Run()
	writer.finish()
	if err != nil {
		if message := strings.TrimSpace(stderr.String()); message != "" {
			return false, false, fmt.Errorf("%s: %w", message, err)
		}
		return false, false, err
	}

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

// IsFullURL reports whether s is a complete git remote URL (git@host:..., git://,
// ssh:// or http(s)://) that must be used as-is.
func IsFullURL(s string) bool {
	for _, p := range []string{"git@", "git://", "ssh://", "https://", "http://"} {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// Clone clones url into dest. Missing parent directories are created by git.
func Clone(url, dest string) (string, error) {
	return proc.Output("git", []string{"clone", url, dest})
}

// BranchExists reports whether branch exists in repo as a local branch or as
// a branch on origin.
func BranchExists(repo, branch string) bool {
	return branchExists([]string{"-C", repo}, branch)
}

// BareBranchExists is BranchExists for a bare repository, which git only
// accepts when given explicitly via --git-dir.
func BareBranchExists(gitDir, branch string) bool {
	return branchExists([]string{"--git-dir", gitDir}, branch)
}

func branchExists(base []string, branch string) bool {
	for _, ref := range []string{"refs/heads/", "refs/remotes/origin/"} {
		args := append(append([]string{}, base...), "rev-parse", "--verify", "--quiet", ref+branch)
		if _, err := proc.Output("git", args); err == nil {
			return true
		}
	}
	return false
}

// ModuleRemote derives a module's remote URL from the bare repo's origin.
// The origin's last path segment (the bare repo name) is replaced by the
// repository name. A repository that is a full URL is returned unchanged.
// If repository is empty, the name is the module directory
// name plus suffix, e.g.
// ("git@git.sr.ht:~me/dotfiles", "-rc", "~/.config/nvim", "") -> "git@git.sr.ht:~me/nvim-rc".
func ModuleRemote(origin, suffix, modulePath, repository string) (string, error) {
	if IsFullURL(repository) {
		return repository, nil
	}
	origin = strings.TrimRight(origin, "/")
	i := strings.LastIndexAny(origin, "/:")
	if i < 0 {
		return "", fmt.Errorf("cannot derive module remote from origin %q", origin)
	}
	name := repository
	if name == "" {
		base := filepath.Base(filepath.Clean(modulePath))
		if base == "." || base == "/" {
			return "", fmt.Errorf("cannot derive repository name from %q", modulePath)
		}
		name = base + suffix
	}
	return origin[:i+1] + name, nil
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
