package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/client9/misspell"
	"github.com/timorunge/espanso"
	"github.com/timorunge/espanso-misspell/internal/fetch"
)

const (
	author = "Timo Runge"
	repo   = "https://github.com/timorunge/espanso-misspell"
	year   = "2019-2026"

	packagesDir = "packages"

	wortlisteURL = "https://de.wikipedia.org/w/index.php?title=Wikipedia:Helferlein/Rechtschreibpr%C3%BCfung/Wortliste&action=raw"
	listadoURL   = "https://es.wikipedia.org/w/index.php?title=Wikipedia:Corrector_ortogr%C3%A1fico/Listado&action=raw"
)

const longDescEN = `# %[1]s

%[1]s is an espanso package which is replacing %[2]s.
The package is based on [github.com/client9/misspell](https://github.com/client9/misspell).

## Installation

` + "```" + `
espanso install %[1]s
espanso restart
` + "```" + `

## Usage

Type ` + "`%[3]s`" + ` and see what's happening.

## License

[BSD 3-Clause "New" or "Revised" License](LICENSE)

Misspell is [MIT](https://github.com/client9/misspell/blob/master/LICENSE).`

const longDescDE = `# misspell-de

misspell-de is an espanso package which is replacing commonly misspelled
german words. The package is auto-generated from the
[Wikipedia Rechtschreibprüfung Wortliste](https://de.wikipedia.org/wiki/Wikipedia:Helferlein/Rechtschreibpr%C3%BCfung/Wortliste).

## Installation

` + "```" + `
espanso install misspell-de
espanso restart
` + "```" + `

## Usage

Type ` + "`Aaachen`" + ` and see it replaced with ` + "`Aachen`" + `.

## License

[BSD 3-Clause "New" or "Revised" License](LICENSE)

Word list data from Wikipedia is licensed under
[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).`

const longDescES = `# misspell-es

misspell-es is an espanso package which is replacing commonly misspelled
spanish words. The package is auto-generated from the
[Wikipedia Corrector ortográfico](https://es.wikipedia.org/wiki/Wikipedia:Corrector_ortogr%C3%A1fico/Listado).

## Installation

` + "```" + `
espanso install misspell-es
espanso restart
` + "```" + `

## Usage

Type ` + "`accomodar`" + ` and see it replaced with ` + "`acomodar`" + `.

## License

[BSD 3-Clause "New" or "Revised" License](LICENSE)

Word list data from Wikipedia is licensed under
[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).`

type pkg struct {
	name      string
	version   string
	title     string
	shortDesc string
	longDesc  string
	fetch     func() (espanso.Matches, error)
}

func dictFetcher(dict []string) func() (espanso.Matches, error) {
	return func() (espanso.Matches, error) {
		return espanso.DictToMatches(dict), nil
	}
}

func wikiFetcher(url, lang string) func() (espanso.Matches, error) {
	return func() (espanso.Matches, error) {
		return fetch.Wikipedia(url, lang)
	}
}

func main() {
	packages := map[string][]pkg{
		"en": {
			{
				name: "misspell-en", version: "0.1.2",
				title: "Misspell EN", shortDesc: "Replace commonly misspelled english words.",
				longDesc: fmt.Sprintf(longDescEN, "misspell-en", "commonly misspelled english words", "yuo"),
				fetch:    dictFetcher(misspell.DictMain),
			},
			{
				name: "misspell-en_UK", version: "0.1.2",
				title: "Misspell en_UK", shortDesc: "Replace american english with british english.",
				longDesc: fmt.Sprintf(longDescEN, "misspell-en_UK", "american english with british english", "color"),
				fetch:    dictFetcher(misspell.DictBritish),
			},
			{
				name: "misspell-en_US", version: "0.1.2",
				title: "Misspell en_US", shortDesc: "Replace british english with american english.",
				longDesc: fmt.Sprintf(longDescEN, "misspell-en_US", "british english with american english", "tyre"),
				fetch:    dictFetcher(misspell.DictAmerican),
			},
		},
		"de": {
			{
				name: "misspell-de", version: "0.1.0",
				title: "Misspell DE", shortDesc: "Replace commonly misspelled german words.",
				longDesc: longDescDE,
				fetch:    wikiFetcher(wortlisteURL, "de"),
			},
		},
		"es": {
			{
				name: "misspell-es", version: "0.1.0",
				title: "Misspell ES", shortDesc: "Replace commonly misspelled spanish words.",
				longDesc: longDescES,
				fetch:    wikiFetcher(listadoURL, "es"),
			},
		},
	}

	targets := os.Args[1:]
	if len(targets) == 0 || slices.Contains(targets, "all") {
		targets = []string{"en", "de", "es"}
	}

	for _, lang := range targets {
		pkgs, ok := packages[lang]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown language: %s (available: en, de, es)\n", lang)
			os.Exit(1)
		}
		for _, p := range pkgs {
			generate(p)
		}
	}
}

func generate(p pkg) {
	matches, err := p.fetch()
	if err != nil {
		panic(fmt.Sprintf("fetch %s: %v", p.name, err))
	}
	matches = matches.SetWord(true).SetPropagateCase(true).Sort().Deduplicate()
	fmt.Printf("%s: %d matches\n", p.name, len(matches))

	dir := filepath.Join(packagesDir, p.name, p.version)

	ep := espanso.Package{
		Name:    p.name,
		Parent:  "default",
		Version: p.version,
		Matches: matches,
	}
	if err := ep.WriteFile(dir); err != nil {
		panic(err)
	}

	metaDir := filepath.Join(packagesDir, p.name)

	r := espanso.Readme{
		Name:      p.name,
		Title:     p.title,
		ShortDesc: p.shortDesc,
		Version:   p.version,
		Author:    author,
		Repo:      repo,
		LongDesc:  p.longDesc,
	}
	if err := r.WriteFile(metaDir); err != nil {
		panic(err)
	}

	l := espanso.BSD3Clause(year, fmt.Sprintf("%s <me@timorunge.com>", author))
	if err := l.WriteFile(metaDir); err != nil {
		panic(err)
	}
}
