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
git-dir = "~/.cfg/" # the location of the bare repo git directory (i.e, git config, hooks, etc).
origin = "<read+write url for origin>" # eg, git@yourhost.com:user/dotfiles.
repo-suffix = "-rc" # appended to a module's directory name to get its repository name.
batch-commit-message = "batch dotf update" # used by `dotf m --prime` for module commit message.

[concurrency]
batch-size = 15 # How many concurrent git operations to run per batch.
                # This helps if your git host rate limits or otherwise restricts your connections.

[modules.default]
paths = ["~/path/to/repo1", { path = "~/path/to/repo2", repository = "repo2-rc" }]
```

## Installation

1. Clone this repo
   1. To install with the latest changes: `make install`
   2. To install with a stable version:
      `git checkout tags/<version> && make install`

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

### Module remotes and cloning

`dotf m --clone` (or `dotf clone --modules`) clones every module that does not
yet exist on disk. Each module's remote is derived from `origin`: the last
element of `origin` (the bare repo name) is replaced by the module's repository
name.

A module entry in `paths` is either a plain string or an object:

```toml
origin = "git@git.sr.ht:~me/dotfiles"
repo-suffix = "-rc"

[modules.default]
paths = [
  "~/.config/nvim",                                                   # -> git@git.sr.ht:~me/nvim-rc
  { path = "~/.config/tmux" },                                        # -> git@git.sr.ht:~me/tmux-rc
  { path = "~/.config/task", repository = "task-config-rc" },         # -> git@git.sr.ht:~me/task-config-rc
  { path = "~/.local/share/task", repository = "task-data-rc" },      # -> git@git.sr.ht:~me/task-data-rc
  { path = "~/.config/foo", repository = "https://github.com/me/foo.git" }, # -> used as-is
]
```

- Without `repository`, the name is the path's last element plus `repo-suffix`.
- With `repository`, the name is used exactly as written. `repo-suffix` is _not_
  appended.
- If `repository` starts with `git@`, `git://`, `ssh://`, `https://` or
  `http://`, it is treated as a complete remote URL and cloned from as-is.
  `origin` is ignored for that module, so it can live on any host.
- Set `repository` whenever two modules share a last path element (like the two
  `task` directories above), otherwise they would map to the same remote.
- Modules are expected to be hosted next to the bare repository, as `origin`
  provides the host and user.
- `path` is required in an object.

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

# Clone modules missing on disk (remote derived from origin, see above)
dotf m --clone
```

### Bare Repo

```bash
# Pull bare dotfiles repository from origin
dotf pull

# Push bare dotfiles repository to origin
dotf push
```

### Clone

```bash
# Clone the bare repository from origin
dotf clone --bare

# Clone all modules missing on disk
dotf clone --modules

# Clone the default branch, then check out <name> (works with --bare, --modules and --all)
dotf clone --modules --branch server

# Clone both
dotf clone --all
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
