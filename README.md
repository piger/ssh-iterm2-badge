# ssh-iterm2-badge

The purpose of this program is to display the hostname of an [ssh](https://www.openssh.com/)
connection as a [badge](https://iterm2.com/documentation-badges.html) in [iTerm2](https://iterm2.com/).

**WARNING**: this is a gigantic hack and I don't really recommend using it.

For an explanation of how this program works you can check my blog post about it
[here](https://blog.spatof.org/posts/2023/05/kqueue-iterm2-and-openssh/).

## Requirements

- macOS
- iTerm2
- macOS' builtin OpenSSH

## Installation

To install this program from its source code you need to have the Go runtime installed; then you can
use `go install` as usual:

```shell
go install github.com/piger/ssh-iterm2-badge@latest
```

Then you need to configure OpenSSH to use this program as a `LocalCommand`; edit `~/.ssh/config` and
add something like this:

```
Host *
  PermitLocalCommand yes
  LocalCommand ssh-iterm2-badge %h
```

You just need to ensure that no other `Host` directive is setting another `LocalCommand` or
disabling `PermitLocalCommand` before this configuration block is reached.
