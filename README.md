# espanso-misspell

`espanso-misspell` generates [espanso](https://espanso.org) packages that
replace commonly misspelled words in multiple languages.

## Packages

| Package | Language | Source | Entries |
|---|---|---|---|
| `misspell-en` | English | [client9/misspell](https://github.com/client9/misspell) DictMain | ~28k |
| `misspell-en_UK` | English | [client9/misspell](https://github.com/client9/misspell) DictBritish | ~1.5k |
| `misspell-en_US` | English | [client9/misspell](https://github.com/client9/misspell) DictAmerican | ~1.6k |
| `misspell-de` | German | [Wikipedia Wortliste](https://de.wikipedia.org/wiki/Wikipedia:Helferlein/Rechtschreibpr%C3%BCfung/Wortliste) | ~5.7k |
| `misspell-es` | Spanish | [Wikipedia Corrector ortográfico](https://es.wikipedia.org/wiki/Wikipedia:Corrector_ortogr%C3%A1fico/Listado) | ~19k |

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
German and Spanish word lists from Wikipedia are
[CC BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
