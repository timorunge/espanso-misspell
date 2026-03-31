# espanso-misspell

`espanso-misspell` generates [espanso](https://espanso.org) packages that
replace commonly misspelled words in multiple languages.

## Packages

| Package | Language | Source | License | Entries |
|---|---|---|---|---|
| `misspell-de` | German | [Wikipedia Wortliste](https://de.wikipedia.org/wiki/Wikipedia:Helferlein/Rechtschreibpr%C3%BCfung/Wortliste) | CC BY-SA 4.0 | ~5.7k |
| `misspell-en` | English | [client9/misspell](https://github.com/client9/misspell) DictMain | MIT | ~28k |
| `misspell-en-codespell` | English | [codespell](https://github.com/codespell-project/codespell) | CC BY-SA 3.0 | ~50k |
| `misspell-en_UK` | English | [client9/misspell](https://github.com/client9/misspell) DictBritish | MIT | ~1.5k |
| `misspell-en_US` | English | [client9/misspell](https://github.com/client9/misspell) DictAmerican | MIT | ~1.6k |
| `misspell-es` | Spanish | [Wikipedia Corrector ortográfico](https://es.wikipedia.org/wiki/Wikipedia:Corrector_ortogr%C3%A1fico/Listado) | CC BY-SA 4.0 | ~19k |
| `misspell-fr` | French | [Wikipédia Fautes d'orthographe](https://fr.wikipedia.org/wiki/Wikip%C3%A9dia:Liste_de_fautes_d%27orthographe_courantes) | CC BY-SA 4.0 | ~91 |
| `misspell-it` | Italian | [Wikipedia Errori comuni](https://it.wikipedia.org/wiki/Wikipedia:Bot/Richieste/Errori_comuni) | CC BY-SA 4.0 | ~331 |

`misspell-en-codespell` is deduplicated against `misspell-en` — install both
for maximum coverage without trigger conflicts.

## Usage

Generate all packages:

```bash
go run ./cmd/generate
```

Generate packages for a specific language:

```bash
go run ./cmd/generate de
go run ./cmd/generate en codespell
```

Output is written to `packages/`.

## License

Generator code: [MIT License](LICENSE)

Each generated package carries its own license matching its data source.
See the table above and individual package LICENSE files.
