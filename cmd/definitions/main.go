package main

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/samber/lo"
	"sort"
	"strings"
)

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

	// Build definition
	builder := strings.Builder{}
	builder.WriteString("---@meta\n\n")
	for _, key := range cats {
		cat := docs.Categories[key]

		builder.WriteString("-- #####################################\n")
		builder.WriteString("-- " + cat.Name + "\n")
		builder.WriteString("-- #####################################\n\n")

		for _, key := range globals {
			glob := docs.Globals[key]
			if glob.Category != cat.Name {
				continue
			}
			builder.WriteString(strings.Join(lo.Map(strings.Split(glob.Description, "\n"), func(item string, index int) string {
				return "--- " + item
			}), "\n"))
			builder.WriteString("\n" + glob.Name + " = \"\"\n\n")
		}

		for _, key := range functions {
			fn := docs.Functions[key]
			if fn.Category != cat.Name {
				continue
			}
			builder.WriteString(strings.Join(lo.Map(strings.Split(fn.Description, "\n"), func(item string, index int) string {
				return "--- " + item
			}), "\n"))
			if len(fn.Args) > 0 {
				builder.WriteString("\n" + strings.Join(lo.Map(fn.Args, func(item string, index int) string {
					isOptional := strings.HasPrefix(item, "(optional) ")
					if isOptional {
						item = strings.TrimPrefix(item, "(optional) ")
					}

					t := "any"
					split := strings.Split(item, ":")
					if len(split) > 1 {
						t = strings.TrimSpace(split[1])
					}

					return "---@param " + strings.TrimSpace(split[0]) + lo.Ternary(isOptional, "?", "") + " " + t
				}), "\n"))
			}
			if fn.Return != "" {
				builder.WriteString("\n---@return " + fn.Return)
			}
			builder.WriteString("\nfunction " + fn.Name + "(")
			builder.WriteString(strings.Join(lo.Map(fn.Args, func(item string, index int) string {
				isOptional := strings.HasPrefix(item, "(optional) ")
				if isOptional {
					item = strings.TrimPrefix(item, "(optional) ")
				}
				return strings.TrimSpace(strings.Split(item, ":")[0])
			}), ", "))
			builder.WriteString(") end\n\n")
		}
	}

	// Output docs
	fmt.Print(builder.String())
}
