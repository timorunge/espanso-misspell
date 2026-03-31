# espanso-misspell-en

`espanso-misspell-en` is holding [espanso](https://espanso.org) packages which
are replacing commonly misspelled english words.

The packages are auto-generated based on word lists from
[github.com/client9/misspell](https://github.com/client9/misspell).

## Usage

Generate all espanso packages:

```bash
go run .
```

This creates `misspell-en/`, `misspell-en_UK/`, and `misspell-en_US/`
directories, each containing `package.yml`, `README.md`, and `LICENSE`.

## License

[BSD 3-Clause "New" or "Revised" License](LICENSE)

Misspell is [MIT](https://github.com/client9/misspell/blob/master/LICENSE).
