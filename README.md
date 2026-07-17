# efind

[中文](README.zh-CN.md)

A small command-line tool written in Go to locate executables in your `PATH` — an enhanced alternative to `which` / `where`.

`efind` resolves the absolute path of a command, optionally evaluates symlinks, and understands tools managed by [Scoop](https://scoop.sh/) and [mise](https://mise.jdx.dev/).

## Features

- Locate one or more executables by name (like `which`).
- Optionally resolve symlinks / shim targets to their real file (`--eval`).
- Seamless support for executables installed via **Scoop** (`scoop which`) and **mise** (`mise which`).
- Multiple output formats: aligned human-readable table, plain path, or structured **JSON**.
- Multiple names can be queried in a single invocation.

## Installation

Requires Go 1.26+.

```bash
go install github.com/cnk3x/efind@latest
```

Or build from source:

```bash
git clone https://github.com/cnk3x/efind
cd efind
go build -o efind .
```

## Usage

```
efind [flags] <name> [<name> ...]
```

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--eval` | `-e` | Evaluate symlinks / shim targets, printing the real path. |
| `--json` | `-j` | Output results as JSON. |
| `--noe` |      | Suppress error output for executables that are not found. |
| `--verbose` | `-v` | Force verbose (aligned) output even for a single name. |

If no name is provided, usage information is printed.

### Examples

Find a single command:

```bash
$ efind go
C:\Users\wen\sdk\go\bin\go.exe
```

Resolve the real path behind a shim / symlink:

```bash
$ efind -e node
node => C:\Users\wen\bin\mise\bin\node => C:\Users\wen\.local\share\mise\installs\node\21.0.0\bin\node.exe
```

Query multiple commands at once (verbose aligned output is enabled automatically):

```bash
$ efind git python rustc
git    => C:\Program Files\Git\cmd\git.exe
python => C:\Users\wen\AppData\Local\Programs\Python\python.exe
rustc  => executable file not found
```

Suppress not-found errors:

```bash
$ efind --noe rustc
```

Get machine-readable output:

```bash
$ efind -j git
[
  {
    "name": "git",
    "full": "C:\\Program Files\\Git\\cmd\\git.exe"
  }
]
```

## How it works

1. `efind` uses `exec.LookPath` to find the executable in `PATH`.
2. When `--eval` is set, it tries to resolve the real location in this order:
   - **Scoop** shim via `scoop which <name>`
   - **mise** shim via `mise which <name>`
   - Otherwise falls back to `filepath.EvalSymlinks`.

## License

[MIT](LICENSE)
