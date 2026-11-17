# dotf Changelog

## v0.6.4

- Add a GitHub Actions workflow to publish Linux and macOS binaries for amd64 and
  arm64 when a version tag is pushed.

## v0.6.3

- Add `dotf clone --branch` to select a branch for the bare repository and modules,
  falling back to the default branch with a warning if it does not exist.
- Configure cloned bare repositories to hide untracked files.
- Apply the configured concurrency batch size to module cloning.

## v0.6.2

- Rename the module object key `repository` to `repo`. Existing module objects
  must use the new key.

## v0.6.1

- Allow fully qualified Git URLs as module repository overrides.

## v0.6.0

- Add `dotf clone` for cloning the bare repository, modules, or both, and
  `dotf m --clone` for cloning modules missing on disk.
- Derive module remotes from `origin`, the module directory name, and
  `repo-suffix`.
- Support module entries as strings or `{ path, repository }` objects with
  explicit repository overrides.
- Fix output from concurrent operations.

## v0.5.2

- Fixed module operations hanging if concurrency was unset in config or set to <= 0
- Display absolute paths in `dotf m --list`.

## v0.5.1

- Improve concurrency batch handling.
- Removed concurrency sleep duration, as new worker group handling suffices.

## v0.5.0

- Add config options to control concurrent module operations (batching and
  delaying)
- Attach interactive Git processes to the terminal so commands such as
  `rebase --interactive` do not hang.

## v0.4.6

- Add missing error handling in modules.prime()

## v0.4.5

- Require an explicit `--all`, `--bare`, or `--modules` flag for `dotf sync`.
- Use go-arg's built-in help and version handling.

## v0.4.4

- Check Git command errors using exit status.
- Improve terminal output.

## v0.4.3

- Only push the bare repository when it has commits to push.

## v0.4.2

- Propagate errors from bare repository operations.

## v0.4.1

- Propagate Git command errors and improve error handling in concurrent module
  operations.

## v0.4.0

- Add `dotf m --status` to display short Git status for modules.

## v0.3.7

- Improve sync-state detection and bare repository sync checks.
- Improve sync handling and terminal output.

## v0.3.6

- Include the prime step when syncing modules and the bare repository.
- Improve terminal output.

## v0.3.5

- Support environment variable and tilde expansion in `git-dir` and `worktree`
  config paths.

## v0.3.4

- Add `dotf sync` to pull and push the bare repository, modules, or both.
- Add version output through `--version` and `-v`.
- Add a SourceHut build job to mirror the repository.

## v0.3.3

- Run module subcommands concurrently using goroutines.

## v0.3.2

- Fall back to the default module group when `DOTF_ENV` is unset.

## v0.3.1

- Add `--quiet` to suppress operation output.
- Add more detailed logging for module pushes and pulls.

## v0.3.0

- Add named module groups selected through `DOTF_ENV`, merged with the default
  module group.
- Refresh upstream state and fix detection of whether repositories need pushing
  or pulling.
- Print operation logs to stdout.
- Update the required Go version to 1.25.

## v0.2.2

- Fix module push and pull detection when more than one commit is pending.
- Return stderr from failed Git commands.

## v0.2.1

- Replace submodule-specific operations with modules: independent repositories
  listed in the configuration.
- Add module listing, editing, pushing, pulling, and batch staging and committing
  through `dotf m`.
- Fix Git command handling and module priming.

## v0.2.0

- Remove mirroring functionality completely.

## v0.1.4

- Added the 'mode' config option (to push to mirror with `--mirror` or not).

## v0.1.3

- Fix handling of relative submodule paths.

## v0.1.2

- fix: `dotf sm --edit` returned a nil pointer.
- Add `dotf sm --edit` to select a submodule with fzf and open it in an editor.
- Handle cancellation and exit errors from fzf.

## v0.1.1

- Stop adding extra quotes to configured batch commit messages.
- Add installation instructions and SourceHut build configuration.

## v0.1.0

- Initial release with a configurable bare repository Git wrapper.
- Add submodule listing and batch staging and committing with `dotf sm`.
- Push to the origin and mirror concurrently.
