// Command generate produces espanso misspell packages for multiple languages.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/client9/misspell"
	"github.com/timorunge/espanso"

	"github.com/timorunge/espanso-misspell/internal/fetch"
)

const (
	author      = "Timo Runge"
	authorEmail = author + " <me@timorunge.com>"
	repo        = "https://github.com/timorunge/espanso-misspell"
	year        = "2019-2026"

	packagesDir = "packages"

	wikiDeURL    = "https://de.wikipedia.org/w/index.php?title=Wikipedia:Helferlein/Rechtschreibpr%C3%BCfung/Wortliste&action=raw"
	wikiEsURL    = "https://es.wikipedia.org/w/index.php?title=Wikipedia:Corrector_ortogr%C3%A1fico/Listado&action=raw"
	wikiFrURL    = "https://fr.wikipedia.org/w/index.php?title=Wikip%C3%A9dia:Liste_de_fautes_d%27orthographe_courantes&action=raw"
	wikiItURL    = "https://it.wikipedia.org/w/index.php?title=Wikipedia:Bot/Richieste/Errori_comuni&action=raw"
	codespellURL = "https://raw.githubusercontent.com/codespell-project/codespell/master/codespell_lib/data/dictionary.txt"
)

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

[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/)

Word list data from Wikipedia is licensed under
[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).`

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

[MIT License](LICENSE)

Misspell is [MIT](https://github.com/client9/misspell/blob/master/LICENSE).`

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

[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/)

Word list data from Wikipedia is licensed under
[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).`

const longDescFR = `# misspell-fr

misspell-fr is an espanso package which is replacing commonly misspelled
french words. The package is auto-generated from the
[Wikipédia Liste de fautes d'orthographe courantes](https://fr.wikipedia.org/wiki/Wikip%C3%A9dia:Liste_de_fautes_d%27orthographe_courantes).

## Installation

` + "```" + `
espanso install misspell-fr
espanso restart
` + "```" + `

## Usage

Type ` + "`aigü`" + ` and see it replaced with ` + "`aigu`" + `.

## License

[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/)

Word list data from Wikipedia is licensed under
[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).`

const longDescIT = `# misspell-it

misspell-it is an espanso package which is replacing commonly misspelled
italian words. The package is auto-generated from the
[Wikipedia Errori comuni](https://it.wikipedia.org/wiki/Wikipedia:Bot/Richieste/Errori_comuni).

## Installation

` + "```" + `
espanso install misspell-it
espanso restart
` + "```" + `

## Usage

Type ` + "`conbattere`" + ` and see it replaced with ` + "`combattere`" + `.

## License

[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/)

Word list data from Wikipedia is licensed under
[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).`

const longDescCodespell = `# misspell-en-codespell

misspell-en-codespell is an espanso package which is replacing commonly
misspelled english words. The package is auto-generated from the
[codespell](https://github.com/codespell-project/codespell) dictionary,
deduplicated against misspell-en to avoid trigger conflicts.

## Installation

` + "```" + `
espanso install misspell-en-codespell
espanso restart
` + "```" + `

## Usage

Type ` + "`acessibility`" + ` and see it replaced with ` + "`accessibility`" + `.

## License

[CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/)

Dictionary data from codespell is derived from English Wikipedia and licensed
under [CC BY-SA 3.0](https://creativecommons.org/licenses/by-sa/3.0/).`

type pkg struct {
	name      string
	version   string
	title     string
	shortDesc string
	longDesc  string
	license   espanso.License
	fetch     func(context.Context) (espanso.Matches, error)
}

func dictFetcher(dict []string) func(context.Context) (espanso.Matches, error) {
	return func(_ context.Context) (espanso.Matches, error) {
		return espanso.DictToMatches(dict)
	}
}

func urlFetcher(url, lang string, parse fetch.LineParser) func(context.Context) (espanso.Matches, error) {
	return func(ctx context.Context) (espanso.Matches, error) {
		return fetch.Fetch(ctx, url, lang, parse)
	}
}

func main() {
	mitLicense := espanso.MIT(year, authorEmail)
	ccBySa40 := espanso.CCBYSA40(year, authorEmail)
	ccBySa30 := espanso.CCBYSA30(year, authorEmail)

	packages := map[string][]pkg{
		"en": {
			{
				name: "misspell-en", version: "0.1.2",
				title: "Misspell EN", shortDesc: "Replace commonly misspelled english words.",
				longDesc: fmt.Sprintf(longDescEN, "misspell-en", "commonly misspelled english words", "yuo"),
				license:  mitLicense,
				fetch:    dictFetcher(misspell.DictMain),
			},
			{
				name: "misspell-en_UK", version: "0.1.2",
				title: "Misspell en_UK", shortDesc: "Replace american english with british english.",
				longDesc: fmt.Sprintf(longDescEN, "misspell-en_UK", "american english with british english", "color"),
				license:  mitLicense,
				fetch:    dictFetcher(misspell.DictBritish),
			},
			{
				name: "misspell-en_US", version: "0.1.2",
				title: "Misspell en_US", shortDesc: "Replace british english with american english.",
				longDesc: fmt.Sprintf(longDescEN, "misspell-en_US", "british english with american english", "tyre"),
				license:  mitLicense,
				fetch:    dictFetcher(misspell.DictAmerican),
			},
		},
		"codespell": {
			{
				name: "misspell-en-codespell", version: "0.1.0",
				title: "Misspell EN Codespell", shortDesc: "Replace commonly misspelled english words (codespell).",
				longDesc: longDescCodespell,
				license:  ccBySa30,
				fetch:    urlFetcher(codespellURL, "codespell", fetch.CodespellLine),
			},
		},
		"de": {
			{
				name: "misspell-de", version: "0.1.0",
				title: "Misspell DE", shortDesc: "Replace commonly misspelled german words.",
				longDesc: longDescDE,
				license:  ccBySa40,
				fetch:    urlFetcher(wikiDeURL, "de", fetch.PipeLine),
			},
		},
		"es": {
			{
				name: "misspell-es", version: "0.1.0",
				title: "Misspell ES", shortDesc: "Replace commonly misspelled spanish words.",
				longDesc: longDescES,
				license:  ccBySa40,
				fetch:    urlFetcher(wikiEsURL, "es", fetch.PipeLine),
			},
		},
		"fr": {
			{
				name: "misspell-fr", version: "0.1.0",
				title: "Misspell FR", shortDesc: "Replace commonly misspelled french words.",
				longDesc: longDescFR,
				license:  ccBySa40,
				fetch:    urlFetcher(wikiFrURL, "fr", fetch.TemplateLine("Suggestion")),
			},
		},
		"it": {
			{
				name: "misspell-it", version: "0.1.0",
				title: "Misspell IT", shortDesc: "Replace commonly misspelled italian words.",
				longDesc: longDescIT,
				license:  ccBySa40,
				fetch:    urlFetcher(wikiItURL, "it", fetch.TemplateLine("BR1")),
			},
		},
	}

	targets := os.Args[1:]
	if len(targets) == 0 || slices.Contains(targets, "all") {
		targets = []string{"en", "codespell", "de", "es", "fr", "it"}
	}

	// Collect misspell-en triggers for codespell deduplication.
	var enTriggers map[string]struct{}
	if slices.Contains(targets, "codespell") {
		enTriggers = make(map[string]struct{})
		for _, dict := range [][]string{misspell.DictMain, misspell.DictBritish, misspell.DictAmerican} {
			for i := 0; i < len(dict); i += 2 {
				enTriggers[dict[i]] = struct{}{}
			}
		}
	}

	ctx := context.Background()
	for _, lang := range targets {
		pkgs, ok := packages[lang]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown language: %s (available: en, codespell, de, es, fr, it)\n", lang)
			os.Exit(1)
		}
		for _, p := range pkgs {
			var exclude map[string]struct{}
			if lang == "codespell" {
				exclude = enTriggers
			}
			if err := generate(ctx, p, exclude); err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", p.name, err)
				os.Exit(1)
			}
		}
	}
}

func generate(ctx context.Context, p pkg, exclude map[string]struct{}) error {
	matches, err := p.fetch(ctx)
	if err != nil {
		return fmt.Errorf("fetch: %w", err)
	}

	if len(exclude) > 0 {
		matches = matches.Filter(func(m espanso.Match) bool {
			for _, t := range m.Triggers {
				if _, ok := exclude[t]; ok {
					return false
				}
			}
			return true
		})
	}

	matches = matches.SetWord(true).SetPropagateCase(true).Sort().Deduplicate()
	fmt.Printf("%s: %d matches\n", p.name, len(matches))

	metaDir := filepath.Join(packagesDir, p.name)
	pkgDir := filepath.Join(metaDir, p.version)

	manifest := espanso.Manifest{
		Name:        p.name,
		Title:       p.title,
		Description: p.shortDesc,
		Version:     p.version,
		Author:      author,
	}
	if err := manifest.WriteFile(metaDir); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}

	ep := espanso.Package{
		Name:    p.name,
		Parent:  "default",
		Version: p.version,
		Matches: matches,
	}
	if err := ep.WriteFile(pkgDir); err != nil {
		return fmt.Errorf("write package: %w", err)
	}

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
		return fmt.Errorf("write readme: %w", err)
	}

	if err := p.license.WriteFile(metaDir); err != nil {
		return fmt.Errorf("write license: %w", err)
	}

	return nil
}
