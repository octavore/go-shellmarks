# go-shellmarks

A quick way to jump around in your shell

## Setup

Required. Add this to your .bashrc to use go-shellmarks as `z`:

```shell
source <(go-shellmarks @shell --alias z)
```

This change will take effect when you relaunch your terminal.

## Usage

```
z <bookmark>   jumps to <bookmark>
z @add         bookmark the current directory as <bookmark>
z @get         show stored value for <bookmark>
z @ls          list all bookmarks
z @rm          remove a stored bookmark
```

## License

MIT
