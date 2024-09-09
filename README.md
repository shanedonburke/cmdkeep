# CmdKeep

CmdKeep is a tool for conveniently saving and reusing shell commands.

Use it for commands you often forget, are tedious to type, or have several parameters.

Commands can be parameterized, and previous arguments are remembered.

# Installation

Check the [Releases](https://github.com/shanedonburke/cmdkeep/releases) section for the latest executables.

Add `ck` to your PATH environment variable for convenience.

# Usage Examples

Run `ck -h` to see available operations, or `ck <operation> -h` for help with an individual operations.

# Command Templates

Commands are specified as template strings, where sets of braces indicate parameters.

Braces may optionally contain a name for the parameter. Multiple parameters with the same name will share a value.

## Add a Command

This creates a `cat` command with a `file` parameter.

```shell
ck add cat "cat {file}"
```

This creates a `diff` command with two unnamed parameters.

```shell
ck add diff "diff {} {}"
```

## Running Commands

This runs a saved `cat` command. You will be prompted for any arguments.

```shell
ck cat
# OR
ck run cat
```

This runs the `diff` command with the given arguments (to skip prompting).

```shell
ck diff --args=firstarg,secondarg
```

This runs the last command executed, prompting the user for arguments.

```shell
ck last
```

This reruns the last command executed with the same arguments.
The user will not be prompted.

```shell
ck last -dy
```

This runs an arbitrary command without saving it.

```shell
ck -c 'git commit -m "{}"'
```

This prints the assembled `diff` command instead of executing it.

```shell
ck diff -p --args=firstarg,secondarg
```

## Listing Commands

This will list all saved commands.

```shell
ck commands
```

## Deleting Saved Comamnds

This deletes `cat` and `diff` commands that were previously saved using `ck add`.

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
# Print the command with the given argument and evaluate it.
eval "$(ck run source -p --args=install.sh)"
```

