package main

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/acarl005/stripansi"
	"github.com/alexeyco/simpletable"
	"github.com/samber/lo"
	"regexp"
	"sort"
	"strings"
)

func stripSpecial(s string) string {
	// 1. collapse consecutive newlines
	newlineRegex := regexp.MustCompile(`\n+`)
	s = newlineRegex.ReplaceAllString(s, " - ")

	// 2. replace [ with **
	return strings.ReplaceAll(strings.ReplaceAll(s, "[", "**"), "]", "**")
}

func makeCode(s string) string {
	return "``" + s + "``"
}

func makeCodes(s []string) string {
	return strings.Join(lo.Map(s, func(s string, i int) string {
		return makeCode(s)
	}), ", ")
}

func makeCodeBlock(s string) string {
	return "<code style='white-space: pre;'>" + strings.ReplaceAll(strings.ReplaceAll(s, "\n", "</br>"), "\\", "\\\\") + " </code>"
}

func buildGameContentDocs() {
	res := game.NewSession().GetResources()

	fmt.Println("# Game Content")

	fmt.Println(`This document contains all game content that is available in the game. It is automatically generated and may be out of date.
Content that is dynamically generated at runtime is not included in this document, only content that is registered at the beginning of a session.

`)

	// Stats

	fmt.Println("\n\n## Stats\n")

	statsTable := simpletable.New()
	statsTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "Type"},
			{Align: simpletable.AlignCenter, Text: "Count"},
		},
	}
	statsTable.SetStyle(simpletable.StyleMarkdown)

	countArtifacts := len(res.Artifacts)
	countCards := len(res.Cards)
	countStatusEffects := len(res.StatusEffects)
	countEnemies := len(res.Enemies)
	countEvents := len(res.Events)

	statsTable.Body.Cells = append(statsTable.Body.Cells, []*simpletable.Cell{
		{Text: "Artifacts"},
		{Text: fmt.Sprint(countArtifacts)},
	})

	statsTable.Body.Cells = append(statsTable.Body.Cells, []*simpletable.Cell{
		{Text: "Cards"},
		{Text: fmt.Sprint(countCards)},
	})

	statsTable.Body.Cells = append(statsTable.Body.Cells, []*simpletable.Cell{
		{Text: "Status Effects"},
		{Text: fmt.Sprint(countStatusEffects)},
	})

	statsTable.Body.Cells = append(statsTable.Body.Cells, []*simpletable.Cell{
		{Text: "Enemies"},
		{Text: fmt.Sprint(countEnemies)},
	})

	statsTable.Body.Cells = append(statsTable.Body.Cells, []*simpletable.Cell{
		{Text: "Events"},
		{Text: fmt.Sprint(countEvents)},
	})

	statsTable.Println()

	// Artifacts

	artifactTable := simpletable.New()
	artifactTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Description"},
			{Align: simpletable.AlignCenter, Text: "Price"},
			{Align: simpletable.AlignCenter, Text: "Tags"},
			{Align: simpletable.AlignCenter, Text: "Test Present"},
		},
	}
	artifactTable.SetStyle(simpletable.StyleMarkdown)

	artifacts := lo.Values(res.Artifacts)
	sort.SliceStable(artifacts, func(i, j int) bool {
		return artifacts[i].ID < artifacts[j].ID
	})
	sort.SliceStable(artifacts, func(i, j int) bool {
		return artifacts[i].Price < artifacts[j].Price
	})

	for _, v := range artifacts {
		r := []*simpletable.Cell{
			{Text: "``" + v.ID + "``"},
			{Text: v.Name},
			{Text: stripSpecial(stripansi.Strip(v.Description))},
			{Text: fmt.Sprint(v.Price)},
			{Text: strings.Join(v.Tags, ", ")},
			{Text: lo.Ternary(v.Test != nil, ":heavy_check_mark:", ":no_entry_sign:")},
		}

		artifactTable.Body.Cells = append(artifactTable.Body.Cells, r)
	}

	fmt.Println("\n\n## Artifacts\n")
	artifactTable.Println()

	// Cards

	cardTable := simpletable.New()
	cardTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Description"},
			{Align: simpletable.AlignCenter, Text: "Action Points"},
			{Align: simpletable.AlignCenter, Text: "Exhaust"},
			{Align: simpletable.AlignCenter, Text: "Consumable"},
			{Align: simpletable.AlignCenter, Text: "Max Level"},
			{Align: simpletable.AlignCenter, Text: "Price"},
			{Align: simpletable.AlignCenter, Text: "Tags"},
			{Align: simpletable.AlignCenter, Text: "Color"},
			{Align: simpletable.AlignCenter, Text: "Used Callbacks"},
			{Align: simpletable.AlignCenter, Text: "Test Present"},
		},
	}
	cardTable.SetStyle(simpletable.StyleMarkdown)

	actionPoints := make(map[int]int)
	cardType := make(map[string]int)

	cards := lo.Values(res.Cards)
	sort.SliceStable(cards, func(i, j int) bool {
		return cards[i].ID < cards[j].ID
	})
	sort.SliceStable(cards, func(i, j int) bool {
		return cards[i].Price < cards[j].Price
	})

	for _, v := range cards {
		actionPoints[v.PointCost]++

		if v.DoesExhaust {
			cardType["Exhaust"]++
		} else if v.DoesConsume {
			cardType["Consume"]++
		} else {
			cardType["Normal"]++
		}

		r := []*simpletable.Cell{
			{Text: "``" + v.ID + "``"},
			{Text: v.Name},
			{Text: stripSpecial(stripansi.Strip(v.Description))},
			{Text: fmt.Sprint(v.PointCost)},
			{Text: lo.Ternary(v.DoesExhaust, ":heavy_check_mark:", ":no_entry_sign:")},
			{Text: lo.Ternary(v.DoesConsume, ":heavy_check_mark:", ":no_entry_sign:")},
			{Text: fmt.Sprint(v.MaxLevel)},
			{Text: fmt.Sprint(v.Price)},
			{Text: strings.Join(v.Tags, ", ")},
			{Text: v.Color},
			{Text: makeCodes(lo.Keys(v.Callbacks))},
			{Text: lo.Ternary(v.Test != nil, ":heavy_check_mark:", ":no_entry_sign:")},
		}

		cardTable.Body.Cells = append(cardTable.Body.Cells, r)
	}

	fmt.Println("\n\n## Cards\n")
	cardTable.Println()

	fmt.Println("\n\n### Action Points\n")

	fmt.Println(fmt.Sprintf("```mermaid\npie\ntitle Action Points\n%s\n```\n", strings.Join(lo.Map(lo.Entries(actionPoints), func(e lo.Entry[int, int], i int) string {
		return fmt.Sprintf("\"%d AP\" : %d", e.Key, e.Value)
	}), "\n")))

	fmt.Println("\n\n### Card Types\n")

	fmt.Println(fmt.Sprintf("```mermaid\npie\ntitle Card Types\n%s\n```\n", strings.Join(lo.Map(lo.Entries(cardType), func(e lo.Entry[string, int], i int) string {
		return fmt.Sprintf("\"%s\" : %d", e.Key, e.Value)
	}), "\n")))

	// Status Effects

	statusEffectTable := simpletable.New()
	statusEffectTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Description"},
			{Align: simpletable.AlignCenter, Text: "Look"},
			{Align: simpletable.AlignCenter, Text: "Foreground"},
			{Align: simpletable.AlignCenter, Text: "Can Stack"},
			{Align: simpletable.AlignCenter, Text: "Decay"},
			{Align: simpletable.AlignCenter, Text: "Rounds"},
			{Align: simpletable.AlignCenter, Text: "Used Callbacks"},
			{Align: simpletable.AlignCenter, Text: "Test Present"},
		},
	}
	statusEffectTable.SetStyle(simpletable.StyleMarkdown)

	statusEffects := lo.Values(res.StatusEffects)
	sort.SliceStable(statusEffects, func(i, j int) bool {
		return statusEffects[i].ID < statusEffects[j].ID
	})
	sort.SliceStable(statusEffects, func(i, j int) bool {
		return statusEffects[i].Order < statusEffects[j].Order
	})

	for _, v := range statusEffects {
		r := []*simpletable.Cell{
			{Text: "``" + v.ID + "``"},
			{Text: v.Name},
			{Text: stripSpecial(stripansi.Strip(v.Description))},
			{Text: v.Look},
			{Text: v.Foreground},
			{Text: lo.Ternary(v.CanStack, ":heavy_check_mark:", ":no_entry_sign:")},
			{Text: string(v.Decay)},
			{Text: fmt.Sprint(v.Rounds)},
			{Text: makeCodes(lo.Keys(v.Callbacks))},
			{Text: lo.Ternary(v.Test != nil, ":heavy_check_mark:", ":no_entry_sign:")},
		}

		statusEffectTable.Body.Cells = append(statusEffectTable.Body.Cells, r)
	}

	fmt.Println("\n\n## Status Effects\n")
	statusEffectTable.Println()

	// Enemies

	enemyTable := simpletable.New()
	enemyTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Description"},
			{Align: simpletable.AlignCenter, Text: "Initial HP"},
			{Align: simpletable.AlignCenter, Text: "Max HP"},
			{Align: simpletable.AlignCenter, Text: "Look"},
			{Align: simpletable.AlignCenter, Text: "Color"},
			{Align: simpletable.AlignCenter, Text: "Used Callbacks"},
			{Align: simpletable.AlignCenter, Text: "Test Present"},
		},
	}
	enemyTable.SetStyle(simpletable.StyleMarkdown)

	enemies := lo.Values(res.Enemies)
	sort.SliceStable(enemies, func(i, j int) bool {
		return enemies[i].ID < enemies[j].ID
	})
	sort.SliceStable(enemies, func(i, j int) bool {
		return enemies[i].Name < enemies[j].Name
	})

	for _, v := range enemies {
		r := []*simpletable.Cell{
			{Text: "``" + v.ID + "``"},
			{Text: v.Name},
			{Text: stripSpecial(stripansi.Strip(v.Description))},
			{Text: fmt.Sprint(v.InitialHP)},
			{Text: fmt.Sprint(v.MaxHP)},
			{Text: makeCodeBlock(v.Look)},
			{Text: v.Color},
			{Text: makeCodes(lo.Keys(v.Callbacks))},
			{Text: lo.Ternary(v.Test != nil, ":heavy_check_mark:", ":no_entry_sign:")},
		}

		enemyTable.Body.Cells = append(enemyTable.Body.Cells, r)
	}

	fmt.Println("\n\n## Enemies\n")
	enemyTable.Println()

	// Events

	eventTable := simpletable.New()
	eventTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "ID"},
			{Align: simpletable.AlignCenter, Text: "Name"},
			{Align: simpletable.AlignCenter, Text: "Description"},
			{Align: simpletable.AlignCenter, Text: "Tags"},
			{Align: simpletable.AlignCenter, Text: "Choices"},
			{Align: simpletable.AlignCenter, Text: "Test Present"},
		},
	}
	eventTable.SetStyle(simpletable.StyleMarkdown)

	events := lo.Values(res.Events)
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].ID < events[j].ID
	})
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].Name < events[j].Name
	})

	for _, v := range events {
		r := []*simpletable.Cell{
			{Text: "``" + v.ID + "``"},
			{Text: v.Name},
			{Text: stripSpecial(stripansi.Strip(v.Description))},
			{Text: strings.Join(v.Tags, ", ")},
			{Text: "<ul>" + strings.Join(lo.Map(v.Choices, func(c game.EventChoice, i int) string {
				return "<li>" + makeCode(stripansi.Strip(c.Description)) + "</li>"
			}), " ") + "</ul>"},
			{Text: lo.Ternary(v.Test != nil, ":heavy_check_mark:", ":no_entry_sign:")},
		}

		eventTable.Body.Cells = append(eventTable.Body.Cells, r)
	}

	fmt.Println("\n\n## Events\n")
	eventTable.Println()
}
