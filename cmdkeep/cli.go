package cmdkeep

type CLI struct {
	Add struct {
		Key     string `arg help:"Key that will be used to run the command."`
		Command string `arg help:"Command template."`
	} `cmd help:"Save a new command."`
	Run struct {
		Key         string   "arg optional:\"\" help:\"Command key previously saved using `ck add`\""
		Command     string   `optional:"" short:"c" help:"Command template to run instead of saved key."`
		Args        []string `optional:"" short:"a" help:"Command arguments to skip prompting."`
		UseDefaults bool     "optional:\"\" short:\"d\" help:\"Reuse arguments from last execution of this command. Overridden by `--args`.\""
		PrintOnly   bool     `short:"p" default:"false" help:"Prints the built command instead of executing it."`
	} `cmd help:"Run a saved command or a command template."`
	Last struct {
		Args        []string `optional:"" short:"a" help:"Command arguments to skip prompting."`
		UseDefaults bool     "optional:\"\" short:\"d\" help:\"Reuse arguments from last execution of this command. Overridden by `--args`.\""
		PrintOnly   bool     `short:"p" default:"false" help:"Prints the built command instead of executing it."`
		NoPrompt    bool     `name:"yes" short:"y" default:"false" help:"Skip confirmation prompt."`
	} `cmd help:"Rerun the last command that was executed."`
	Commands struct {
	} `cmd help:"List all saved commands."`
	RM struct {
		Keys []string `arg help:"Keys of commands to delete."`
	} `cmd help:"Remove a saved command."`
}
