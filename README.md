# dotf

A `git` wrapper to make tracking dotfiles with a bare git repository even more
convenient.

## Config

`dotf` expects to find a config file at `${XDG_CONFIG_HOME}/dotf/config.toml`.

The below settings are for demonstration in keeping with the
[bare repo dotfiles tutorial on Atlassian](https://www.atlassian.com/git/tutorials/dotfiles).

Go through the tutorial first, if the concept of tracking system configurations
with a bare git repostiory is new to you.

```toml
env = "example" # or `export DOTF_ENV=example` in shell config.
worktree = "${HOME}" # the bare repo's root directory (all files deeper can be tracked).
gitdir = "~/.cfg/" # the location of the bare repo git directory (i.e, git config, hooks, etc).
origin = "<read+write url for origin>" # eg, git@yourhost.com:user/dotfiles.
batch-commit-message = "batch dotf update" # used by `dotf m --prime` for module commit message.

[modules.default]
paths = ["~/path/to/repo1", "~/path/to/repo2"]
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

The `modules.default` group will _always_ be included. You can use this as a
base group of modules used across systems. Any repos specific to certain systems
(e.g., work and personal) can be included in their own groups.

A non-default grouped module will only be included if the `DOTF_ENV` shell
variable or the `env` key in the `config.toml` are set. This allows you your
`dotf` config to be portable across systems while still enabling alternate
module setups.

If you do not want to have any modules included by default, simply don't use the
`modules.default` group.

```bash
export DOTF_ENV=work
```

```toml
[modules.default] # present on all systems you work on
paths = ["~/path/to/repo1", "~/path/to/repo2"]

[modules.work]
paths = ["~/path/to/work/repo1"]

[modules.personal] # will not be included on work computer
paths = ["~/path/to/personal/repo1"]
```

Submodules are still available for you to use as normal if you prefer:
`dotf submodule add ...`, etc. As all Git commands apart from those intercepted
in this program are _passed to your bare dotfiles Git repository_.

## Usage

All `git` commands are passed as normal. Some are intercepted and handled
differently, some are unique. `dotf m ...` commands range over module paths in
your `config.toml`, and run as goroutines:


### Modules

```bash
# List all tracked modules
dotf m --list

# Search module directories using fzf, open selected in $EDITOR (default: vim)
dotf m --edit

# Check status of modules `git status -s`
dotf m --status

# Pull upstream changes from all modules
dotf m --pull

# Stage all changes (`git add -A`) and commit them for all modules
# Commit message is defined in config.toml
dotf m --prime

# Push local changes of all modules
dotf m --push
```

### Bare Repo

```bash
# Pull bare dotfiles repository from origin
dotf pull

# Push bare dotfiles repository to origin
dotf push
```

### Sync

```bash
# Sync bare repository: pull and push
dotf sync --bare

# Sync all modules: pull and push
dotf sync --modules

# Sync both bare repo and modules
dotf sync
```

### Misc

```bash
# Run commands quietly (errors only)
dotf ... -q

# Show interactive help for dotf or git
dotf --help

# All other commands will be passed as-is to Git (operating on the bare repo).

dotf status
dotf log --oneline
dotf rebase -i HEAD~2
...
```

## Author

Tillman Jex\
[www.tjex.net](https://tjex.net)
