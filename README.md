# dotf

A `git` wrapper to make tracking dotfiles with a bare git repository even more
convenient.

## Config

`dotf` expects to find a config file at `${XDG_CONFIG_HOME}/dotf/config.toml`.

The below settings are for demonstration in keeping with the
[bare repo dotfiles tutorial on Atlassian](https://www.atlassian.com/git/tutorials/dotfiles)

```toml
worktree = "/Users/<user>" # note: must be absolute path (no $HOME, ~/, etc.. yet)
gitdir = "/Users/<user>/.cfg/" # the bare git repo root
origin = "<read+write url for origin>" # eg, git@yourhost.com:user/dotfiles
batch-commit-message = "batch dotf update" # used by `dotf m --prime` for module commit message
modules = ["~/path/to/repo1", "~/path/to/repo2"] # like git submodles, but not a part of  the bare repo
```

For example, my config
[is here](https://git.sr.ht/~tjex/dotfiles/tree/mac/item/.config/dotf/config.toml).

## Installation

1. Clone this repo
   1. To install with the latest changes: `go install`
   2. To install with a stable version:
      `git checkout tags/<version> && go install`

To get version tags, first fetch all tags with `git fetch --all --tags` and then
run `git tag`.

This repo is also [mirrored on GitHub](https://github.com/tjex/dotf).

## Usage

All `git` commands are passed as normal. Some are intercepted and handled
differently, some are unique:

```markdown
`dotf m --prime`: add (with `git add -A`) and commit all changes to all
modules. Commit message is set in `config.toml`.

`dotf m --list`: list all tracked modules.

`dotf m --edit`: search module directories with `fzf`, opening selected with
$EDITOR (defaults to vim).

`dotf push`: push to origin.

`dotf --help`: display help for dotf or git (interactively)
```

All flags have shorthand as well: `--prime` / `-p`.


Run all other git commands as normal:

```bash
dotf status
dotf log --oneline
dotf rebase -i HEAD~2
...

```

## TODO

- Module pull
- Module push

## Author

Tillman Jex \
[www.tjex.net](https://tjex.net)
