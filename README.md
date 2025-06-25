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

## Modules

Modules are a replacement for Git's submodules feature. I found that adding
extra repositories to my bare repo to be frustrating at times as it requires the
`git submodules for each ...` way of working.

Adding submodules also creates strict relationships between the submodule and
the bare repository. Once this broke my config an was very difficult to fix.
This is a little scary when the bare repository is in control of a large portion
of your systems configuration.

Instead `modules` are regular git repositories, in the sense that they have no
relationship to your bare dotfiles respository. You can clone a repo somewhere
on your system, add its path to the `modules` array in the `config.toml` and
then list, edit, pull, commit and push with the `m` command: `dotf m --list`.

Submodules are therefore still available for you to use as normal: `dotf
submodule add ...`, etc. As all Git commands apart from those intercepted in
this program are _passed to your bare dotfiles Git repository_.

## Usage

All `git` commands are passed as normal. Some are intercepted and handled
differently, some are unique:

```markdown

`dotf m --list`: list all tracked modules.

`dotf m --edit`: search module directories with `fzf`, opening selected with
$EDITOR (defaults to vim).

`dotf m --pull`: pull upstream changes from all modules.

`dotf m --commit`: add (with `git add -A`) and commit all changes to all
modules. Commit message is set in `config.toml`.

`dotf m --push`: push local changes of all modules.

`dotf pull`: pull bare dotfiles repository from origin.

`dotf push`: push bare dotfiles repository to origin.

`dotf --help`: display help for dotf or git (interactively)
```

All flags have shorthand as well: `--commit` / `-c`.


Run all other git commands as normal:

```bash
dotf status
dotf log --oneline
dotf rebase -i HEAD~2
...

```

## Author

Tillman Jex \
[www.tjex.net](https://tjex.net)
