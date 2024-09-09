# CmdKeep

CmdKeep is a tool for conveniently saving and reusing shell commands.

Use it for commands you often forget, are tedious to type, or have several parameters.

Commands can be parameterized, and previous arguments are remembered.

# Installation

Check the [Releases](https://github.com/shanedonburke/cmdkeep/releases/latest) section for the latest executables.

Add `ck` to your PATH environment variable for convenience.

# Usage

Run `ck -h` to see available operations, or `ck <operation> -h` for help with an individual operation.

## Command Templates

Commands are specified as template strings, where sets of braces indicate parameters.

Braces may optionally contain a name for the parameter. Multiple parameters with the same name will share a value.

## Adding Commands

★ Creates a `cat` command with a `file` parameter:

```shell
ck add cat "cat {file}"
```

★ Creates a `diff` command with two unnamed parameters:

```shell
ck add diff "diff {} {}"
```

## Running Commands

★ Runs a saved `cat` command. You will be prompted for any arguments:

```shell
ck cat
# OR
ck run cat
```

★ Runs the `diff` command with the given arguments (to skip prompting):

```shell
ck diff --args=firstarg,secondarg
```

★ Reruns the last command executed, prompting the user for arguments:

```shell
ck last
```

★ Reruns the last command with the same arguments.
The user will not be prompted:

```shell
ck last -dy
```

★ Runs an arbitrary command without saving it:

```shell
ck -c 'git commit -m "{}"'
```

★ Prints the assembled `diff` command instead of executing it:

```shell
ck diff -p --args=firstarg,secondarg
```

## Listing Commands

★ Lists all saved commands:

```shell
ck commands
```

## Deleting Saved Comamnds

★ Deletes `cat` and `diff` commands that were previously saved using `ck add`:

```shell
ck rm cat diff
```

## Other Use Cases

### Run Multiple Commands in One Template (Bash)

```shell
ck add commit_push "bash -c 'git add . && git commit && git push'"
```

### Evaluate a Command in the Current Shell (Bash)

```shell
# Add a `source` command with a 'file' parameter.
ck add source "source {file}"
# Print the command with the given argument, then evaluate the result.
eval "$(ck source -p --args=install.sh)"
```

# Building from Source

1. Clone the repository.
2. Run `go build -o bin/ck` (platform-dependent).

# Contributing

Pull requests are welcome!

