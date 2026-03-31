package main

import (
	"fmt"
	"path/filepath"

	"github.com/client9/misspell"
	"github.com/timorunge/espanso"
)

const (
	author  = "Timo Runge"
	repo    = "https://github.com/timorunge/espanso-misspell-en"
	version = "0.1.2"
	year    = "2019-2026"
)

const longDescFmt = `# %[1]s

%[1]s is a espanso package which is replacing %[2]s.
The package is based on [github.com/client9/misspell](https://github.com/client9/misspell).

## Installation

Install the package with:

` + "```" + `
espanso install %[1]s
espanso restart
` + "```" + `

## Usage

Type ` + "`%[3]s`" + ` and see what's happening.

## License

[BSD 3-Clause "New" or "Revised" License](LICENSE)

Misspell is [MIT](https://github.com/client9/misspell/blob/master/LICENSE).`

type pkg struct {
	name      string
	dict      []string
	title     string
	shortDesc string
	teaser    string
	example   string
}

func main() {
	packages := []pkg{
		{"misspell-en", misspell.DictMain, "Misspell EN",
			"Replace commonly misspelled english words.",
			"commonly misspelled english words", "yuo"},
		{"misspell-en_UK", misspell.DictBritish, "Misspell en_UK",
			"Replace american english with british english.",
			"american english with british english", "color"},
		{"misspell-en_US", misspell.DictAmerican, "Misspell en_US",
			"Replace british english with american english.",
			"british english with american english", "tyre"},
	}
	for _, p := range packages {
		generate(p)
	}
}

func generate(p pkg) {
	dir := filepath.Join(p.name, version)

	ep := espanso.Package{
		Name:    p.name,
		Parent:  "default",
		Version: version,
		Matches: espanso.DictToMatches(p.dict).SetWord(true).SetPropagateCase(true),
	}
	if err := ep.WriteFile(dir); err != nil {
		panic(err)
	}

	r := espanso.Readme{
		Name:      p.name,
		Title:     p.title,
		ShortDesc: p.shortDesc,
		Version:   version,
		Author:    author,
		Repo:      repo,
		LongDesc:  fmt.Sprintf(longDescFmt, p.name, p.teaser, p.example),
	}
	if err := r.WriteFile(p.name); err != nil {
		panic(err)
	}

	l := espanso.BSD3Clause(year, fmt.Sprintf("%s <me@timorunge.com>", author))
	if err := l.WriteFile(p.name); err != nil {
		panic(err)
	}
}
