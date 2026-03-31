# espanso-misspell

`espanso-misspell` generates [espanso](https://espanso.org) packages that
replace commonly misspelled words in multiple languages.

## Packages

| Package | Language | Source | Entries |
|---|---|---|---|
| `misspell-de` | German | [Wikipedia Wortliste](https://de.wikipedia.org/wiki/Wikipedia:Helferlein/Rechtschreibpr%C3%BCfung/Wortliste) | ~5.7k |
| `misspell-en` | English | [client9/misspell](https://github.com/client9/misspell) DictMain | ~28k |
| `misspell-en_UK` | English | [client9/misspell](https://github.com/client9/misspell) DictBritish | ~1.5k |
| `misspell-en_US` | English | [client9/misspell](https://github.com/client9/misspell) DictAmerican | ~1.6k |
| `misspell-es` | Spanish | [Wikipedia Corrector ortográfico](https://es.wikipedia.org/wiki/Wikipedia:Corrector_ortogr%C3%A1fico/Listado) | ~19k |
| `misspell-fr` | French | [Wikipédia Fautes d'orthographe](https://fr.wikipedia.org/wiki/Wikip%C3%A9dia:Liste_de_fautes_d%27orthographe_courantes) | ~91 |
| `misspell-it` | Italian | [Wikipedia Errori comuni](https://it.wikipedia.org/wiki/Wikipedia:Bot/Richieste/Errori_comuni) | ~331 |

## Usage

Generate all packages:

```bash
go run ./cmd/generate
```

Generate packages for a specific language:

```bash
go run ./cmd/generate de
go run ./cmd/generate en es
```

Output is written to `packages/`.

## License

[BSD 3-Clause "New" or "Revised" License](LICENSE)

English word lists from misspell are
[MIT](https://github.com/client9/misspell/blob/master/LICENSE).
German, Spanish, French, and Italian word lists from Wikipedia are
[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
