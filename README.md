# CmdKeep

CmdKeep is a tool for conveniently storing console commands.

Use it for commands you often forget, are tedious to type, or have several parameters.

Commands can be parameterized, and previous arguments are remembered.

# Usage

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

This will run an arbitrary command without saving it.

```shell
ck run --command 'git commit -m "{message}"'
```

## List Commands

This will list all saved commands.

```shell
ck commands
```

## Delete a Saved Command

This will delete a `cat` command that was previously saved using `ck add`.

```shell
ck rm cat
```