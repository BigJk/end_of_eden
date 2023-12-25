package main

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/internal/lua/ludoc"
	"github.com/samber/lo"
	"regexp"
	"sort"
	"strings"
)

const globTemplate = `<details> <summary><b><code>%s</code></b> </summary> <br/>

%s

</details>

`

const fnTemplate = `<details> <summary><b><code>%s</code></b> </summary> <br/>

%s

**Signature:**

%s

</details>

`

var linkPattern = regexp.MustCompile("[^a-zA-Z\\d\\s:-]")

func functionToMd(fn ludoc.Function) string {
	signature := fmt.Sprintf("%s(%s) -> %s", fn.Name, strings.Join(fn.Args, ", "), lo.Ternary(len(fn.Return) == 0, "None", fn.Return))
	return fmt.Sprintf(fnTemplate, fn.Name, fn.Description, "```\n"+signature+"\n```")
}

func globToMd(glob ludoc.Global) string {
	return fmt.Sprintf(globTemplate, glob.Name, glob.Description)
}

func main() {
	docs := game.NewSession().LuaDocs()

	// Fetch the globals and functions keys
	globals := lo.Keys(docs.Globals)
	functions := lo.Keys(docs.Functions)

	// Sort them
	sort.Strings(globals)
	sort.Strings(functions)

	// Get the sorted categories
	cats := lo.Keys(docs.Categories)
	sort.SliceStable(cats, func(i, j int) bool {
		return docs.Categories[cats[i]].Order < docs.Categories[cats[j]].Order
	})

	// Build index
	builder := strings.Builder{}
	builder.WriteString("# End Of Eden Lua Docs\n")

	builder.WriteString("## Index\n\n")
	for _, key := range cats {
		cat := docs.Categories[key]
		builder.WriteString(fmt.Sprintf("- [%s](#%s)\n", cat.Name, linkPattern.ReplaceAllString(strings.ToLower(strings.ReplaceAll(cat.Name, " ", "-")), "")))
	}

	builder.WriteString("\n")

	// Build docs
	for _, key := range cats {
		cat := docs.Categories[key]

		builder.WriteString("## " + cat.Name + "\n\n")
		builder.WriteString(cat.Description + "\n\n")

		builder.WriteString("### Globals\n")

		hasGlobals := false
		for _, key := range globals {
			glob := docs.Globals[key]
			if glob.Category != cat.Name {
				continue
			}
			builder.WriteString(globToMd(glob))
			hasGlobals = true
		}

		if !hasGlobals {
			builder.WriteString("\nNone\n\n")
		}

		builder.WriteString("### Functions\n")

		hasFunctions := false
		for _, key := range functions {
			fn := docs.Functions[key]
			if fn.Category != cat.Name {
				continue
			}
			builder.WriteString(functionToMd(fn))
			hasFunctions = true
		}

		if !hasFunctions {
			builder.WriteString("\nNone\n\n")
		}

	}

	// Output docs
	fmt.Print(builder.String())
}
