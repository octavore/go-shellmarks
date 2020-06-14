# go-shellmarks

A quick way to jump around in your shell

## Installation

Add this to your `~/.bashrc` file:

```bash
source <(go-shellmarks -shell)
```

## Usage

`go-shellmarks -shell` needs to be sourced when your shell initializes (see
Installation above). This registers a `g` command in your shell, with the
following flags:

```
g <bookmark>       jumps to <bookmark>
g -a <bookmark>    add a new bookmark to the current folder
g -d <bookmark>    delete a bookmark
g -l               list bookmarks
g -h               show help
```

## License

MIT
