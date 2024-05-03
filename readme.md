# dotf

A little git wrapper to make some dotfile operations more convenient.

## config

Configuration is done by editing the `config.toml` file in `internal/config`.
Currently `dotf` gets built with the configuration you set.

As an example, these settings are in keeping with the
[bare repo dotfiles tutorial on Atlassian](https://www.atlassian.com/git/tutorials/dotfiles)

`/internal/config/config.toml`

```toml
worktree = "--work-tree=/Users/<user>" # note: must be absolute path (and no $HOME, etc)
gitdir = "--git-dir=/Users/<user>/.cfg/"
mirror = "ssh read/write url for mirror"
origin = "<ssh read/write url for origin>" # eg, git@git.sr.ht:~tjex/dotfiles
remoteName = "<remote name>" # eg, origin, remote, etc
```

## usage

`dotf push` will push to origin and mirror concurrently.

All other regular `git` commands will be passed to git as normal.

## that's it?

For now yes. But I want to implement a concurrent submodule sync feature.

Note: This is as much a learning task as it is a valuable tool...
