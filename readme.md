# dotf

A `git` wrapper to make dotfile tracking with a bare git repository even more
convenient.

## Config

`dotf` expects to find a config file at `${XDG_CONFIG_HOME}/dotf/config.toml`.

The below settings are for demonstration in keeping with the
[bare repo dotfiles tutorial on Atlassian](https://www.atlassian.com/git/tutorials/dotfiles)

```toml
worktree = "/Users/<user>" # note: must be absolute path (no $HOME, ~/, etc.. yet)
gitdir = "/Users/<user>/.cfg/" # the bare git repo root
origin = "<read+write url for origin>" # eg, git@yourhost.com:user/dotfiles
mirror = "read+write url for mirror" # eg, git@yourmirror.com:user/dotfiles
batch-commit-message = "batch dotf update" # used by `dotf prime` for submodule commit message
```

## Usage

All `git` commands are passed as normal. Some are intercepted and handled
differently, some are unique:

`prime` - add (with `git add -A`) and commit all changes to all submodules.
Commit message is set in `config.toml`.

`push` - push to origin and mirror concurrently.

A regular workflow would then look like the following. From anywhere in your
file system:

```bash
dotf prime # add and commit all changes within submodules
dotf add -u # add all changes to tracked dotfiles
dotf commit -m "update all dotfiles"

dotf push
```

Run all other git commands as normal:

```bash
dotf status
dotf log --oneline
dotf rebase -i HEAD~2
...
```

