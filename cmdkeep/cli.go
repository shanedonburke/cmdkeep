package cmdkeep

type add struct {
	Key     string `arg help:"Key that will be used to run the command."`
	Command string `arg help:"Command template."`
}

func (a *add) Help() string {
	return "Examples:\n" +
		" $ ck add cat \"cat {file}\"\n" +
		" $ ck add diff \"diff {} {}\"\n"
}

type run struct {
	Key         string   "arg optional:\"\" help:\"Command key previously saved using `ck add`\""
	Command     string   `optional:"" short:"c" help:"Command template to run instead of saved key."`
	Args        []string `optional:"" short:"a" help:"Command arguments to skip prompting."`
	UseDefaults bool     "optional:\"\" short:\"d\" help:\"Reuse arguments from last execution of this command. Overridden by `--args`.\""
	PrintOnly   bool     `short:"p" default:"false" help:"Prints the built command instead of executing it."`
}

func (r *run) Help() string {
	return "Note: `run` can be omitted - running a command is the default behavior.\n\n" +
		"Examples:\n" +
		" $ ck run cat                (Runs a `cat` command saved using `ck add`)\n" +
		" $ ck cat                    (Also runs `cat`)\n" +
		" $ ck cat --args=myfile.txt  (Runs `cat` with an argument to skip prompting)\n" +
		" $ ck cat -d                 (Runs `cat` with its most recent arguments)\n" +
		" $ ck cat -p                 (Builds `cat` template, then prints it)\n" +
		" $ ck -c \"echo {message}\"    (Runs the specified template as if it were a saved command)\n"
}

type last struct {
	Args        []string `optional:"" short:"a" help:"Command arguments to skip prompting."`
	UseDefaults bool     "optional:\"\" short:\"d\" help:\"Reuse arguments from last execution of this command. Overridden by `--args`.\""
	PrintOnly   bool     `short:"p" default:"false" help:"Prints the built command instead of executing it."`
	NoPrompt    bool     `name:"yes" short:"y" default:"false" help:"Skip confirmation prompt."`
}

func (l *last) Help() string {
	return "Examples:\n" +
		" $ ck last     (Runs the last command executed)\n" +
		" $ ck last -d  (Runs last command with same arguments)\n" +
		" $ ck last -y  (Runs last command without confirmation)\n" +
		" $ ck last -p  (Prints last command instead of running it)\n"
}

type rm struct {
	Keys []string `arg help:"Keys of commands to delete."`
}

func (r *rm) Help() string {
	return "Examples:\n" +
		" $ rm cat        (Deletes one command)\n" +
		" $ rm echo diff  (Deletes two commands)\n"
}

type CLI struct {
	Add      add `cmd help:"Save a new command."`
	Commands struct {
	} `cmd help:"List all saved commands."`
	Run  run  `cmd default:"withargs" help:"Run a saved command or a command template."`
	Last last `cmd help:"Rerun the last command that was executed."`
	RM   rm   `cmd help:"Remove one or more saved commands."`
}
