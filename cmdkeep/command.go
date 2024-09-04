package cmdkeep

type CKCommand interface {
	Run(cli *CLI)
}
