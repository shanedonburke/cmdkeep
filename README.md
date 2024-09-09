# CmdKeep

CmdKeep is a tool for conveniently storing console commands.

Use it for commands you often forget, are tedious to type, or have several parameters.

Commands can be parameterized, and previous arguments are remembered.

# Installation

Check the [Releases](https://github.com/shanedonburke/cmdkeep/releases) section for the latest executables.

Add `ck` to your PATH environment variable for convenience.

# Usage

Run `ck --help` to see available operations, or `ck <command> --help` for help with an individual command.

# Command Templates

Commands are specified as template strings, where sets of braces indicate parameters.

Braces may optionally contain a name for the parameter.

## Add a Command

This creates a `cat` command with a `file` parameter.

```shell
ck add cat "cat {file}"
```

## Run a Command

This will run the `cat` command. You will be prompted for any arguments.

```shell
ck run cat
```


This will rerun the last command executed with the given arguments.

```shell
ck last --args=myfile.txt,someotherarg
```

This will rerun the last command executed with that execution's arguments.
The user will not be prompted.

```shell
ck last -dy
```

This will run an arbitrary command without saving it.

```shell
ck run --command 'git commit -m "{}"'
```

## List Commands

This will list all saved commands.

```shell
ck commands
```

## Delete Saved Comamnds

This will delete `cat` and `echo` commands that were previously saved using `ck add`.

```shell
ck rm cat echo
```

## Other Use Cases

### Evaluate a Command in the Current Shell (Bash)

```shell
# Add a `source` command with one parameter.
ck add source "source {file}"
# Print the command with the given argument, and evaluate it.
eval "$(ck run source -p --args=install.sh)"
```
