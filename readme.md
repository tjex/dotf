# dotf

A little git wrapper to make some dotfile operations more convenient.

## config

`dotf` looks for configuration in `${XDG_CONFIG_HOME}/dotf/config.toml`.

As an example, these settings are in keeping with the
[bare repo dotfiles tutorial on Atlassian](https://www.atlassian.com/git/tutorials/dotfiles)

```toml
worktree = "/Users/<user>" # note: must be absolute path (and no $HOME, etc.. yet)
gitdir = "/Users/<user>/.cfg/"
mirror = "ssh read/write url for mirror"
origin = "<ssh read/write url for origin>" # eg, git@git.sr.ht:~tjex/dotfiles
batch-commit-message = "batch dotf update" # used by `dotf prime`
```

## usage

All `git` commands are passed as normal. Some are intercepted and handled
differently, some are unique:

`push` - push to origin and mirror concurrently. `prime` - add (via
`git add -u`) and commit all changes to all submodules. Commit message is set in
`config.toml`

A regular workflow would then look like the following. From anywhere in your
file system:

```bash
dotf add -u
dotf commit -m "add all unstaged changes to tracked files"

dotf prime

dotf push
```
