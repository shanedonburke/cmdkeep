package cmdkeep

import (
	"cmdkeep/cli"
	"cmdkeep/model"
	"fmt"
	"maps"
	"sort"
	"strings"
)

type ListCommand struct{}

func (lc *ListCommand) Run(cl *cli.CLI, m *model.Model) {
	keys := make([]string, len(m.Commands))
	i := 0
	for key := range maps.Keys(m.Commands) {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	if m.Last != "" {
		lastParts := strings.Split(m.Last, ":")
		fmt.Print("Last command used: ")
		if lastParts[0] == "key" {
			fmt.Printf("%s\n\n", lastParts[1])
		} else {
			fmt.Printf("`%s`\n\n", lastParts[1])
		}
	}
	if len(m.Commands) > 0 {
		fmt.Println("Commands:")
		maxKeyLen := 0
		for _, key := range keys {
			if len(key) > maxKeyLen {
				maxKeyLen = len(key)
			}
		}
		maxKeyLen++ // Colon at the end
		keyFormat := fmt.Sprintf("%%-%ds", maxKeyLen)
		for _, key := range keys {
			formattedKey := fmt.Sprintf(keyFormat, key+":")
			fmt.Printf("  %s  %s\n", formattedKey, m.Commands[key].Template)
		}
	} else {
		fmt.Println("No commands saved - try `ck add`.")
	}
}
